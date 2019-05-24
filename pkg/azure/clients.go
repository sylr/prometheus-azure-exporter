package azure

import (
	"net/http"
	"os"
	"sync"

	"github.com/Azure/azure-sdk-for-go/services/batch/2018-08-01.7.0/batch"
	azurebatch "github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	graph "github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2018-07-01/storage"
	"github.com/Azure/go-autorest/autorest"
	log "github.com/sirupsen/logrus"
)

var (
	mutex = sync.RWMutex{}
)

// AzureClients Collection of Azure clients
type AzureClients struct {
	mutex                       sync.RWMutex
	batchAccountClients         map[string]*azurebatch.AccountClient
	batchPoolClients            map[string]*azurebatch.PoolClient
	batchJobClients             map[string]*batch.JobClient
	subscriptionsClients        map[string]*subscription.SubscriptionsClient
	applicationsClients         map[string]*graph.ApplicationsClient
	servicePrincipalsClients    map[string]*graph.ServicePrincipalsClient
	storageAccountsClients      map[string]*storage.AccountsClient
	storageAccountUsagesClients map[string]*storage.UsagesClient
	blobContainersClients       map[string]*storage.BlobContainersClient
	groupClients                map[string]*resources.GroupsClient
}

// NewAzureClients makes new AzureClients object
func NewAzureClients() *AzureClients {
	azc := &AzureClients{
		mutex:                       sync.RWMutex{},
		batchAccountClients:         make(map[string]*azurebatch.AccountClient),
		batchPoolClients:            make(map[string]*azurebatch.PoolClient),
		batchJobClients:             make(map[string]*batch.JobClient),
		subscriptionsClients:        make(map[string]*subscription.SubscriptionsClient),
		applicationsClients:         make(map[string]*graph.ApplicationsClient),
		servicePrincipalsClients:    make(map[string]*graph.ServicePrincipalsClient),
		storageAccountsClients:      make(map[string]*storage.AccountsClient),
		storageAccountUsagesClients: make(map[string]*storage.UsagesClient),
		blobContainersClients:       make(map[string]*storage.BlobContainersClient),
		groupClients:                make(map[string]*resources.GroupsClient),
	}

	return azc
}

