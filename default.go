package errors

import "net/http"

var (
	DefaultHTTPErrorCode       int    = http.StatusBadRequest
	DefaultMaskedHTTPErrorCode int    = http.StatusInternalServerError
	DefaultMaskMessage         string = "Sorry, there are internal server error occured, please try again later. (%d)"
)
