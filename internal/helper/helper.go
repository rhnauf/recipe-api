package helper

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func HandleResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	res := &Response{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	}

	js, err := json.Marshal(res)
	if err != nil {
		HandleInternalServerError(w)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)
}

func HandleInternalServerError(w http.ResponseWriter) {
	res := &Response{
		StatusCode: http.StatusInternalServerError,
		Message:    "internal server error",
		Data:       nil,
	}

	js, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("error marshalling response")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(js)
}
