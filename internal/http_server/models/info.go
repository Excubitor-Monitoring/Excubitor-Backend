package models

import ctx "github.com/Excubitor-Monitoring/Excubitor-Backend/internal/context"

type InfoResponse struct {
	Authentication Authentication `json:"authentication"`
	Modules        []ctx.Module   `json:"modules"`
}

type Authentication struct {
	Method string `json:"method"`
}

func NewInfoResponse(authenticationMethod string, modules []ctx.Module) InfoResponse {
	return InfoResponse{
		Authentication: Authentication{
			Method: authenticationMethod,
		},
		Modules: modules,
	}
}
