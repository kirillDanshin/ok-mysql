package ok

import (
	"fmt"
	"log"

	"github.com/kirillDanshin/dlog"
	"github.com/kirillDanshin/myutils"
	"github.com/kirillDanshin/ok-mysql/defaults"

	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/pcap"
)

// Config for ok-mysql
type Config struct {
	// Address to read from
	Address string

	// SnapshotLength for pcap packet capture
	SnapshotLength int32

	// Lazy?
	Lazy bool
}

// NewInstance is a constructor for Instance
func NewInstance(config *Config) (*Instance, error) {
	if config == nil {
		return nil, fmt.Errorf("Instance config required, nil provided")
	}
	return newInst(config)
}

// Instance instance
type Instance struct {
	Addr    *net.TCPAddr
	SnapLen int32
	// Lazy?
	Lazy bool

	defragger *ip4defrag.IPv4Defragmenter
}

// Run the instance
func (i *Instance) Run() error {
	defer myutils.CPUProf()()

	var (
		device      = "lo" //TODO <kirilldanshin> device founder
		snapLen     = i.SnapLen
		promiscuous = defaults.Promiscuous
		timeout     = defaults.Timeout
		err         error
		handle      *pcap.Handle

		port = i.Addr.Port
	)

	i.defragger = ip4defrag.NewIPv4Defragmenter()

	// Open device
	handle, err = pcap.OpenLive(device, snapLen, promiscuous, timeout)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	filter := fmt.Sprintf("port %d", port)
	err = handle.SetBPFFilter(filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Only capturing port 3306 packets.")

	go syncPrinter(syncPrint)

	pSrc := gopacket.NewPacketSource(
		handle,
		handle.LinkType(),
	)
	pSrc.Lazy = i.Lazy
	count := 0
	bytes := int64(0)
	// start := time.Now()
	for packet := range pSrc.Packets() {
		count++
		bytes += int64(len(packet.Data()))
		// fmt.Printf("%s\n\n\n", string(packet.Data()))

		i.processPacket(packet)
	}

	close(syncPrint)

	dlog.F("Processed %d packets (%d bytes)", count, bytes)

	// ifaces, err := net.Interfaces()
	// if err != nil {
	// 	return err
	// }
	//
	// var wg sync.WaitGroup
	// for _, iface := range ifaces {
	// 	wg.Add(1)
	// 	go func(iface net.Interface) {
	// 		defer wg.Done()
	// 		if err := scan(&iface); err != nil {
	// 			log.Printf("interface %v: %v", iface.Name, err)
	// 		}
	// 	}(iface)
	// }
	//
	// // wait for all interfaces' scan to complete.
	// wg.Wait()

	return nil
}
