package handlers

import (
	api "lab-tracker/internal/delivery/http/gen"
	"lab-tracker/internal/domain"
)

func newErrorResponse(code api.ErrorResponseErrorCode, message string) api.ErrorResponse {
	return api.ErrorResponse{
		Error: struct {
			Code    api.ErrorResponseErrorCode `json:"code"`
			Message string                     `json:"message"`
		}{
			Code:    code,
			Message: message,
		},
	}
}

func unauthorized() api.ErrorResponse {
	return newErrorResponse(api.ErrorResponseErrorCodeUNAUTHORIZED, domain.ErrUnauthorized.Error())
}

func forbidden() api.ErrorResponse {
	return newErrorResponse(api.ErrorResponseErrorCodeFORBIDDEN, domain.ErrForbidden.Error())
}

func invalidRequest(message string) api.ErrorResponse {
	if message == "" {
		message = domain.ErrInvalidRequest.Error()
	}
	return newErrorResponse(api.ErrorResponseErrorCodeINVALIDREQUEST, message)
}

func userNotFound() api.ErrorResponse {
	return newErrorResponse(api.ErrorResponseErrorCodeUSERNOTFOUND, domain.ErrUserNotFound.Error())
}

func groupNotFound() api.ErrorResponse {
	return newErrorResponse(api.ErrorResponseErrorCodeGROUPNOTFOUND, domain.ErrGroupNotFound.Error())
}

func labWorkNotFound() api.ErrorResponse {
	return newErrorResponse(api.ErrorResponseErrorCodeLABWORKNOTFOUND, domain.ErrLabWorkNotFound.Error())
}

func internalError() api.ErrorResponse {
	return newErrorResponse(api.ErrorResponseErrorCodeINTERNALERROR, domain.ErrInternal.Error())
}
