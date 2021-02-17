module github.com/sylr/prometheus-azure-exporter

go 1.14

require (
	github.com/Azure/azure-sdk-for-go v50.2.0+incompatible
	github.com/Azure/azure-storage-blob-go v0.13.0
	github.com/Azure/go-autorest/autorest v0.11.18
	github.com/Azure/go-autorest/autorest/adal v0.9.13
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.7
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/jessevdk/go-flags v1.4.0
	github.com/prometheus/client_golang v1.14.1
	github.com/prometheus/common v0.15.0 // indirect
	github.com/prometheus/procfs v0.3.0 // indirect
	github.com/satori/go.uuid v1.2.0 // indirect
	github.com/sirupsen/logrus v1.7.1
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/text v0.3.5 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.4.0
	sylr.dev/libqd/cache v0.0.0-20210116223609-0430c5632a32
	sylr.dev/libqd/sync v0.0.0-20210116223455-05eb9c839987
)

replace github.com/prometheus/client_golang => github.com/sylr/prometheus-client-golang v1.6.1-0.20200515191553-9c85e674da94
