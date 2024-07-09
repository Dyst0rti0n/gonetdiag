package icmp

import (
    "net"
	"time"

    "golang.org/x/net/ipv4"
)

func SendICMPRequest(conn *ipv4.PacketConn, destAddr *net.IPAddr, seq int) error {
    msg := make([]byte, 8)
    msg[0] = 8 // Echo request
    msg[1] = 0 // Code 0
    msg[2] = 0 // Checksum placeholder
    msg[3] = 0 // Checksum placeholder
    msg[6] = byte(seq >> 8)
    msg[7] = byte(seq & 0xff)

    csum := Checksum(msg)
    msg[2] = byte(csum >> 8)
    msg[3] = byte(csum & 0xff)

    _, err := conn.WriteTo(msg, nil, destAddr)
    return err
}

func ReceiveICMPReply(conn *ipv4.PacketConn, timeout time.Duration) error {
    conn.SetReadDeadline(time.Now().Add(timeout))
    reply := make([]byte, 20+8)
    _, _, _, err := conn.ReadFrom(reply)
    return err
}

func Checksum(data []byte) uint16 {
    var sum uint32
    for i := 0; i < len(data)-1; i += 2 {
        sum += uint32(data[i])<<8 | uint32(data[i+1])
    }
    if len(data)%2 == 1 {
        sum += uint32(data[len(data)-1]) << 8
    }
    for (sum >> 16) > 0 {
        sum = (sum & 0xffff) + (sum >> 16)
    }
    return ^uint16(sum)
}
