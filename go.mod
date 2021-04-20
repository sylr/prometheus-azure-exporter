module github.com/sylr/prometheus-azure-exporter

go 1.14

require (
	github.com/Azure/azure-sdk-for-go v53.3.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.13.0
	github.com/Azure/go-autorest/autorest v0.11.18
	github.com/Azure/go-autorest/autorest/adal v0.9.13
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/jessevdk/go-flags v1.5.0
	github.com/prometheus/client_golang v1.14.1
	github.com/prometheus/common v0.20.0 // indirect
	github.com/prometheus/procfs v0.6.0 // indirect
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/crypto v0.0.0-20210415154028-4f45737414dc // indirect
	golang.org/x/net v0.0.0-20210415231046-e915ea6b2b7d // indirect
	golang.org/x/sys v0.0.0-20210415045647-66c3f260301c // indirect
	gopkg.in/yaml.v2 v2.4.0
	sylr.dev/libqd/cache v0.0.0-20210116223609-0430c5632a32
	sylr.dev/libqd/sync v0.0.0-20210116223455-05eb9c839987
)

replace github.com/prometheus/client_golang => github.com/sylr/prometheus-client-golang v1.6.1-0.20200515191553-9c85e674da94
