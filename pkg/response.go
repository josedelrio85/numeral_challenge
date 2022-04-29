package numeral

import (
	"encoding/json"
	"log"
	"net/http"
)

// ResponseAPI represents the data structure needed to create a response
type ResponseAPI struct {
	Status      int               `json:"status"`
	Description string            `json:"description,omitempty"`
	Success     bool              `json:"success"`
	Data        map[string]string `json:"data,omitempty"`
}

// response sets the params to generate a JSON response
func response(w http.ResponseWriter, ra ResponseAPI) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ra.Status)

	json.NewEncoder(w).Encode(ra)
}

// responseError prints the error using log and returns a response
func responseError(w http.ResponseWriter, description string) {
	log.Println(description)

	ra := ResponseAPI{
		Status:      http.StatusInternalServerError,
		Description: description,
		Success:     false,
	}
	response(w, ra)
}

// responseBadRequest calls response function with proper data to generate a Bad Request response
// func responseBadRequest(w http.ResponseWriter) {
// 	ra := ResponseAPI{
// 		Status:  http.StatusBadRequest,
// 		Success: false,
// 	}

// 	response(w, ra)
// }

// responseUnprocessable calls response function with proper data to generate a Unprocessable Entity response
func responseUnprocessable(w http.ResponseWriter, description string) {
	ra := ResponseAPI{
		Status:      http.StatusUnprocessableEntity,
		Description: description,
		Success:     false,
	}

	response(w, ra)
}

// responseOk calls response function with proper data to generate an OK response
func responseOk(w http.ResponseWriter, success bool, userid string) {
	ra := ResponseAPI{
		Status:  http.StatusOK,
		Success: success,
		Data: map[string]string{
			"id": userid,
		},
	}
	response(w, ra)
}

// responseOk calls response function with proper data to generate an Unauthorized response
func responseUnauthorized(w http.ResponseWriter) {
	ra := ResponseAPI{
		Status:      http.StatusUnauthorized,
		Success:     false,
		Description: "Unauthorized",
	}
	response(w, ra)
}
