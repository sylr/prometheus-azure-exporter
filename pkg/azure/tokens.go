package azure

import (
	"context"
	"os"
	"time"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	log "github.com/sirupsen/logrus"
	"github.com/sylr/prometheus-azure-exporter/pkg/tools"
)

var (
	cacheKeyStorageToken = `adal-token`
)

// GetStorageToken ...
func GetStorageToken(ctx context.Context) (*adal.ServicePrincipalToken, error) {
	c := tools.GetCache(1 * time.Hour)
	cacheKey := cacheKeyStorageToken

	contextLogger := log.WithFields(log.Fields{
		"_id": ctx.Value("id").(string),
	})

	if ctoken, ok := c.Get(cacheKey); ok {
		if token, ok := ctoken.(*adal.ServicePrincipalToken); !ok {
			contextLogger.Errorf("Failed to cast object from cache back to *adal.ServicePrincipalToken")
		} else {
			return token, nil
		}
	}

	envName := os.Getenv("AZURE_ENVIRONMENT")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")

	if len(envName) == 0 {
		envName = azure.PublicCloud.Name
	}

	env, err := azure.EnvironmentFromName(envName)

	if err != nil {
		return nil, err
	}

	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, os.Getenv("AZURE_TENANT_ID"))

	if err != nil {
		return nil, err
	}

	token, err := adal.NewServicePrincipalToken(*oauthConfig, clientID, clientSecret, "https://storage.azure.com/")

	if err != nil {
		return nil, err
	}

	token.SetAutoRefresh(true)
	token.Refresh()

	c.SetDefault(cacheKey, token)

	return token, nil
}
