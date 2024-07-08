package packetloss

import (
    "fmt"
    "time"

    "github.com/go-ping/ping"
)

func DetectPacketLoss(target string) (string, error) {
    pinger, err := ping.NewPinger(target)
    if err != nil {
        return "", fmt.Errorf("failed to create pinger: %w", err)
    }

    pinger.Count = 20
    pinger.Interval = time.Second
    pinger.Timeout = 25 * time.Second

    err = pinger.Run()
    if err != nil {
        return "", fmt.Errorf("failed to ping target: %w", err)
    }

    stats := pinger.Statistics()
    return fmt.Sprintf("Packet loss to %s: %v%% (Sent: %d, Received: %d, Lost: %d)",
        target, stats.PacketLoss, stats.PacketsSent, stats.PacketsRecv, stats.PacketsSent-stats.PacketsRecv), nil
}
