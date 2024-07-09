package packetloss

import (
    "fmt"
    "net"
    "time"

    "golang.org/x/net/ipv4"
    "github.com/Dyst0rti0n/gonetdiag/internal/icmp"
)

func DetectPacketLoss(target string, count int, timeout time.Duration) (string, error) {
    destAddr, err := net.ResolveIPAddr("ip4", target)
    if err != nil {
        return "", fmt.Errorf("failed to resolve target: %w", err)
    }

    conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
    if err != nil {
        return "", fmt.Errorf("failed to listen on packet: %w", err)
    }
    defer conn.Close()

    pconn := ipv4.NewPacketConn(conn)
    var packetsSent, packetsRecv int

    for i := 0; i < count; i++ {
        if err := icmp.SendICMPRequest(pconn, destAddr, i); err != nil {
            return "", err
        }
        if err := icmp.ReceiveICMPReply(pconn, timeout); err != nil {
            continue // Count as a lost packet
        }
        packetsRecv++
        packetsSent++
    }

    packetLoss := float64(packetsSent-packetsRecv) / float64(packetsSent) * 100
    return fmt.Sprintf("Packet loss to %s: %.2f%% (Sent: %d, Received: %d, Lost: %d)",
        target, packetLoss, packetsSent, packetsRecv, packetsSent-packetsRecv), nil
}
