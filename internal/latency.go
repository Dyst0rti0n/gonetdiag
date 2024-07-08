package latency

import (
    "fmt"
    "time"

    "github.com/go-ping/ping"
)

func AnalyzeLatency(target string) (string, error) {
    pinger, err := ping.NewPinger(target)
    if err != nil {
        return "", fmt.Errorf("failed to create pinger: %w", err)
    }

    pinger.Count = 10
    pinger.Interval = time.Second
    pinger.Timeout = 15 * time.Second

    err = pinger.Run()
    if err != nil {
        return "", fmt.Errorf("failed to ping target: %w", err)
    }

    stats := pinger.Statistics()
    return fmt.Sprintf("Latency to %s: Avg %v, Max %v, Min %v, StdDev %v",
        target, stats.AvgRtt, stats.MaxRtt, stats.MinRtt, stats.StdDevRtt), nil
}
