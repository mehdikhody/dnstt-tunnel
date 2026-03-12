package packet

type Packet struct {
	StreamID uint8
	PacketID uint16
	Domain   string
	Payload  []byte
}
