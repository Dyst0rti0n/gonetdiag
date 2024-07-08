package traceroute

import (
    "fmt"
    "log"
    "strings"

    "github.com/aeden/traceroute"
)

func TraceRoute(target string) (string, error) {
    options := &traceroute.TracerouteOptions{}
    result, err := traceroute.Traceroute(target, options)
    if err != nil {
        return "", fmt.Errorf("failed to perform traceroute: %w", err)
    }

    var sb strings.Builder
    for _, hop := range result.Hops {
        sb.WriteString(hop.String() + "\n")
    }
    return sb.String(), nil
}
