package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// WrapEmptyJSON toma un flujo de bytes y si no hay datos
// permitirá insertar un objeto JSON vacío. Útil
// cuando las llamadas a la API se omiten
func WrapEmptyJSON(data []byte) []byte {
	if len(data) > 0 {
		return data
	}
	return []byte("{}")
}

func JSONError(w http.ResponseWriter, errorCode int, errorMessages ...string) {
	w.WriteHeader(errorCode)
	if len(errorMessages) > 1 {
		err := json.NewEncoder(w).Encode(struct {
			Status string   `json:"status,omitempty"`
			Errors []string `json:"errors,omitempty"`
		}{
			Status: fmt.Sprintf("%d / %s", errorCode, http.StatusText(errorCode)),
			Errors: errorMessages,
		})
		if err != nil {
			log.Fatalf("Error encode struct on JSONError %v", err)
		}
		return
	}

	err := json.NewEncoder(w).Encode(struct {
		Status string `json:"status,omitempty"`
		Error  string `json:"error,omitempty"`
	}{
		Status: fmt.Sprintf("%d / %s", errorCode, http.StatusText(errorCode)),
		Error:  errorMessages[0],
	})
	if err != nil {
		log.Fatalf("Error encode struct on JSONError %v", err)
	}
}

func JSONMessage(w http.ResponseWriter, code int, messages ...string) {
	w.WriteHeader(code)
	if len(messages) > 1 {
		err := json.NewEncoder(w).Encode(struct {
			Status   string   `json:"status,omitempty"`
			Messages []string `json:"messages,omitempty"`
		}{
			Status:   fmt.Sprintf("%d / %s", code, http.StatusText(code)),
			Messages: messages,
		})
		if err != nil {
			log.Fatalf("Error encode struct on JSONError %v", err)
		}
		return
	}

	err := json.NewEncoder(w).Encode(struct {
		Status  string `json:"status,omitempty"`
		Message string `jons:"message,omitempty"`
	}{
		Status:  fmt.Sprintf("%d / %s", code, http.StatusText(code)),
		Message: messages[0],
	})
	if err != nil {
		log.Fatalf("Error encode struct on JSONError %v", err)
	}
}

func PrettyJSON(obj interface{}) []byte {
	prettyJSON, err := json.MarshalIndent(obj, "", " ")
	if err != nil {
		log.Println("Failed to generate json", err)
	}
	return prettyJSON
}
