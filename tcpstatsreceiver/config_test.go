package tcpstatsreceiver

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/drewby/tcpstatsreceiver/internal/metadata"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

func TestLoadConfig_Validate_Invalid(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Path = ""
	cfg.ScraperControllerSettings.CollectionInterval = 20
	assert.Error(t, cfg.Validate())
}

func TestLoadConfig_Validate_InvalidPortFilter(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Path = "/test"
	cfg.PortFilter = "1000,invalid"
	cfg.ScraperControllerSettings.CollectionInterval = 10
	assert.Error(t, cfg.Validate())
}

func TestLoadConfig_Validate_Valid(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Path = "/test"
	cfg.ScraperControllerSettings.CollectionInterval = 1
	assert.NoError(t, cfg.Validate())
}

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	assert.NoError(t, err)

	tests := []struct {
		id           component.ID
		expected     component.Config
		errorMessage string
	}{
		{
			id:           component.NewIDWithName("file", ""),
			errorMessage: "path cannot be empty",
		}, {
			id: component.NewIDWithName("file", "1"),
			expected: &Config{
				Path: "./test",
				ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
					CollectionInterval: time.Duration(60) * time.Second,
					InitialDelay:       time.Second,
				},
				MetricsBuilderConfig: metadata.MetricsBuilderConfig{
					Metrics: metadata.MetricsConfig{
						TCPQueueLength: metadata.MetricConfig{
							Enabled: true,
						},
						TCPQueueSize: metadata.MetricConfig{
							Enabled: true,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			if tt.errorMessage != "" {
				assert.EqualError(t, component.ValidateConfig(cfg), tt.errorMessage)
				return
			}

			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}
