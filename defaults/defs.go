// Package defaults provides default read-only settings
// for ok-mysql
package defaults

import "time"

const (
	// SnapLen is default pcap snaphost length
	SnapLen int32 = 65535

	// Timeout is default pcap timeout
	// Timeout = pcap.BlockForever
	Timeout = time.Second * 30

	// Promiscuous is default promiscuous mode status
	Promiscuous = false

	// Net is default network to use
	Net = "tcp4"
)
