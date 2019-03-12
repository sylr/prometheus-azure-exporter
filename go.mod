module github.com/sylr/prometheus-azure-exporter

require (
	contrib.go.opencensus.io/exporter/ocagent v0.2.0 // indirect
	github.com/Azure/azure-sdk-for-go v24.1.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.0.0-20190104215108-45d0c5e3638e
	github.com/Azure/go-autorest v11.5.0+incompatible
	github.com/census-instrumentation/opencensus-proto v0.0.2 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dimchansky/utfbom v1.0.0 // indirect
	github.com/fsnotify/fsnotify v1.4.7
	github.com/jessevdk/go-flags v1.4.0
	github.com/mitchellh/go-homedir v1.0.0 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20161129095857-cc309e4a2223 // indirect
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627
	github.com/prometheus/client_golang v0.9.2
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.4.0
	go.opencensus.io v0.18.0 // indirect
	golang.org/x/crypto v0.0.0-20181015023909-0c41d7ab0a0e // indirect
	golang.org/x/sys v0.0.0-20181021155630-eda9bb28ed51 // indirect
	google.golang.org/api v0.0.0-20181021000519-a2651947f503 // indirect
	google.golang.org/genproto v0.0.0-20181016170114-94acd270e44e // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	github.com/patrickmn/go-cache => github.com/sylr/go-cache v2.1.1-0.20190112150453-7f6fb256aaca+incompatible
	github.com/prometheus/client_golang => github.com/sylr/prometheus-client-golang v0.0.0-20190106175946-16e6956cdb08
)
