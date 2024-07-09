package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "sync"
    "time"
    

    "github.com/Dyst0rti0n/gonetdiag/internal/bandwidth"
    "github.com/Dyst0rti0n/gonetdiag/internal/latency"
    "github.com/Dyst0rti0n/gonetdiag/internal/packetloss"
    "github.com/Dyst0rti0n/gonetdiag/internal/ping"
    "github.com/Dyst0rti0n/gonetdiag/internal/report"
    "github.com/Dyst0rti0n/gonetdiag/internal/traceroute"
    "github.com/fatih/color"
    "github.com/spf13/cobra"
)

func main() {
    var rootCmd = &cobra.Command{Use: "gonetdiag"}

    rootCmd.AddCommand(&cobra.Command{
        Use:   "ping [target]",
        Short: "Ping a target",
        Args:  cobra.MinimumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            target := args[0]
            count, _ := cmd.Flags().GetInt("count")
            timeout, _ := cmd.Flags().GetDuration("timeout")

            result, err := ping.Ping(target, count, timeout)
            if err != nil {
                color.Red("Ping error: %v", err)
                return
            }
            color.Cyan("Ping Result:\n%s", result)
        },
    })

    rootCmd.PersistentFlags().IntP("count", "c", 4, "Number of pings")
    rootCmd.PersistentFlags().DurationP("timeout", "t", 5*time.Second, "Timeout for each ping")

    rootCmd.AddCommand(&cobra.Command{
        Use:   "traceroute [target]",
        Short: "Traceroute to a target",
        Args:  cobra.MinimumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            target := args[0]
            result, err := traceroute.TraceRoute(target)
            if err != nil {
                color.Red("Traceroute error: %v", err)
                return
            }
            color.Cyan("Traceroute Result:\n%s", result)
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "bandwidth [target]",
        Short: "Measure bandwidth to a target",
        Args:  cobra.MinimumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            target := args[0]
            protocol, _ := cmd.Flags().GetString("protocol")
            if protocol == "" {
                protocol = "http"
            }

            result, err := bandwidth.MeasureDownloadBandwidth(target, protocol)
            if err != nil {
                color.Red("Bandwidth measurement error: %v", err)
                return
            }
            color.Cyan("Bandwidth Result:\n%s", result)
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "latency [target]",
        Short: "Analyze latency to a target",
        Args:  cobra.MinimumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            target := args[0]
            count, _ := cmd.Flags().GetInt("count")
            timeout, _ := cmd.Flags().GetDuration("timeout")

            result, err := latency.AnalyzeLatency(target, count, timeout)
            if err != nil {
                color.Red("Latency analysis error: %v", err)
                return
            }
            color.Cyan("Latency Result:\n%s", result)
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "packetloss [target]",
        Short: "Detect packet loss to a target",
        Args:  cobra.MinimumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            target := args[0]
            count, _ := cmd.Flags().GetInt("count")
            timeout, _ := cmd.Flags().GetDuration("timeout")

            result, err := packetloss.DetectPacketLoss(target, count, timeout)
            if err != nil {
                color.Red("Packet loss detection error: %v", err)
                return
            }
            color.Cyan("Packet Loss Result:\n%s", result)
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "report [target]",
        Short: "Generate a network diagnostic report for a target",
        Args:  cobra.MinimumNArgs(1),
        Run: func(cmd *cobra.Command, args []string) {
            target := args[0]

            var wg sync.WaitGroup
            wg.Add(5)

            var pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult string
            var pingErr, traceErr, bandwidthErr, latencyErr, packetLossErr error

            go func() {
                defer wg.Done()
                pingResult, pingErr = ping.Ping(target, 4, 5*time.Second)
            }()

            go func() {
                defer wg.Done()
                traceResult, traceErr = traceroute.TraceRoute(target)
            }()

            go func() {
                defer wg.Done()
                bandwidthResult, bandwidthErr = bandwidth.MeasureUploadBandwidth(target)
            }()

            go func() {
                defer wg.Done()
                latencyResult, latencyErr = latency.AnalyzeLatency(target, 10, 15*time.Second)
            }()

            go func() {
                defer wg.Done()
                packetLossResult, packetLossErr = packetloss.DetectPacketLoss(target, 20, 25*time.Second)
            }()

            wg.Wait()

            if pingErr != nil || traceErr != nil || bandwidthErr != nil || latencyErr != nil || packetLossErr != nil {
                color.Red("Error in generating report: pingErr=%v, traceErr=%v, bandwidthErr=%v, latencyErr=%v, packetLossErr=%v",
                    pingErr, traceErr, bandwidthErr, latencyErr, packetLossErr)
                return
            }

            err := report.GenerateReport(target, pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult)
            if err != nil {
                color.Red("Report generation error: %v", err)
                return
            }
            color.Green("Report generated successfully!")
        },
    })

    rootCmd.AddCommand(&cobra.Command{
        Use:   "interactive",
        Short: "Interactive mode for easier input",
        Run: func(cmd *cobra.Command, args []string) {
            reader := bufio.NewReader(os.Stdin)

            fmt.Print("Enter the target: ")
            target, _ := reader.ReadString('\n')
            target = strings.TrimSpace(target)

            fmt.Print("Enter the test to run (ping, traceroute, bandwidth, latency, packetloss, report): ")
            test, _ := reader.ReadString('\n')
            test = strings.TrimSpace(test)

            switch test {
            case "ping":
                fmt.Print("Enter the number of pings (default 4): ")
                var count int
                fmt.Scanln(&count)
                if count == 0 {
                    count = 4
                }
                result, err := ping.Ping(target, count, 5*time.Second)
                if err != nil {
                    color.Red("Ping error: %v", err)
                } else {
                    color.Cyan("Ping Result:\n%s", result)
                }
            case "traceroute":
                result, err := traceroute.TraceRoute(target)
                if err != nil {
                    color.Red("Traceroute error: %v", err)
                } else {
                    color.Cyan("Traceroute Result:\n%s", result)
                }
            case "bandwidth":
                fmt.Print("Enter the protocol (http or https): ")
                protocol, _ := reader.ReadString('\n')
                protocol = strings.TrimSpace(protocol)
                result, err := bandwidth.MeasureDownloadBandwidth(target, protocol)
                if err != nil {
                    color.Red("Bandwidth measurement error: %v", err)
                } else {
                    color.Cyan("Bandwidth Result:\n%s", result)
                }
            case "latency":
                fmt.Print("Enter the number of pings (default 4): ")
                var count int
                fmt.Scanln(&count)
                if count == 0 {
                    count = 4
                }
                result, err := latency.AnalyzeLatency(target, count, 5*time.Second)
                if err != nil {
                    color.Red("Latency analysis error: %v", err)
                } else {
                    color.Cyan("Latency Result:\n%s", result)
                }
            case "packetloss":
                fmt.Print("Enter the number of pings (default 4): ")
                var count int
                fmt.Scanln(&count)
                if count == 0 {
                    count = 4
                }
                result, err := packetloss.DetectPacketLoss(target, count, 5*time.Second)
                if err != nil {
                    color.Red("Packet loss detection error: %v", err)
                } else {
                    color.Cyan("Packet Loss Result:\n%s", result)
                }
            case "report":
                var wg sync.WaitGroup
                wg.Add(5)

                var pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult string
                var pingErr, traceErr, bandwidthErr, latencyErr, packetLossErr error

                go func() {
                    defer wg.Done()
                    pingResult, pingErr = ping.Ping(target, 4, 5*time.Second)
                }()

                go func() {
                    defer wg.Done()
                    traceResult, traceErr = traceroute.TraceRoute(target)
                }()

                go func() {
                    defer wg.Done()
                    bandwidthResult, bandwidthErr = bandwidth.MeasureUploadBandwidth(target)
                }()

                go func() {
                    defer wg.Done()
                    latencyResult, latencyErr = latency.AnalyzeLatency(target, 10, 15*time.Second)
                }()

                go func() {
                    defer wg.Done()
                    packetLossResult, packetLossErr = packetloss.DetectPacketLoss(target, 20, 25*time.Second)
                }()

                wg.Wait()

                if pingErr != nil || traceErr != nil || bandwidthErr != nil || latencyErr != nil || packetLossErr != nil {
                    color.Red("Error in generating report: pingErr=%v, traceErr=%v, bandwidthErr=%v, latencyErr=%v, packetLossErr=%v",
                        pingErr, traceErr, bandwidthErr, latencyErr, packetLossErr)
                } else {
                    err := report.GenerateReport(target, pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult)
                    if err != nil {
                        color.Red("Report generation error: %v", err)
                    } else {
                        color.Green("Report generated successfully!")
                    }
                }
            default:
                color.Red("Unknown test: %s", test)
            }
        },
    })

    if err := rootCmd.Execute(); err != nil {
        color.Red("CLI error: %v", err)
    }
}
