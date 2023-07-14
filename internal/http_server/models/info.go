package models

import (
	"github.com/Excubitor-Monitoring/Excubitor-Backend/pkg/shared/modules"
)

type InfoResponse struct {
	Authentication Authentication   `json:"authentication"`
	Modules        []modules.Module `json:"modules"`
}

type Authentication struct {
	Method string `json:"method"`
}

func NewInfoResponse(authenticationMethod string, modules []modules.Module) InfoResponse {
	return InfoResponse{
		Authentication: Authentication{
			Method: authenticationMethod,
		},
		Modules: modules,
	}
}
