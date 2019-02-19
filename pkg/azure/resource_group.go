package azure

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/subscription/mgmt/2018-03-01-preview/subscription"
	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/resources"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools/cache"
)

var (
	cacheKeyResourceGroup = "sub-%s-rg-%s"
)

// GetResourceGroup returns a Group
func GetResourceGroup(ctx context.Context, clients *AzureClients, subscription *subscription.Model, name string) (*resources.Group, error) {
	c := cache.GetCache(1 * time.Hour)
	cacheKey := fmt.Sprintf(cacheKeyResourceGroup, *subscription.DisplayName, name)

	if cgroup, ok := c.Get(cacheKey); ok {
		if group, ok := cgroup.(*resources.Group); !ok {
			log.Errorf("Failed to cast object from cache back to *resources.Group")
		} else {
			return group, nil
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	client, err := clients.GetGroupClient(*subscription.SubscriptionID)

	if err != nil {
		return nil, err
	}

	t0 := time.Now()
	group, err := client.Get(ctx, name)
	t1 := time.Since(t0).Seconds()

	ObserveAzureAPICall(t1)

	if err != nil {
		ObserveAzureAPICallFailed(t1)
		return nil, err
	}

	c.SetDefault(cacheKey, &group)

	return &group, nil
}
