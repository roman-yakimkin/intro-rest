package handlers

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Code     int
	Message  string
	Response interface{}
}

func returnErrorResponse(isErrorResponse bool, w http.ResponseWriter, r *http.Request, code int, err error, strMsg string) bool {
	if isErrorResponse {
		message := ""
		if strMsg != "" {
			message = strMsg
		} else if err != nil {
			message = err.Error()
		}
		httpResponse := &ErrorResponse{Code: code, Message: message}
		jsonResponse, err := json.Marshal(httpResponse)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(jsonResponse)
	}
	return isErrorResponse
}

func returnSuccessResponse(w http.ResponseWriter, r *http.Request, message string, response interface{}) {
	var successResponse = SuccessResponse{
		Code:     http.StatusOK,
		Response: response,
	}
	successJSONResponse, err := json.Marshal(successResponse)
	if returnErrorResponse(err != nil, w, r, http.StatusInternalServerError, err, "") {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(successJSONResponse)
}
