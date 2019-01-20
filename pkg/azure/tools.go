package azure

// Original code from: https://gist.github.com/vladbarosan/fb2528754cbd97df51ca11fe7be27d2f

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/sylr/prometheus-azure-exporter/pkg/tools"
)

const (
	cacheKeyResourceID    = `resource-id-%s`
	resourceIDPatternText = `(?i)subscriptions/([^/]+)/resourceGroups/([^/]+)/providers/([^/]+)/([^/]+)/(.+)`
)

var (
	cache             = tools.GetCache(1 * time.Hour)
	resourceIDPattern = regexp.MustCompile(resourceIDPatternText)
)

// ResourceDetails contains details about an Azure resource
type ResourceDetails struct {
	SubscriptionID string
	ResourceGroup  string
	Provider       string
	ResourceType   string
	ResourceName   string
}

// ParseResourceID parses a resource ID into a ResourceDetails struct
func ParseResourceID(resourceID string) (*ResourceDetails, error) {
	cacheKey := fmt.Sprintf(cacheKeyResourceID, resourceID)

	if cdetails, ok := cache.Get(cacheKey); ok {
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
		ResourceType:   match[4],
		ResourceName:   resourceName,
	}

	cache.SetDefault(cacheKey, details)

	return details, nil
}
