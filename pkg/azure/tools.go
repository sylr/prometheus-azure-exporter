package azure

// Original code from: https://gist.github.com/vladbarosan/fb2528754cbd97df51ca11fe7be27d2f

import (
	"fmt"
	"regexp"
	"strings"
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
func ParseResourceID(resourceID string) (ResourceDetails, error) {
	const resourceIDPatternText = `(?i)subscriptions/([^/]+)/resourceGroups/([^/]+)/providers/([^/]+)/([^/]+)/(.+)`
	resourceIDPattern := regexp.MustCompile(resourceIDPatternText)
	match := resourceIDPattern.FindStringSubmatch(resourceID)

	if len(match) == 0 {
		return ResourceDetails{}, fmt.Errorf("parsing failed for %s. Invalid resource Id format", resourceID)
	}

	v := strings.Split(match[5], "/")
	resourceName := v[len(v)-1]

	result := ResourceDetails{
		SubscriptionID: match[1],
		ResourceGroup:  match[2],
		Provider:       match[3],
		ResourceType:   match[4],
		ResourceName:   resourceName,
	}

	return result, nil
}
