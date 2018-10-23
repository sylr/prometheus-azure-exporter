package azure

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
)

var (
	subscriptions = make(map[string]*subscription.Model)
)

// GetSubscription
func GetSubscription(ctx context.Context, clients *AzureClients, subscriptionID string) (*subscription.Model, error) {
	if _, ok := subscriptions[subscriptionID]; ok {
		return subscriptions[subscriptionID], nil
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetSubscriptionClient(subscriptionID)

	if err != nil {
		return nil, err
	}

	sub, err := client.Get(ctx, subscriptionID)
	AzureAPICallsTotal.WithLabelValues().Inc()

	if err != nil {
		AzureAPICallsFailedTotal.WithLabelValues().Inc()
		return nil, err
	}

	subscriptions[subscriptionID] = &sub

	return subscriptions[subscriptionID], nil
}
