package tcpstatsreceiver

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/drewby/tcpstatsreceiver/internal/metadata"
)

type scraper struct {
	logger         *zap.Logger              // Logger to log events
	metricsBuilder *metadata.MetricsBuilder // MetricsBuilder to build metrics
	tcpStats       *tcpStats                // tcpStats to get stats from /proc/net/tcp
}

// newScraper is a constructor function which returns a new scraper instance
func newScraper(metricsBuilder *metadata.MetricsBuilder, path string, portFilter string, logger *zap.Logger) *scraper {
	return &scraper{
		logger:         logger,
		metricsBuilder: metricsBuilder,
		tcpStats:       newTcpStats(path, portFilter, logger),
	}
}

// scrape function that scrapes the files matching the pattern for metrics
func (s *scraper) scrape(ctx context.Context) (pmetric.Metrics, error) {
	s.logger.Debug("Scraping TCP stats at path", zap.String("path", s.tcpStats.path))

	// Call the `get` TcpStats function
	stats, err := s.tcpStats.get()
	if err != nil {
		return pmetric.NewMetrics(), err
	}

	s.logger.Debug("Found TCP stats", zap.Int("count", len(stats)))

	now := pcommon.NewTimestampFromTime(time.Now())
	for _, stat := range stats {
		s.metricsBuilder.RecordTCPQueueSizeDataPoint(now, stat.TxQueue, stat.LocalAddress, stat.LocalPort, "tx")
		s.metricsBuilder.RecordTCPQueueSizeDataPoint(now, stat.RxQueue, stat.LocalAddress, stat.LocalPort, "rx")
		s.metricsBuilder.RecordTCPQueueLengthDataPoint(now, stat.QueueLength, stat.LocalAddress, stat.LocalPort)
	}

	metrics := s.metricsBuilder.Emit()

	s.logger.Debug("Emitting TCP stats", zap.Int("MetricCount", metrics.MetricCount()), zap.Int("DataPointCount", metrics.DataPointCount()))

	return metrics, nil
}
