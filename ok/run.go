package ok

import (
	"fmt"
	"log"

	"github.com/kirillDanshin/dlog"
	"github.com/kirillDanshin/myutils"
	"github.com/kirillDanshin/ok-mysql/defaults"

	"github.com/google/gopacket"
	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/pcap"
)

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
	dlog.Ln("Only capturing port 3306 packets.")

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

	return nil
}
