package api

import (
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
)

func GetClient(server string, token string) (*ClientWithResponses, error) {
	apiToken, err := securityprovider.NewSecurityProviderApiKey("header", "X-Algo-API-Token", token)
	if err != nil {
		return nil, err
	}
	return NewClientWithResponses(server, WithRequestEditorFn(apiToken.Intercept))
}
