package ping

import (
    "github.com/go-ping/ping"
    "log"
)

func Ping(target string) string {
    pinger, err := ping.NewPinger(target)
    if err != nil {
        log.Fatalf("Failed to create pinger: %v", err)
    }

    pinger.Count = 3
    err = pinger.Run() // Blocks until finished.
    if err != nil {
        log.Fatalf("Failed to ping target: %v", err)
    }

    stats := pinger.Statistics()
    return stats.String()
}
