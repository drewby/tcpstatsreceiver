package tcpstatsreceiver

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

const (
	TCP_UNKNOWN     = iota // 0
	TCP_ESTABLISHED        // 1
	TCP_SYN_SENT           // 2
	TCP_SYN_RECV           // 3
	TCP_FIN_WAIT1          // 4
	TCP_FIN_WAIT2          // 5
	TCP_TIME_WAIT          // 6
	TCP_CLOSE              // 7
	TCP_CLOSE_WAIT         // 8
	TCP_LAST_ACK           // 9
	TCP_LISTEN             // 10
	TCP_CLOSING            // 11
	TCP_MAX                // 12
)

type tcpStatsResult struct {
	LocalAddress string
	LocalPort    int64
	TxQueue      int64
	RxQueue      int64
	QueueLength  int64
}

type tcpStatsKey struct {
	LocalAddress string
	LocalPort    int64
}

type tcpStats struct {
	path       string
	portFilter map[int64]bool
	logger     *zap.Logger
}

func newTcpStats(path string, portFilter string, logger *zap.Logger) *tcpStats {
	return &tcpStats{
		path:       path,
		portFilter: parsePortFilter(portFilter, logger),
		logger:     logger,
	}
}

// Parse port filter string into a map of ports.
// Example: "80,443" -> map[80:true,443:true]
func parsePortFilter(portFilter string, logger *zap.Logger) map[int64]bool {
	portFilterMap := make(map[int64]bool)
	if portFilter != "" {
		portFilterStrings := strings.Split(portFilter, ",")
		for _, portFilterString := range portFilterStrings {
			port, err := strconv.ParseInt(portFilterString, 10, 64)
			if err != nil {
				logger.Error("Error parsing port filter", zap.String("port", portFilterString), zap.Error(err))
				continue
			}
			portFilterMap[port] = true
		}
	}
	return portFilterMap
}

// Get TCP stats from /proc/net/tcp (or other path)
func (t *tcpStats) get() ([]tcpStatsResult, error) {
	file, err := os.Open(t.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	statsMap, err := t.parseFile(file)
	if err != nil {
		return nil, err
	}

	stats := make([]tcpStatsResult, 0, len(statsMap))
	for _, stat := range statsMap {
		stats = append(stats, *stat)
	}
	return stats, nil
}

// Parse /proc/net/tcp (or other path) into a map of stats.
// We will use the following columns from /proc/net/tcp:
//
// 1: local_address:port
// 3: status
// 4: tx_queue:rx_queue
//
// For each local_address:port, we will sum the tx_queue and rx_queue to
// get the total queue size in bytes and then count the number of
// connections to get the queue length.
//
// Only connections with status >= TCP_ESTABLISHED (2) and < TCP_TIME_WAIT (6)
// will be counted.
func (t *tcpStats) parseFile(file *os.File) (map[tcpStatsKey]*tcpStatsResult, error) {
	statsMap := make(map[tcpStatsKey]*tcpStatsResult)
	scanner := bufio.NewScanner(file)
	firstLine := true
	for scanner.Scan() {
		if firstLine {
			firstLine = false
			continue // skip first line
		}

		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 9 {
			continue // skip invalid lines
		}

		// Parse local address and port, which is in parts[1]
		localAddress, localPort, err := convertHexIPPort(parts[1])
		if err != nil {
			t.logger.Error("Error parsing local address and port", zap.String("input", parts[1]), zap.Error(err))
			continue
		}

		// If portFilter is set, skip if localPort is not in the filter
		if len(t.portFilter) > 0 {
			if _, ok := t.portFilter[localPort]; !ok {
				continue
			}
		}

		// Parse TX and RX queues, which is in parts[4] split by a :
		queues := strings.Split(parts[4], ":")
		if len(queues) != 2 {
			t.logger.Error("Error parsing queues", zap.String("input", parts[4]))
			continue
		}

		txQueue, err := strconv.ParseInt(queues[0], 16, 64)
		if err != nil {
			t.logger.Error("Error parsing tx queue", zap.String("input", parts[4]), zap.Error(err))
			continue
		}

		rxQueue, err := strconv.ParseInt(queues[1], 16, 64)
		if err != nil {
			t.logger.Error("Error parsing rx queue", zap.String("input", parts[5]), zap.Error(err))
			continue
		}

		// Parse status, which is in parts[3]
		status, err := strconv.ParseInt(parts[3], 16, 64)
		if err != nil {
			t.logger.Error("Error parsing status", zap.String("input", parts[3]), zap.Error(err))
			continue
		}

		// TCP_LISTEN should not be counted as queued.
		queueLength := int64(0)
		if status != TCP_LISTEN {
			queueLength = 1
		}

		if status < TCP_TIME_WAIT || status == TCP_LISTEN {
			key := tcpStatsKey{
				LocalAddress: localAddress,
				LocalPort:    localPort,
			}
			if stat, ok := statsMap[key]; ok {
				stat.TxQueue += txQueue
				stat.RxQueue += rxQueue
				stat.QueueLength += queueLength
			} else {
				statsMap[key] = &tcpStatsResult{
					LocalAddress: localAddress,
					LocalPort:    localPort,
					TxQueue:      txQueue,
					RxQueue:      rxQueue,
					QueueLength:  queueLength,
				}
			}
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return statsMap, nil
}

// Convert hex IP:Port to human readable format.
func convertHexIPPort(hexIPPort string) (string, int64, error) {
	// Split IP and port.
	parts := strings.Split(hexIPPort, ":")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf("invalid IP:Port format")
	}

	// Convert IP.
	ipInt, err := strconv.ParseInt(parts[0], 16, 64)
	if err != nil {
		return "", 0, err
	}
	ip := fmt.Sprintf("%d.%d.%d.%d", byte(ipInt), byte(ipInt>>8), byte(ipInt>>16), byte(ipInt>>24))

	// Convert Port.
	portInt, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return "", 0, err
	}

	return ip, portInt, nil
}
