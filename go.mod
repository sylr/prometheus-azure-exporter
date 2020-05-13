module github.com/sylr/prometheus-azure-exporter

go 1.13

require (
	github.com/Azure/azure-pipeline-go v0.2.2 // indirect
	github.com/Azure/azure-sdk-for-go v40.6.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.8.0
	github.com/Azure/go-autorest/autorest v0.10.1
	github.com/Azure/go-autorest/autorest/adal v0.8.3
	github.com/Azure/go-autorest/autorest/azure/auth v0.4.2
	github.com/Azure/go-autorest/autorest/to v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.2.0 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/golang/protobuf v1.3.4 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/mattn/go-ieproxy v0.0.0-20200203040449-2dbc853185d9 // indirect
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627
	github.com/prometheus/client_golang v1.14.1
	github.com/prometheus/procfs v0.0.10 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/crypto v0.0.0-20200221231518-2aa609cf4a9d // indirect
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
	golang.org/x/sys v0.0.0-20200302083256-062a44052db1 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/patrickmn/go-cache => github.com/sylr/go-cache v2.1.1-0.20190112150453-7f6fb256aaca+incompatible
	github.com/prometheus/client_golang => github.com/sylr/prometheus-client-golang v1.4.2-0.20200226090308-80f92138efbf
)
