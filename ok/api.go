package ok

import (
	"fmt"
	"log"

	"github.com/kirillDanshin/myutils"
	"github.com/kirillDanshin/ok-mysql/defaults"
	"github.com/valyala/bytebufferpool"

	"net"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// Config for ok-mysql
type Config struct {
	// Address to read from
	Address string

	// SnapshotLength for pcap packet capture
	SnapshotLength int32
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

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// fmt.Printf("%s\n\n\n", string(packet.Data()))
		go func(p gopacket.Packet) {
			processPacket(p)
		}(packet)

	}

	close(syncPrint)

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

func processPacket(packet gopacket.Packet) {
	bb := bytebufferpool.Get()
	defer bytebufferpool.Put(bb)
	// Let's see if the packet is an ethernet packet
	ethernetLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethernetLayer != nil {
		fmt.Fprintln(bb, "Ethernet layer detected.")
		ethernetPacket, _ := ethernetLayer.(*layers.Ethernet)
		fmt.Fprintln(bb, "Source MAC: ", ethernetPacket.SrcMAC)
		fmt.Fprintln(bb, "Destination MAC: ", ethernetPacket.DstMAC)
		// Ethernet type is typically IPv4 but could be ARP or other
		fmt.Fprintln(bb, "Ethernet type: ", ethernetPacket.EthernetType)
		fmt.Fprintln(bb)
	}

	// Let's see if the packet is IP (even though the ether type told us)
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	if ipLayer != nil {
		fmt.Fprintln(bb, "IPv4 layer detected.")
		ip, _ := ipLayer.(*layers.IPv4)

		// IP layer variables:
		// Version (Either 4 or 6)
		// IHL (IP Header Length in 32-bit words)
		// TOS, Length, Id, Flags, FragOffset, TTL, Protocol (TCP?),
		// Checksum, SrcIP, DstIP
		fmt.Fprintf(bb, "From %s to %s\n", ip.SrcIP, ip.DstIP)
		fmt.Fprintln(bb, "Protocol: ", ip.Protocol)
		fmt.Fprintln(bb)
	}

	// Let's see if the packet is TCP
	tcpLayer := packet.Layer(layers.LayerTypeTCP)
	if tcpLayer != nil {
		fmt.Fprintln(bb, "TCP layer detected.")
		tcp, _ := tcpLayer.(*layers.TCP)

		// TCP layer variables:
		// SrcPort, DstPort, Seq, Ack, DataOffset, Window, Checksum, Urgent
		// Bool flags: FIN, SYN, RST, PSH, ACK, URG, ECE, CWR, NS
		fmt.Fprintf(bb, "From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
		fmt.Fprintln(bb, "Sequence number: ", tcp.Seq)
		fmt.Fprintln(bb)
	}

	// Iterate over all layers, printing out each layer type
	fmt.Fprintln(bb, "All packet layers:")
	for _, layer := range packet.Layers() {
		fmt.Fprintln(bb, "- ", layer.LayerType())
	}

	// When iterating through packet.Layers() above,
	// if it lists Payload layer then that is the same as
	// this applicationLayer. applicationLayer contains the payload
	applicationLayer := packet.ApplicationLayer()
	if applicationLayer != nil {
		fmt.Fprintln(bb, "Application layer/Payload found.")
		// fmt.Printf("%s\n", applicationLayer.Payload())
		parsePacket(applicationLayer)
	}

	// Check for errors
	if err := packet.ErrorLayer(); err != nil {
		fmt.Fprintln(bb, "Error decoding some part of the packet:", err)
	}

	syncPrint <- fmt.Sprintf("%s", bb.B)
	syncPrint <- fmt.Sprintf("%v", packet.Data())
}

var i int

func parsePacket(app gopacket.ApplicationLayer) error {
	// fmt.Printf("contents: %+#v\n\n", app.LayerContents())
	// fmt.Printf("layer payload: %+#v\n\n", app.LayerPayload())
	// fmt.Printf("payload: %+#v\n\n", app.Payload())
	// fmt.Printf("%s\n\n\n\n\n", app.LayerType())
	// i++
	// if i == 2 {
	// 	os.Exit(0)
	// }

	return nil
}

var (
	syncPrint = make(chan string, 128)
)

func syncPrinter(c chan string) {
	for s := range c {
		fmt.Println(s)
		fmt.Print("\n\n\n\n")
	}
}
