package ok

import (
	"net"

	"github.com/google/gopacket/ip4defrag"
)

type (
	// Config for ok-mysql
	Config struct {
		// Address to read from
		Address string

		// SnapshotLength for pcap packet capture
		SnapshotLength int32

		// Lazy?
		Lazy bool
	}

	// Instance instance
	Instance struct {
		Addr    *net.TCPAddr
		SnapLen int32
		// Lazy?
		Lazy bool

		defragger *ip4defrag.IPv4Defragmenter
		device    string
	}
)
