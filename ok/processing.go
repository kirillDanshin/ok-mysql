package ok

import (
	"fmt"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/kirillDanshin/dlog"
	"github.com/valyala/bytebufferpool"
)

func (i *Instance) processPacket(packet gopacket.Packet) {

	ip4Layer := packet.Layer(layers.LayerTypeIPv4)
	if ip4Layer == nil {
		return
	}
	ip4 := ip4Layer.(*layers.IPv4)
	l := ip4.Length
	newip4, err := i.defragger.DefragIPv4(ip4)
	if err != nil {
		dlog.Ln("Error while defragging", err)
	} else if newip4 == nil {
		dlog.Ln("Recieved a fragment")
		return
	}
	if newip4.Length != l {
		dlog.F("Decoding re-assembled packet: %s\n", newip4.NextLayerType())
		pb, ok := packet.(gopacket.PacketBuilder)
		if !ok {
			dlog.Ln("Error while getting packet builder: it's not a PacketBuilder")
		}
		nextDecoder := newip4.NextLayerType()
		nextDecoder.Decode(newip4.Payload, pb)
	}

	bb := bytebufferpool.Get()
	defer func() {
		syncPrint <- fmt.Sprintf("bb.B: %s", bb.B)
		syncPrint <- fmt.Sprintf("packet: %s", packet)
		syncPrint <- fmt.Sprintf("data: %s", packet.Data())
		bytebufferpool.Put(bb)
	}()
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

		// little hacky \x00 decoding
		s, _ := strconv.Unquote(fmt.Sprintf(`"%s"`, string(tcp.LayerPayload())))
		fmt.Fprintf(bb, "tcp.LayerContents(): %#+v", s)
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

}

func parsePacket(app gopacket.ApplicationLayer) error {
	syncPrint <- fmt.Sprintf("contents: %+#v\n\n", app.LayerContents())
	syncPrint <- fmt.Sprintf("layer payload: %+#v\n\n", app.LayerPayload())
	syncPrint <- fmt.Sprintf("payload: %+#v\n\n", app.Payload())
	syncPrint <- fmt.Sprintf("%s\n\n\n\n\n", app.LayerType())

	return nil
}
