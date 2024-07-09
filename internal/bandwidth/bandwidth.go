package bandwidth

import (
    "fmt"
    "net"
    "time"
    "strings"
)

func MeasureBandwidth(target string) (string, error) {
    if !strings.Contains(target, ":") {
        target = fmt.Sprintf("%s:80", target) // Default to port 80 if no port is specified
    }

    conn, err := net.DialTimeout("tcp", target, 5*time.Second)
    if err != nil {
        return "", fmt.Errorf("failed to dial target: %w", err)
    }
    defer conn.Close()

    start := time.Now()
    data := make([]byte, 1024)
    var totalBytes int

    for time.Since(start) < time.Second {
        n, err := conn.Write(data)
        if err != nil {
            return "", err
        }
        totalBytes += n
    }

    bandwidth := float64(totalBytes) / time.Since(start).Seconds() / 1024 / 1024
    return fmt.Sprintf("Measured bandwidth to %s: %.2f Mbps", target, bandwidth), nil
}
