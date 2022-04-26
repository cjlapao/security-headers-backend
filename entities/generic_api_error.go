package entities

type GenericApiError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"message"`
}
