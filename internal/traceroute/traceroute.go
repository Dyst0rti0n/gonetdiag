package traceroute

import (
    "fmt"
    "net"
    "strings"
    "time"

    "golang.org/x/net/ipv4"
    "github.com/Dyst0rti0n/gonetdiag/internal/icmp"
)

func TraceRoute(target string) (string, error) {
    var sb strings.Builder

    destAddr, err := net.ResolveIPAddr("ip4", target)
    if err != nil {
        return "", fmt.Errorf("failed to resolve target: %w", err)
    }

    for ttl := 1; ttl <= 30; ttl++ {
        conn, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
        if err != nil {
            return "", fmt.Errorf("failed to listen on packet: %w", err)
        }
        defer conn.Close()

        pconn := ipv4.NewPacketConn(conn)
        if err := pconn.SetTTL(ttl); err != nil {
            return "", fmt.Errorf("failed to set TTL: %w", err)
        }

        start := time.Now()
        if err := icmp.SendICMPRequest(pconn, destAddr, ttl); err != nil {
            return "", fmt.Errorf("failed to send ICMP request: %w", err)
        }

        addr, err := receiveAddress(pconn, time.Second)
        if err != nil {
            sb.WriteString(fmt.Sprintf("%d: * * *\n", ttl))
            continue
        }

        host, err := net.LookupAddr(addr)
        if err != nil || len(host) == 0 {
            host = []string{addr}
        }

        RTT := time.Since(start)
        sb.WriteString(fmt.Sprintf("%d: %s, RTT = %v\n", ttl, host[0], RTT))
        if addr == destAddr.String() {
            break
        }
    }

    return sb.String(), nil
}

func receiveAddress(conn *ipv4.PacketConn, timeout time.Duration) (string, error) {
    conn.SetReadDeadline(time.Now().Add(timeout))
    buf := make([]byte, 512)
    _, cm, src, err := conn.ReadFrom(buf)
    if err != nil {
        return "", err
    }
    if cm == nil {
        return "", fmt.Errorf("control message is nil")
    }
    if cm.TTL == 0 {
        return "", fmt.Errorf("TTL expired")
    }
    return src.String(), nil
}
