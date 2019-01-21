package azure

import (
	"os"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

var (
	authorizer                    autorest.Authorizer
	graphAuthorizer               autorest.Authorizer
	batchAuthorizer               autorest.Authorizer
	batchAuthorizerWithResource   autorest.Authorizer
	storageAuthorizer             autorest.Authorizer
	storageAuthorizerWithResource autorest.Authorizer
)

// GetAuthorizer get graph authorizer
func GetAuthorizer() (autorest.Authorizer, error) {
	if authorizer != nil {
		return authorizer, nil
	}

	var err error

	authorizer, err = auth.NewAuthorizerFromEnvironment()

	if err != nil {
		return nil, err
	}

	return authorizer, err
}

// GetGraphAuthorizer get graph authorizer
func GetGraphAuthorizer() (autorest.Authorizer, error) {
	if graphAuthorizer != nil {
		return graphAuthorizer, nil
	}

	var err error

	envName := os.Getenv("AZURE_ENVIRONMENT")

	if len(envName) == 0 {
		envName = azure.PublicCloud.Name
	}

	env, err := azure.EnvironmentFromName(envName)

	graphAuthorizer, err = auth.NewAuthorizerFromEnvironmentWithResource(env.GraphEndpoint)

	if err != nil {
		return nil, err
	}

	return graphAuthorizer, err
}

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

// GetStorageAuthorizer get storage authorizer
func GetStorageAuthorizer() (autorest.Authorizer, error) {
	if storageAuthorizer != nil {
		return storageAuthorizer, nil
	}

	var err error

	storageAuthorizer, err = auth.NewAuthorizerFromEnvironment()

	if err != nil {
		return nil, err
	}

	return storageAuthorizer, err
}

// GetStorageAuthorizerWithResource get storage authorizer with resource
func GetStorageAuthorizerWithResource(resource string) (autorest.Authorizer, error) {
	if storageAuthorizerWithResource != nil {
		return storageAuthorizerWithResource, nil
	}

	var err error

	storageAuthorizerWithResource, err = auth.NewAuthorizerFromEnvironmentWithResource(resource)

	if err != nil {
		return nil, err
	}

	return storageAuthorizerWithResource, err
}