// GetSubscriptionClient return subscription client
func (azc *AzureClients) GetSubscriptionClient(subscriptionID string) (*subscription.SubscriptionsClient, error) {
	if _, ok := azc.subscriptionsClients[subscriptionID]; ok {
		return azc.subscriptionsClients[subscriptionID], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetBatchAuthorizer()

	if err != nil {
		return nil, err
	}

	client := subscription.NewSubscriptionsClient()
	azc.subscriptionsClients[subscriptionID] = &client
	azc.subscriptionsClients[subscriptionID].Authorizer = auth
	azc.subscriptionsClients[subscriptionID].ResponseInspector = respondInspect(subscriptionID)

	return azc.subscriptionsClients[subscriptionID], nil
}

// GetGroupClient return group client
func (azc *AzureClients) GetGroupClient(subscriptionID string) (*resources.GroupsClient, error) {
	if _, ok := azc.groupClients[subscriptionID]; ok {
		return azc.groupClients[subscriptionID], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetAuthorizer()

	if err != nil {
		return nil, err
	}

	client := resources.NewGroupsClient(subscriptionID)
	azc.groupClients[subscriptionID] = &client
	azc.groupClients[subscriptionID].Authorizer = auth
	azc.groupClients[subscriptionID].ResponseInspector = respondInspect(subscriptionID)

	return azc.groupClients[subscriptionID], nil
}

// GetBatchAccountClient return batch account client for specific subscription
func (azc *AzureClients) GetBatchAccountClient(subscriptionID string) (*azurebatch.AccountClient, error) {
	if _, ok := azc.batchAccountClients[subscriptionID]; ok {
		return azc.batchAccountClients[subscriptionID], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetBatchAuthorizer()

	if err != nil {
		return nil, err
	}

	client := azurebatch.NewAccountClient(subscriptionID)
	azc.batchAccountClients[subscriptionID] = &client
	azc.batchAccountClients[subscriptionID].Authorizer = auth
	azc.batchAccountClients[subscriptionID].ResponseInspector = respondInspect(subscriptionID)

	return azc.batchAccountClients[subscriptionID], nil
}

// GetBatchPoolClient get batch pool client
func (azc *AzureClients) GetBatchPoolClient(subscriptionID string) (*azurebatch.PoolClient, error) {
	if _, ok := azc.batchPoolClients[subscriptionID]; ok {
		return azc.batchPoolClients[subscriptionID], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetBatchAuthorizer()

	if err != nil {
		return nil, err
	}

	client := azurebatch.NewPoolClient(subscriptionID)
	azc.batchPoolClients[subscriptionID] = &client
	azc.batchPoolClients[subscriptionID].Authorizer = auth
	azc.batchPoolClients[subscriptionID].ResponseInspector = respondInspect(subscriptionID)

	return azc.batchPoolClients[subscriptionID], nil
}

// GetBatchJobClient get batch job client
func (azc *AzureClients) GetBatchJobClient(accountEndpoint string) (*batch.JobClient, error) {
	if _, ok := azc.batchJobClients[accountEndpoint]; ok {
		return azc.batchJobClients[accountEndpoint], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetBatchAuthorizer()

	if err != nil {
		return nil, err
	}

	client := batch.NewJobClientWithBaseURI("https://" + accountEndpoint)
	azc.batchJobClients[accountEndpoint] = &client
	azc.batchJobClients[accountEndpoint].Authorizer = auth
	// azc.batchJobClients[accountEndpoint].ResponseInspector = respondInspectDebug()

	return azc.batchJobClients[accountEndpoint], nil
}

// GetBatchJobClientWithResource get job client with resource
func (azc *AzureClients) GetBatchJobClientWithResource(accountEndpoint string, resource string) (*batch.JobClient, error) {
	if _, ok := azc.batchJobClients[accountEndpoint+resource]; ok {
		return azc.batchJobClients[accountEndpoint+resource], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetBatchAuthorizerWithResource(resource)

	if err != nil {
		return nil, err
	}

	client := batch.NewJobClientWithBaseURI("https://" + accountEndpoint)
	azc.batchJobClients[accountEndpoint+resource] = &client
	azc.batchJobClients[accountEndpoint+resource].Authorizer = auth
	// azc.batchJobClients[accountEndpoint+resource].ResponseInspector = respondInspectDebug()

	return azc.batchJobClients[accountEndpoint+resource], nil
}

// GetApplicationsClient get applications client
func (azc *AzureClients) GetApplicationsClient(tenantID string) (*graph.ApplicationsClient, error) {
	if _, ok := azc.applicationsClients[tenantID]; ok {
		return azc.applicationsClients[tenantID], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetGraphAuthorizer()

	if err != nil {
		return nil, err
	}

	client := graph.NewApplicationsClient(tenantID)
	azc.applicationsClients[tenantID] = &client
	azc.applicationsClients[tenantID].Authorizer = auth
	// azc.applicationsClients[tenantID].ResponseInspector = respondInspectDebug()

	return azc.applicationsClients[tenantID], nil
}

// GetStorageAccountsClient get storage account client
func (azc *AzureClients) GetStorageAccountsClient(subscriptionID string) (*storage.AccountsClient, error) {
	if _, ok := azc.storageAccountsClients[subscriptionID]; ok {
		return azc.storageAccountsClients[subscriptionID], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetStorageAuthorizer()

	if err != nil {
		return nil, err
	}

	client := storage.NewAccountsClient(subscriptionID)
	azc.storageAccountsClients[subscriptionID] = &client
	azc.storageAccountsClients[subscriptionID].Authorizer = auth
	azc.storageAccountsClients[subscriptionID].ResponseInspector = respondInspect(subscriptionID)

	return azc.storageAccountsClients[subscriptionID], nil
}

// GetStorageAccountsClientWithResource get storage account client
func (azc *AzureClients) GetStorageAccountsClientWithResource(subscriptionID string, accountEndpoint string, resource string) (*storage.AccountsClient, error) {
	if _, ok := azc.storageAccountsClients[accountEndpoint+resource]; ok {
		return azc.storageAccountsClients[accountEndpoint+resource], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetStorageAuthorizerWithResource(resource)

	if err != nil {
		return nil, err
	}

	client := storage.NewAccountsClientWithBaseURI(accountEndpoint, subscriptionID)
	azc.storageAccountsClients[accountEndpoint+resource] = &client
	azc.storageAccountsClients[accountEndpoint+resource].Authorizer = auth
	azc.storageAccountsClients[accountEndpoint+resource].ResponseInspector = respondInspect(subscriptionID)

	return azc.storageAccountsClients[accountEndpoint+resource], nil
}

// GetStorageAccountUsagesClient get storage account client
func (azc *AzureClients) GetStorageAccountUsagesClient(subscriptionID string) (*storage.UsagesClient, error) {
	if _, ok := azc.storageAccountUsagesClients[subscriptionID]; ok {
		return azc.storageAccountUsagesClients[subscriptionID], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetStorageAuthorizer()

	if err != nil {
		return nil, err
	}

	client := storage.NewUsagesClient(subscriptionID)
	azc.storageAccountUsagesClients[subscriptionID] = &client
	azc.storageAccountUsagesClients[subscriptionID].Authorizer = auth
	azc.storageAccountUsagesClients[subscriptionID].ResponseInspector = respondInspect(subscriptionID)

	return azc.storageAccountUsagesClients[subscriptionID], nil
}

// GetBlobContainersClient get storage account client
func (azc *AzureClients) GetBlobContainersClient(subscriptionID string) (*storage.BlobContainersClient, error) {
	if _, ok := azc.blobContainersClients[subscriptionID]; ok {
		return azc.blobContainersClients[subscriptionID], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetStorageAuthorizer()

	if err != nil {
		return nil, err
	}

	client := storage.NewBlobContainersClient(subscriptionID)
	azc.blobContainersClients[subscriptionID] = &client
	azc.blobContainersClients[subscriptionID].Authorizer = auth
	azc.blobContainersClients[subscriptionID].ResponseInspector = respondInspect(subscriptionID)

	return azc.blobContainersClients[subscriptionID], nil
}

// GetBlobContainersClientWithResource get storage account client
func (azc *AzureClients) GetBlobContainersClientWithResource(subscriptionID string, accountEndpoint string, resource string) (*storage.BlobContainersClient, error) {
	if _, ok := azc.blobContainersClients[accountEndpoint+resource]; ok {
		return azc.blobContainersClients[accountEndpoint+resource], nil
	}

	azc.mutex.Lock()
	defer azc.mutex.Unlock()

	auth, err := GetStorageAuthorizerWithResource(resource)

	if err != nil {
		return nil, err
	}

	client := storage.NewBlobContainersClientWithBaseURI(accountEndpoint, subscriptionID)
	azc.blobContainersClients[accountEndpoint+resource] = &client
	azc.blobContainersClients[accountEndpoint+resource].Authorizer = auth
	azc.blobContainersClients[accountEndpoint+resource].ResponseInspector = respondInspect(subscriptionID)

	return azc.blobContainersClients[accountEndpoint+resource], nil
}

// ----------------------------------------------------------------------------

func respondInspect(subscription string) autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {
			SetReadRateLimitRemaining(os.Getenv("AZURE_TENANT_ID"), subscription, resp)
			SetWriteRateLimitRemaining(os.Getenv("AZURE_TENANT_ID"), subscription, resp)
			return r.Respond(resp)
		})
	}
}

func respondInspectDebug() autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {
			for key, val := range resp.Header {
				log.Debugf("HEADER %v: %v", key, val)
			}
			return r.Respond(resp)
		})
	}
}
