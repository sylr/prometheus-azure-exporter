module github.com/sylr/prometheus-azure-exporter

require (
	contrib.go.opencensus.io/exporter/ocagent v0.2.0 // indirect
	github.com/Azure/azure-sdk-for-go v24.0.0+incompatible
	github.com/Azure/go-autorest v11.2.8+incompatible
	github.com/census-instrumentation/opencensus-proto v0.0.2 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/dimchansky/utfbom v1.0.0 // indirect
	github.com/gogo/protobuf v1.1.1 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/mitchellh/go-homedir v1.0.0 // indirect
	github.com/patrickmn/go-cache v0.0.0-20180815053127-5633e0862627
	github.com/prometheus/client_golang v0.9.2
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.2.0
	go.opencensus.io v0.18.0 // indirect
	golang.org/x/crypto v0.0.0-20181015023909-0c41d7ab0a0e // indirect
	golang.org/x/sys v0.0.0-20181021155630-eda9bb28ed51 // indirect
	google.golang.org/api v0.0.0-20181021000519-a2651947f503 // indirect
	google.golang.org/genproto v0.0.0-20181016170114-94acd270e44e // indirect
)

replace github.com/prometheus/client_golang => github.com/sylr/prometheus-client-golang v0.0.0-20190106175946-16e6956cdb08
