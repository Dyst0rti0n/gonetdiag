package latency

import (
    "fmt"
    "net"
    "time"

    "golang.org/x/net/ipv4"
    "github.com/Dyst0rti0n/gonetdiag/internal/icmp"
)

func AnalyzeLatency(target string, count int, timeout time.Duration) (string, error) {
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
    var minRTT, maxRTT, totalRTT time.Duration
    var packetsRecv int

    for i := 0; i < count; i++ {
        start := time.Now()
        if err := icmp.SendICMPRequest(pconn, destAddr, i); err != nil {
            return "", err
        }
        if err := icmp.ReceiveICMPReply(pconn, timeout); err != nil {
            continue // Count as a lost packet
        }
        RTT := time.Since(start)
        if packetsRecv == 0 || RTT < minRTT {
            minRTT = RTT
        }
        if RTT > maxRTT {
            maxRTT = RTT
        }
        totalRTT += RTT
        packetsRecv++
    }

    avgRTT := totalRTT / time.Duration(packetsRecv)
    return fmt.Sprintf("Latency to %s: Avg %v, Max %v, Min %v",
        target, avgRTT, maxRTT, minRTT), nil
}
