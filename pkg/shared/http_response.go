package shared

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"reflect"
)

// HTTPResponse format
type HTTPResponse struct {
	Success bool        `json:"success"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Meta    interface{} `json:"meta,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// NewHTTPResponse for create common response
func NewHTTPResponse(code int, message string, params ...interface{}) *HTTPResponse {
	response := new(HTTPResponse)

	for _, param := range params {

		// get value param if type is pointer
		refValue := reflect.ValueOf(param)
		if refValue.Kind() == reflect.Ptr {
			refValue = refValue.Elem()
		}
		param = refValue.Interface()

		switch val := param.(type) {
		case Meta:
			response.Meta = val
		default:
			response.Data = param
		}
	}

	if code < http.StatusBadRequest {
		response.Success = true
	}
	response.Code = code
	response.Message = message
	return response
}

// JSON for set http JSON response (Content-Type: application/json) with parameter is http response writer
func (resp *HTTPResponse) JSON(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)
	return json.NewEncoder(w).Encode(resp)
}

// XML for set http XML response (Content-Type: application/xml)
func (resp *HTTPResponse) XML(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(resp.Code)
	return xml.NewEncoder(w).Encode(resp)
}
