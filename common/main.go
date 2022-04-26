package common

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/security-headers-backend/entities"
)

func WriteError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	errorResponse := entities.GenericApiError{
		ErrorCode:    "BAD_BODY",
		ErrorMessage: err.Error(),
	}
	json.NewEncoder(w).Encode(errorResponse)
}
