module github.com/sylr/prometheus-azure-exporter

go 1.17

require (
	github.com/Azure/azure-sdk-for-go v65.0.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.15.0
	github.com/Azure/go-autorest/autorest v0.11.27
	github.com/Azure/go-autorest/autorest/adal v0.9.20
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.11
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/fsnotify/fsnotify v1.5.4
	github.com/gofrs/uuid v4.1.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jessevdk/go-flags v1.5.0
	github.com/prometheus/client_golang v1.14.1
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	gopkg.in/yaml.v2 v2.4.0
	sylr.dev/libqd/cache v0.0.0-20210116223609-0430c5632a32
	sylr.dev/libqd/sync v0.0.0-20210116223455-05eb9c839987
)

require (
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.5 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/golang-jwt/jwt/v4 v4.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
	sylr.dev/cache/v2 v2.3.0 // indirect
)

replace github.com/prometheus/client_golang => github.com/sylr/prometheus-client-golang v1.6.1-0.20200515191553-9c85e674da94
