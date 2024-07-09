package bandwidth

import (
    "fmt"
    "io"
    "net"
    "net/http"
    "strings"
    "time"
)

func MeasureUploadBandwidth(target string) (string, error) {
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
    return fmt.Sprintf("Measured upload bandwidth to %s: %.2f Mbps", target, bandwidth), nil
}

func MeasureDownloadBandwidth(target, protocol string) (string, error) {
    if protocol != "http" && protocol != "https" {
        return "", fmt.Errorf("invalid protocol specified: %s", protocol)
    }

    if !strings.HasPrefix(target, "http://") && !strings.HasPrefix(target, "https://") {
        target = fmt.Sprintf("%s://%s", protocol, target)
    }

    client := &http.Client{
        Timeout: 10 * time.Second,
    }

    start := time.Now()
    resp, err := client.Get(target)
    if err != nil {
        return "", fmt.Errorf("failed to perform GET request: %w", err)
    }
    defer resp.Body.Close()

    var totalBytes int
    buffer := make([]byte, 32*1024)
    for {
        n, err := resp.Body.Read(buffer)
        totalBytes += n
        if err != nil {
            if err == io.EOF {
                break
            }
            return "", fmt.Errorf("error while reading response body: %w", err)
        }
    }

    bandwidth := float64(totalBytes) / time.Since(start).Seconds() / 1024 / 1024
    return fmt.Sprintf("Measured download bandwidth from %s: %.2f Mbps", target, bandwidth), nil
}
