package traceroute

import (
    "github.com/aeden/traceroute"
    "log"
    "strings"
)

func TraceRoute(target string) string {
    options := &traceroute.TracerouteOptions{}
    result, err := traceroute.Traceroute(target, options)
    if err != nil {
        log.Fatalf("Failed to perform traceroute: %v", err)
    }

    var sb strings.Builder
    for _, hop := range result.Hops {
        sb.WriteString(hop.String() + "\n")
    }
    return sb.String()
}
