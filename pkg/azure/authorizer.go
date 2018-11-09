package azure

import (
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var (
	batchAuthorizer             autorest.Authorizer
	batchAuthorizerWithResource autorest.Authorizer
)

// GetBatchAuthorizer get batch authorizer
func GetBatchAuthorizer() (autorest.Authorizer, error) {
	if batchAuthorizer != nil {
		return batchAuthorizer, nil
	}

	var err error

	batchAuthorizer, err = auth.NewAuthorizerFromEnvironment()

	if err != nil {
		return nil, err
	}

	return batchAuthorizer, err
}

// GetBatchAuthorizerWithResource get batch authorizer with resource
func GetBatchAuthorizerWithResource(resource string) (autorest.Authorizer, error) {
	if batchAuthorizerWithResource != nil {
		return batchAuthorizerWithResource, nil
	}

	var err error

	batchAuthorizerWithResource, err = auth.NewAuthorizerFromEnvironmentWithResource(resource)

	if err != nil {
		return nil, err
	}

	return batchAuthorizerWithResource, err
}
