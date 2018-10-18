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

	var a autorest.Authorizer
	var err error

	a, err = auth.NewAuthorizerFromEnvironment()

	if err == nil {
		// cache
		batchAuthorizer = a
	} else {
		// clear cache
		batchAuthorizer = nil
	}

	return batchAuthorizer, err
}

// GetBatchAuthorizerWithResource get batch authorizer with resource
func GetBatchAuthorizerWithResource(resource string) (autorest.Authorizer, error) {
	if batchAuthorizerWithResource != nil {
		return batchAuthorizerWithResource, nil
	}

	var a autorest.Authorizer
	var err error

	a, err = auth.NewAuthorizerFromEnvironmentWithResource(resource)

	if err == nil {
		// cache
		batchAuthorizerWithResource = a
	} else {
		// clear cache
		batchAuthorizerWithResource = nil
	}

	return batchAuthorizerWithResource, err
}
