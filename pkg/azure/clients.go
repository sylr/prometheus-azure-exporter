package azure

import (
	"github.com/Azure/azure-sdk-for-go/services/batch/2018-08-01.7.0/batch"
	azurebatch "github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
)

type AzureClients struct {
	batchAccountClients map[string]*azurebatch.AccountClient
	batchPoolClients    map[string]*azurebatch.PoolClient
	batchJobClients     map[string]*batch.JobClient
}

// GetNewAzureClients makes new AzureClients object
func GetNewAzureClients() *AzureClients {
	azc := &AzureClients{
		batchAccountClients: make(map[string]*azurebatch.AccountClient),
		batchPoolClients:    make(map[string]*azurebatch.PoolClient),
		batchJobClients:     make(map[string]*batch.JobClient),
	}

	return azc
}

// GetAccountClient return batch account client for specific subscription
func (azc *AzureClients) GetBatchAccountClient(subscriptionID string) (*azurebatch.AccountClient, error) {
	if _, ok := azc.batchAccountClients[subscriptionID]; ok {
		return azc.batchAccountClients[subscriptionID], nil
	}

	auth, err := GetBatchAuthorizer()

	if err != nil {
		return nil, err
	}

	client := azurebatch.NewAccountClient(subscriptionID)
	azc.batchAccountClients[subscriptionID] = &client
	azc.batchAccountClients[subscriptionID].Authorizer = auth

	return azc.batchAccountClients[subscriptionID], nil
}

// GetBatchPoolClient get batch pool client
func (azc *AzureClients) GetBatchPoolClient(subscriptionID string) (*azurebatch.PoolClient, error) {
	if _, ok := azc.batchPoolClients[subscriptionID]; ok {
		return azc.batchPoolClients[subscriptionID], nil
	}

	auth, err := GetBatchAuthorizer()

	if err != nil {
		return nil, err
	}

	client := azurebatch.NewPoolClient(subscriptionID)
	azc.batchPoolClients[subscriptionID] = &client
	azc.batchPoolClients[subscriptionID].Authorizer = auth

	return azc.batchPoolClients[subscriptionID], nil
}

// GetBatchJobClient get batch job client
func (azc *AzureClients) GetBatchJobClient(accountEndpoint string) (*batch.JobClient, error) {
	if _, ok := azc.batchJobClients[accountEndpoint]; ok {
		return azc.batchJobClients[accountEndpoint], nil
	}

	auth, err := GetBatchAuthorizer()

	if err != nil {
		return nil, err
	}

	client := batch.NewJobClientWithBaseURI("https://" + accountEndpoint)
	azc.batchJobClients[accountEndpoint] = &client
	azc.batchJobClients[accountEndpoint].Authorizer = auth

	return azc.batchJobClients[accountEndpoint], nil
}

// GetBatchJobClientWithResource get job client with resource
func (azc *AzureClients) GetBatchJobClientWithResource(accountEndpoint string, resource string) (*batch.JobClient, error) {
	if _, ok := azc.batchJobClients[accountEndpoint+resource]; ok {
		return azc.batchJobClients[accountEndpoint+resource], nil
	}

	auth, err := GetBatchAuthorizerWithResource(resource)

	if err != nil {
		return nil, err
	}

	client := batch.NewJobClientWithBaseURI("https://" + accountEndpoint)
	azc.batchJobClients[accountEndpoint+resource] = &client
	azc.batchJobClients[accountEndpoint+resource].Authorizer = auth

	return azc.batchJobClients[accountEndpoint+resource], nil
}
