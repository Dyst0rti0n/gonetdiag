package main

import (
    "fmt"
    "github.com/Dyst0rti0n/gonetdiag/internal/ping"
    "github.com/Dyst0rti0n/gonetdiag/internal/traceroute"
)

func main() {
    fmt.Println("GoNetDiag - Advanced Network Diagnostic Tool")
    
    // Example usage
    pingResult := ping.Ping("8.8.8.8")
    fmt.Println("Ping Result:", pingResult)
    
    traceResult := traceroute.TraceRoute("8.8.8.8")
    fmt.Println("Traceroute Result:", traceResult)
}
