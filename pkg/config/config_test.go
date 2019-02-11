package config

import (
	"testing"
)

func TestMustDiscoverBasedOnTags(t *testing.T) {
	var tags map[string]*string
	strue := "True"
	sfalse := "False"
	tag := "prometheus_io_azure_exporter_discover"

	// AutoDiscoveryMode: All
	// tag: True
	CurrentConfig = &PrometheusAzureExporterConfig{
		AutoDiscoveryMode: "All",
		AutoDiscoveryTag:  tag,
	}

	tags = map[string]*string{
		tag: &strue,
	}

	if b := MustDiscoverBasedOnTags(tags); !b {
		t.Fatalf("Expected %v but got %v", true, b)
	}

	// AutoDiscoveryMode: All
	// tag: False
	CurrentConfig = &PrometheusAzureExporterConfig{
		AutoDiscoveryMode: "All",
		AutoDiscoveryTag:  tag,
	}

	tags = map[string]*string{
		tag: &sfalse,
	}

	if b := MustDiscoverBasedOnTags(tags); b {
		t.Fatalf("Expected %v but got %v", false, b)
	}

	// AutoDiscoveryMode: All
	// tag: nil
	CurrentConfig = &PrometheusAzureExporterConfig{
		AutoDiscoveryMode: "All",
		AutoDiscoveryTag:  tag,
	}

	tags = map[string]*string{}

	if b := MustDiscoverBasedOnTags(tags); !b {
		t.Fatalf("Expected %v but got %v", false, b)
	}

	// AutoDiscoveryMode: None
	// tag: true
	CurrentConfig = &PrometheusAzureExporterConfig{
		AutoDiscoveryMode: "None",
		AutoDiscoveryTag:  tag,
	}

	tags = map[string]*string{
		tag: &strue,
	}

	if b := MustDiscoverBasedOnTags(tags); !b {
		t.Fatalf("Expected %v but got %v", true, b)
	}

	// AutoDiscoveryMode: None
	// tag: false
	CurrentConfig = &PrometheusAzureExporterConfig{
		AutoDiscoveryMode: "None",
		AutoDiscoveryTag:  tag,
	}

	tags = map[string]*string{
		tag: &sfalse,
	}

	if b := MustDiscoverBasedOnTags(tags); b {
		t.Fatalf("Expected %v but got %v", false, b)
	}

	// AutoDiscoveryMode: None
	// tag: nil
	CurrentConfig = &PrometheusAzureExporterConfig{
		AutoDiscoveryMode: "None",
		AutoDiscoveryTag:  tag,
	}

	tags = map[string]*string{}

	if b := MustDiscoverBasedOnTags(tags); b {
		t.Fatalf("Expected %v but got %v", false, b)
	}
}
