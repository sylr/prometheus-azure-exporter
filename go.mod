module github.com/sylr/prometheus-azure-exporter

go 1.13

require (
	github.com/Azure/azure-sdk-for-go v33.4.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.8.0
	github.com/Azure/go-autorest/autorest v0.9.2
	github.com/Azure/go-autorest/autorest/adal v0.8.0
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.0
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/fsnotify/fsnotify v1.4.7
	github.com/jessevdk/go-flags v1.4.0
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627
	github.com/pkg/errors v0.8.1 // indirect
	github.com/prometheus/client_golang v0.9.3
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	golang.org/x/crypto v0.0.0-20190618222545-ea8f1a30c443 // indirect
	gopkg.in/yaml.v2 v2.2.7
)

replace (
	github.com/patrickmn/go-cache => github.com/sylr/go-cache v2.1.1-0.20190112150453-7f6fb256aaca+incompatible
	github.com/prometheus/client_golang => github.com/sylr/prometheus-client-golang v0.0.0-20190106175946-16e6956cdb08
)
