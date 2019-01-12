package azure

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/Azure/azure-sdk-for-go/services/batch/2018-08-01.7.0/batch"
	azurebatch "github.com/Azure/azure-sdk-for-go/services/batch/mgmt/2017-09-01/batch"
	graph "github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	"github.com/Azure/go-autorest/autorest"
)

// AzureClients Collection of Azure clients
type AzureClients struct {
	batchAccountClients      map[string]*azurebatch.AccountClient
	batchPoolClients         map[string]*azurebatch.PoolClient
	batchJobClients          map[string]*batch.JobClient
	subscriptionsClients     map[string]*subscription.SubscriptionsClient
	applicationsClients      map[string]*graph.ApplicationsClient
	servicePrincipalsClients map[string]*graph.ServicePrincipalsClient
}

// NewAzureClients makes new AzureClients object
func NewAzureClients() *AzureClients {
	azc := &AzureClients{
		batchAccountClients:      make(map[string]*azurebatch.AccountClient),
		batchPoolClients:         make(map[string]*azurebatch.PoolClient),
		batchJobClients:          make(map[string]*batch.JobClient),
		subscriptionsClients:     make(map[string]*subscription.SubscriptionsClient),
		applicationsClients:      make(map[string]*graph.ApplicationsClient),
		servicePrincipalsClients: make(map[string]*graph.ServicePrincipalsClient),
	}

	return azc
}

// GetSubscriptionClient return subscription client
func (azc *AzureClients) GetSubscriptionClient(subscriptionID string) (*subscription.SubscriptionsClient, error) {
	if _, ok := azc.subscriptionsClients[subscriptionID]; ok {
		return azc.subscriptionsClients[subscriptionID], nil
	}

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

// GetBatchAccountClient return batch account client for specific subscription
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
	azc.batchAccountClients[subscriptionID].ResponseInspector = respondInspect(subscriptionID)

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
	azc.batchPoolClients[subscriptionID].ResponseInspector = respondInspect(subscriptionID)

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
	//azc.batchJobClients[accountEndpoint].ResponseInspector = respondInspectDebug()

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
	//azc.batchJobClients[accountEndpoint+resource].ResponseInspector = respondInspectDebug()

	return azc.batchJobClients[accountEndpoint+resource], nil
}

// GetApplicationsClient get applications client
func (azc *AzureClients) GetApplicationsClient(tenantID string) (*graph.ApplicationsClient, error) {
	if _, ok := azc.applicationsClients[tenantID]; ok {
		return azc.applicationsClients[tenantID], nil
	}

	auth, err := GetGraphAuthorizer()

	if err != nil {
		return nil, err
	}

	client := graph.NewApplicationsClient(tenantID)
	azc.applicationsClients[tenantID] = &client
	azc.applicationsClients[tenantID].Authorizer = auth
	//azc.applicationsClients[tenantID].ResponseInspector = respondInspectDebug()

	return azc.applicationsClients[tenantID], nil
}

func respondInspect(subscription string) autorest.RespondDecorator {
	return func(r autorest.Responder) autorest.Responder {
		return autorest.ResponderFunc(func(resp *http.Response) error {
			SetReadRateLimitRemaining(subscription, resp)
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
