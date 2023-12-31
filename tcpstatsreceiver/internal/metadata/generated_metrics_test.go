// Code generated by mdatagen. DO NOT EDIT.

package metadata

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type testConfigCollection int

const (
	testSetDefault testConfigCollection = iota
	testSetAll
	testSetNone
)

func TestMetricsBuilder(t *testing.T) {
	tests := []struct {
		name      string
		configSet testConfigCollection
	}{
		{
			name:      "default",
			configSet: testSetDefault,
		},
		{
			name:      "all_set",
			configSet: testSetAll,
		},
		{
			name:      "none_set",
			configSet: testSetNone,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			start := pcommon.Timestamp(1_000_000_000)
			ts := pcommon.Timestamp(1_000_001_000)
			observedZapCore, observedLogs := observer.New(zap.WarnLevel)
			settings := receivertest.NewNopCreateSettings()
			settings.Logger = zap.New(observedZapCore)
			mb := NewMetricsBuilder(loadMetricsBuilderConfig(t, test.name), settings, WithStartTime(start))

			expectedWarnings := 0
			assert.Equal(t, expectedWarnings, observedLogs.Len())

			defaultMetricsCount := 0
			allMetricsCount := 0

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordTCPQueueLengthDataPoint(ts, 1, "attr-val", 1)

			defaultMetricsCount++
			allMetricsCount++
			mb.RecordTCPQueueSizeDataPoint(ts, 1, "attr-val", 1, "attr-val")

			metrics := mb.Emit()

			if test.configSet == testSetNone {
				assert.Equal(t, 0, metrics.ResourceMetrics().Len())
				return
			}

			assert.Equal(t, 1, metrics.ResourceMetrics().Len())
			rm := metrics.ResourceMetrics().At(0)
			attrCount := 0
			enabledAttrCount := 0
			assert.Equal(t, enabledAttrCount, rm.Resource().Attributes().Len())
			assert.Equal(t, attrCount, 0)

			assert.Equal(t, 1, rm.ScopeMetrics().Len())
			ms := rm.ScopeMetrics().At(0).Metrics()
			if test.configSet == testSetDefault {
				assert.Equal(t, defaultMetricsCount, ms.Len())
			}
			if test.configSet == testSetAll {
				assert.Equal(t, allMetricsCount, ms.Len())
			}
			validatedMetrics := make(map[string]bool)
			for i := 0; i < ms.Len(); i++ {
				switch ms.At(i).Name() {
				case "tcp.queue.length":
					assert.False(t, validatedMetrics["tcp.queue.length"], "Found a duplicate in the metrics slice: tcp.queue.length")
					validatedMetrics["tcp.queue.length"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Indicates the total number of sockets currently present in the TCP transmit and receive queues.", ms.At(i).Description())
					assert.Equal(t, "sockets", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("local.address")
					assert.True(t, ok)
					assert.EqualValues(t, "attr-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("local.port")
					assert.True(t, ok)
					assert.EqualValues(t, 1, attrVal.Int())
				case "tcp.queue.size":
					assert.False(t, validatedMetrics["tcp.queue.size"], "Found a duplicate in the metrics slice: tcp.queue.size")
					validatedMetrics["tcp.queue.size"] = true
					assert.Equal(t, pmetric.MetricTypeGauge, ms.At(i).Type())
					assert.Equal(t, 1, ms.At(i).Gauge().DataPoints().Len())
					assert.Equal(t, "Represents the total data size in bytes in the TCP transmit (tx) and receive (rx) queues.", ms.At(i).Description())
					assert.Equal(t, "bytes", ms.At(i).Unit())
					dp := ms.At(i).Gauge().DataPoints().At(0)
					assert.Equal(t, start, dp.StartTimestamp())
					assert.Equal(t, ts, dp.Timestamp())
					assert.Equal(t, pmetric.NumberDataPointValueTypeInt, dp.ValueType())
					assert.Equal(t, int64(1), dp.IntValue())
					attrVal, ok := dp.Attributes().Get("local.address")
					assert.True(t, ok)
					assert.EqualValues(t, "attr-val", attrVal.Str())
					attrVal, ok = dp.Attributes().Get("local.port")
					assert.True(t, ok)
					assert.EqualValues(t, 1, attrVal.Int())
					attrVal, ok = dp.Attributes().Get("queue")
					assert.True(t, ok)
					assert.EqualValues(t, "attr-val", attrVal.Str())
				}
			}
		})
	}
}
