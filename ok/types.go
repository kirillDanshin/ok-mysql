package ok

import (
	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/ip4defrag"
	"github.com/ivpusic/grpool"
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

		// registry registry
		pool  *grpool.Pool
		queue chan gopacket.Packet
	}

	// map[ip:port][]packetInfo
	// registry map[string]*packets

	// packets struct {
	// 	info []packetInfo
	//
	// 	mtx *sync.Mutex
	// }
	//
	// packetSrc uint
	//
	// packetInfo struct {
	// 	From packetSrc
	// 	Time time.Time
	// 	ACK  bool
	// 	Ack  uint32
	// 	FIN  bool
	// }

)
