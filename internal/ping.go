package ping

import (
    "fmt"
    "time"

    "github.com/go-ping/ping"
)

func Ping(target string, count int, timeout time.Duration) (string, error) {
    pinger, err := ping.NewPinger(target)
    if err != nil {
        return "", fmt.Errorf("failed to create pinger: %w", err)
    }

    pinger.Count = count
    pinger.Timeout = timeout

    err = pinger.Run()
    if err != nil {
        return "", fmt.Errorf("failed to ping target: %w", err)
    }

    stats := pinger.Statistics()
    return stats.String(), nil
}
