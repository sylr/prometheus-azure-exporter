package azure

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools/cache"
)

// GetSubscription returns a subscription
func GetSubscription(ctx context.Context, clients *AzureClients, subscriptionID string) (*subscription.Model, error) {
	c := cache.GetCache(30 * time.Second)

	if csub, ok := c.Get(subscriptionID); ok {
		if sub, ok := csub.(*subscription.Model); !ok {
			log.WithField("subscription", subscriptionID).Errorf("Failed to cast object from cache back to *subscription.Model")
		} else {
			return sub, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetSubscriptionClient(subscriptionID)

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	sub, err := client.Get(ctx, subscriptionID)
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		return nil, err
	}

	c.SetDefault(subscriptionID, &sub)

	return &sub, nil
}
