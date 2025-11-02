package renderresponse

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ErrorDetails struct {
	Code    string `jsone:"code,omitempty"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

type JSONResponse struct {
	Status  string         `json:"status"`
	Data    interface{}    `json:"data,omitempty"`
	Message string         `json:"message,omitempty"`
	Errors  []ErrorDetails `json:"errors,omitempty"`
}

// renderResponse - базовая функция для отправки JSON
func renderResponse(w http.ResponseWriter, response JSONResponse, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "failed to encode JSON: %v", err)
	}
}

func RenderJSON(w http.ResponseWriter, data interface{}, status int) {
	response := JSONResponse{
		Status: "success",
		Data:   data,
	}
	renderResponse(w, response, status)
}

func RenderError(w http.ResponseWriter, message string, status int, errors ...error) {
	response := JSONResponse{
		Status:  "error",
		Message: message,
	}

	if len(errors) > 0 {
		errorDetails := make([]ErrorDetails, len(errors))
		for i, err := range errors {
			errorDetails[i] = ErrorDetails{
				Message: err.Error(),
			}
		}
		response.Errors = errorDetails
	}

	renderResponse(w, response, status)
}

// ValidateMethod Валидатор метода
func ValidateMethod(w http.ResponseWriter, r *http.Request, allowedMethods ...string) bool {
	for _, method := range allowedMethods {
		if r.Method == method {
			return true
		}
	}
	RenderError(w,
		fmt.Sprintf("Method %s is not allowed", r.Method),
		http.StatusMethodNotAllowed,
		fmt.Errorf("allowed methods: %v", allowedMethods),
	)
	return false
}
