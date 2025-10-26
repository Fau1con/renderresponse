package renderresponse

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type JSONResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
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

func RenderError(w http.ResponseWriter, message string, status int) {
	response := JSONResponse{
		Status:  "error",
		Message: message,
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
		fmt.Sprintf("Method %s not allowed. Allowed: %v", r.Method, allowedMethods),
		http.StatusMethodNotAllowed,
	)
	return false
}
