package tcpstatsreceiver

import (
	"path/filepath"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestGetTcpStats(t *testing.T) {
	expected := []tcpStatsResult{
		{
			LocalAddress: "192.168.1.16",
			LocalPort:    45334,
			TxQueue:      0,
			RxQueue:      12,
			QueueLength:  2,
		},
		{
			LocalAddress: "192.168.1.10",
			LocalPort:    39112,
			TxQueue:      10,
			RxQueue:      4,
			QueueLength:  2,
		},
	}

	f := filepath.Join("testdata", "tcp")
	s := newTcpStats(f, "", zap.NewNop())

	stats, err := s.get()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(stats, expected) {
		t.Errorf("Expected: %+v, but got %+v", expected, stats)
	}
}

func TestGetTcpStatsWithPortFilter(t *testing.T) {
	expected := []tcpStatsResult{
		{
			LocalAddress: "192.168.1.16",
			LocalPort:    45334,
			TxQueue:      0,
			RxQueue:      12,
			QueueLength:  2,
		},
	}

	f := filepath.Join("testdata", "tcp")
	s := newTcpStats(f, "45334", zap.NewNop())

	stats, err := s.get()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(stats, expected) {
		t.Errorf("Expected: %+v, but got %+v", expected, stats)
	}
}
