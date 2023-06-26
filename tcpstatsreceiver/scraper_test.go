package tcpstatsreceiver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"

	"github.com/drewby/tcpstatsreceiver/internal/metadata"
)

func TestNewScraper(t *testing.T) {
	logger := zap.NewNop()
	metricsBuilder := &metadata.MetricsBuilder{}
	path := "/path/to/scrape"

	s := newScraper(metricsBuilder, path, "", logger)

	assert.NotNil(t, s)
	assert.Equal(t, s.tcpStats.path, path)
	assert.Equal(t, s.logger, logger)
	assert.Equal(t, s.metricsBuilder, metricsBuilder)
}

func TestScrape(t *testing.T) {
	logger := zaptest.NewLogger(t)
	metricsBuilder := metadata.NewMetricsBuilder(metadata.DefaultMetricsBuilderConfig(), receivertest.NewNopCreateSettings())
	path := "testdata/tcp"

	s := newScraper(metricsBuilder, path, "", logger)

	ctx := context.Background()
	metrics, err := s.scrape(ctx)

	assert.NotNil(t, metrics)
	assert.Nil(t, err)

	assert.Equal(t, 2, metrics.MetricCount())
}
