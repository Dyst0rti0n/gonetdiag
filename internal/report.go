package report

import (
    "fmt"
    "os"
    "time"
)

func GenerateReport(target, pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult string) error {
    report := fmt.Sprintf(
        "Network Diagnostic Report for %s\n\nTimestamp: %s\n\nPing Result:\n%s\n\nTraceroute Result:\n%s\n\nBandwidth Result:\n%s\n\nLatency Result:\n%s\n\nPacket Loss Result:\n%s\n",
        target, time.Now().Format(time.RFC3339), pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult)

    file, err := os.Create(fmt.Sprintf("%s_report.txt", target))
    if err != nil {
        return fmt.Errorf("failed to create report file: %w", err)
    }
    defer file.Close()

    _, err = file.WriteString(report)
    if err != nil {
        return fmt.Errorf("failed to write to report file: %w", err)
    }

    fmt.Println("Report generated successfully!")
    return nil
}
