package tcpstatsreceiver

import (
	"errors"
	"strconv"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/drewby/tcpstatsreceiver/internal/metadata"
)

const (
	defaultPath = "/proc/net/tcp"
)

// Config defines the configuration for the TCP stats receiver.
type Config struct {
	Path                                    string                   `mapstructure:"path"`       // Path to the file to be scraped for metrics (default: /proc/net/tcp)
	PortFilter                              string                   `mapstructure:"portfilter"` // Comma-separated list of ports to filter on (default: "")
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"` // ScraperControllerSettings to configure scraping interval (default: 10s)
	metadata.MetricsBuilderConfig           `mapstructure:",squash"` // MetricsBuilderConfig to enable/disable specific metrics (default: all enabled)
}

func createDefaultConfig() component.Config {
	return &Config{
		Path:                      defaultPath,
		PortFilter:                "",
		ScraperControllerSettings: scraperhelper.NewDefaultScraperControllerSettings(metadata.Type),
		MetricsBuilderConfig:      metadata.DefaultMetricsBuilderConfig(),
	}
}

func (c Config) Validate() error {
	if c.Path == "" {
		return errors.New("path cannot be empty")
	}

	if c.PortFilter != "" {
		filters := strings.Split(c.PortFilter, ",")
		for _, filter := range filters {
			port, err := strconv.ParseInt(filter, 10, 64)
			if err != nil {
				return err
			}
			if port < 0 || port > 65535 {
				return errors.New("port filter must be between 0 and 65535")
			}
		}
	}

	return nil
}
