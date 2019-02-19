package azure

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/sylr/prometheus-azure-exporter/pkg/tools/cache"
)

const (
	cacheKeyResourceID    = `resource-id-%s`
	resourceIDPatternText = `(?i)subscriptions/([^/]+)/resourceGroups/([^/]+)/providers/([^/]+)/([^/]+)/(.+)`
)

var (
	c                 = cache.GetCache(1 * time.Hour)
	resourceIDPattern = regexp.MustCompile(resourceIDPatternText)
)

// ResourceDetails contains details about an Azure resource
type ResourceDetails struct {
	SubscriptionID string
	ResourceGroup  string
	Provider       string
	Type           string
	Name           string
}

// ParseResourceID parses a resource ID into a ResourceDetails struct
// Original code from: https://gist.github.com/vladbarosan/fb2528754cbd97df51ca11fe7be27d2f
func ParseResourceID(resourceID string) (*ResourceDetails, error) {
	cacheKey := fmt.Sprintf(cacheKeyResourceID, resourceID)

	if cdetails, ok := c.Get(cacheKey); ok {
		if details, ok := cdetails.(*ResourceDetails); ok {
			return details, nil
		}
	}

	match := resourceIDPattern.FindStringSubmatch(resourceID)

	if len(match) == 0 {
		return nil, fmt.Errorf("parsing failed for %s. Invalid resource Id format", resourceID)
	}

	v := strings.Split(match[5], "/")
	resourceName := v[len(v)-1]

	details := &ResourceDetails{
		SubscriptionID: match[1],
		ResourceGroup:  match[2],
		Provider:       match[3],
		Type:           match[4],
		Name:           resourceName,
	}

	c.SetDefault(cacheKey, details)

	return details, nil
}
