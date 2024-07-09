package report

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "os"
)

type Report struct {
    Target        string `json:"target"`
    PingResult    string `json:"ping_result"`
    TraceResult   string `json:"trace_result"`
    BandwidthResult string `json:"bandwidth_result"`
    LatencyResult string `json:"latency_result"`
    PacketLossResult string `json:"packet_loss_result"`
}

func GenerateReport(target, pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult string) error {
    report := Report{
        Target: target,
        PingResult: pingResult,
        TraceResult: traceResult,
        BandwidthResult: bandwidthResult,
        LatencyResult: latencyResult,
        PacketLossResult: packetLossResult,
    }

    // Save JSON report
    jsonFile, err := os.Create(fmt.Sprintf("%s_report.json", target))
    if err != nil {
        return fmt.Errorf("failed to create JSON report file: %w", err)
    }
    defer jsonFile.Close()

    jsonEncoder := json.NewEncoder(jsonFile)
    if err := jsonEncoder.Encode(report); err != nil {
        return fmt.Errorf("failed to encode JSON report: %w", err)
    }

    // Save CSV report
    csvFile, err := os.Create(fmt.Sprintf("%s_report.csv", target))
    if err != nil {
        return fmt.Errorf("failed to create CSV report file: %w", err)
    }
    defer csvFile.Close()

    csvWriter := csv.NewWriter(csvFile)
    defer csvWriter.Flush()

    if err := csvWriter.Write([]string{"Target", "PingResult", "TraceResult", "BandwidthResult", "LatencyResult", "PacketLossResult"}); err != nil {
        return fmt.Errorf("failed to write CSV header: %w", err)
    }
    if err := csvWriter.Write([]string{target, pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult}); err != nil {
        return fmt.Errorf("failed to write CSV record: %w", err)
    }

    return nil
}
