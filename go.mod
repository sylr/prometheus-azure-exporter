module github.com/sylr/prometheus-azure-exporter

go 1.14

require (
	github.com/Azure/azure-pipeline-go v0.2.2 // indirect
	github.com/Azure/azure-sdk-for-go v42.2.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.8.0
	github.com/Azure/go-autorest/autorest v0.10.1
	github.com/Azure/go-autorest/autorest/adal v0.8.3
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627
	github.com/prometheus/client_golang v1.14.1
	github.com/prometheus/common v0.10.0 // indirect
	github.com/prometheus/procfs v0.0.11 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37 // indirect
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120 // indirect
	golang.org/x/sys v0.0.0-20200515095857-1151b9dac4a9 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/patrickmn/go-cache => github.com/sylr/go-cache v2.1.1-0.20190112150453-7f6fb256aaca+incompatible
	github.com/prometheus/client_golang => github.com/sylr/prometheus-client-golang v1.6.1-0.20200515191553-9c85e674da94
)
