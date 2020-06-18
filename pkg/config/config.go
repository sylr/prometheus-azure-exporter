package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"time"

	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var (
	// ConfigFromFlagParser ...
	ConfigFromFlagParser *PrometheusAzureExporterConfig
	// CurrentConfig ...
	CurrentConfig *PrometheusAzureExporterConfig
	// AutoDiscoveryModeAll ...
	AutoDiscoveryModeAll = regexp.MustCompile(`^([Aa]ll)$`)
	// AutoDiscoveryModeTagged ...
	AutoDiscoveryModeTagged = regexp.MustCompile(`^([Tt]agged|[Nn]one)$`)
	// AutoDiscoveryTagTrue ...
	AutoDiscoveryTagTrue = regexp.MustCompile(`^([Tt]rue|[Yy]es)$`)
	// AutoDiscoveryTagFalse ...
	AutoDiscoveryTagFalse = regexp.MustCompile(`^([Ff]alse|[Nn]o)$`)
)

// PrometheusAzureExporterConfig ...
type PrometheusAzureExporterConfig struct {
	ConfigFile        string        `                          short:"f"   long:"config"               description:"Yaml config"`
	Verbose           []bool        `yaml:"verbose"            short:"v"   long:"verbose"              description:"Show verbose debug information"`
	JSONOutput        bool          `yaml:"json_output"        short:"j"   long:"json"                 description:"Use json format for output"`
	Version           bool          `                                      long:"version"              description:"Show version"`
	ListeningAddress  string        `yaml:"listening_address"  short:"a"   long:"address"              description:"Listening address" env:"LISTENING_ADDRESS" default:"0.0.0.0"`
	ListeningPort     uint          `yaml:"listening_port"     short:"p"   long:"port"                 description:"Listening port" env:"LISTENING_PORT" default:"9000"`
	UpdateInterval    time.Duration `yaml:"update_interval"    short:"i"   long:"interval"             description:"Number of seconds between metrics updates" default:"125s"`
	NoCache           bool          `yaml:"no-cache"                       long:"no-cache"             description:"Disable internal caching"`
	AutoDiscoveryMode string        `yaml:"autodiscovery_mode" short:"m"   long:"autodiscovery-mode"   description:"Which Azure resources should we pocess: All, Tagged" default:"All"`
	AutoDiscoveryTag  string        `yaml:"autodiscovery_tag"  short:"t"   long:"autodiscovery-tag"    description:"If discovery mode set to Tagged we process Azure Resources with this tag set to True, If discovery mode set to All, resources with this tag set to False will be discarded" default:"prometheus_io_azure_exporter_discover"`

	// Env vars used for Azure Authent, see
	// https://github.com/Azure/go-autorest/blob/v13.3.0/autorest/azure/auth/auth.go#L41-L51
	AzureTenantID            string `env:"AZURE_TENANT_ID"              description:"Azure tenant id"`
	AzureSubscriptionID      string `env:"AZURE_SUBSCRIPTION_ID"        description:"Azure subscription id"`
	AzureClientID            string `env:"AZURE_CLIENT_ID"              description:"Azure client id"`
	AzureClientSecret        string `env:"AZURE_CLIENT_SECRET"          description:"Azure client secret"`
	AzureCertificatePath     string `env:"AZURE_CERTIFICATE_PATH"       description:"Azure certificate path"`
	AzureCertificatePassword string `env:"AZURE_CERTIFICATE_PASSWORD"   description:"Azure certificate password"`
	AzureUsername            string `env:"AZURE_USERNAME"               description:"Azure username"`
	AzurePassword            string `env:"AZURE_PASSWORD"               description:"Azure password"`
	AzureEnvironment         string `env:"AZURE_ENVIRONMENT"            description:"Azure environment"`
	AzureADResource          string `env:"AZURE_AD_RESOURCE"            description:"Azure AD resource"`

	UpdateMetricsFunctions []UpdateMetricsFunctionConfig `yaml:"update_metrics_functions,omitempty"`
}

// UpdateMetricsFunctionConfig ...
type UpdateMetricsFunctionConfig struct {
	Name     string        `yaml:"name,omitempty"`
	Interval time.Duration `yaml:"interval,omitempty"`
}

// ParseConfigFile parses the config file defined by -f/--config
func ParseConfigFile() (*PrometheusAzureExporterConfig, error) {
	if ConfigFromFlagParser == nil || len(ConfigFromFlagParser.ConfigFile) == 0 {
		return ConfigFromFlagParser, nil
	}

	conf, err := LoadFile(ConfigFromFlagParser.ConfigFile)

	if err != nil {
		return nil, err
	}

	return conf, nil
}

// LoadFile parses the given YAML file into a Config.
func LoadFile(filename string) (*PrometheusAzureExporterConfig, error) {
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	cfg, err := parseYAML(content)

	if err != nil {
		return nil, fmt.Errorf("parsing YAML file %s: %v", filename, err)
	}

	return cfg, nil
}

// parseYAML parses the YAML input s into a Config.
func parseYAML(bytes []byte) (*PrometheusAzureExporterConfig, error) {
	cfg := *ConfigFromFlagParser
	err := yaml.UnmarshalStrict([]byte(bytes), &cfg)

	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ParseOptions loads config from cli arguments
func ParseOptions() {
	if ConfigFromFlagParser == nil {
		ConfigFromFlagParser = &PrometheusAzureExporterConfig{}
	}

	parser := flags.NewParser(ConfigFromFlagParser, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			log.Fatal(err)
			os.Exit(1)
		}
	}
}

// ValidateConfig returns a []error if config file contains configuration
// which does not make sens or cannot be applied
func ValidateConfig(conf *PrometheusAzureExporterConfig) []error {
	errs := make([]error, 0)

	if CurrentConfig != nil && conf.ListeningAddress != CurrentConfig.ListeningAddress {
		errs = append(errs, errors.New("config: cannot change listening address"))
	}

	if CurrentConfig != nil && conf.ListeningPort != CurrentConfig.ListeningPort {
		errs = append(errs, errors.New("config: cannot change listening port"))
	}

	switch {
	case AutoDiscoveryModeAll.MatchString(conf.AutoDiscoveryMode):
	case AutoDiscoveryModeTagged.MatchString(conf.AutoDiscoveryMode):
	default:
		str := fmt.Sprintf("config: `%s` is not a valid autodiscovery mode", conf.AutoDiscoveryMode)
		errs = append(errs, errors.New(str))
	}

	return errs
}

// MustDiscoverBasedOnTags tags an map of tags returns True if the object
// must be discovered based on autodiscovery mode.
func MustDiscoverBasedOnTags(tags map[string]*string) bool {
	if CurrentConfig != nil {
		tag := CurrentConfig.AutoDiscoveryTag
		switch {
		// All
		case AutoDiscoveryModeAll.MatchString(CurrentConfig.AutoDiscoveryMode):
			if val, ok := tags[tag]; ok {
				if AutoDiscoveryTagFalse.MatchString(*val) {
					return false
				}
			}
		// None
		case AutoDiscoveryModeTagged.MatchString(CurrentConfig.AutoDiscoveryMode):
			if val, ok := tags[tag]; ok {
				if AutoDiscoveryTagTrue.MatchString(*val) {
					return true
				}
			}

			return false
		}
	}

	return true
}
