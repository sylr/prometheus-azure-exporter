module github.com/sylr/prometheus-azure-exporter

go 1.17

require (
	github.com/Azure/azure-sdk-for-go v53.4.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.14.0
	github.com/Azure/go-autorest/autorest v0.11.22
	github.com/Azure/go-autorest/autorest/adal v0.9.16
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.8
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/fsnotify/fsnotify v1.5.1
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jessevdk/go-flags v1.5.0
	github.com/prometheus/client_golang v1.14.1
	github.com/prometheus/common v0.20.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/net v0.0.0-20210415231046-e915ea6b2b7d // indirect
	gopkg.in/yaml.v2 v2.4.0
	sylr.dev/libqd/cache v0.0.0-20210116223609-0430c5632a32
	sylr.dev/libqd/sync v0.0.0-20210116223455-05eb9c839987
)

require (
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.2 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.0.0 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
	sylr.dev/cache/v2 v2.3.0 // indirect
)

replace github.com/prometheus/client_golang => github.com/sylr/prometheus-client-golang v1.6.1-0.20200515191553-9c85e674da94
