package bandwidth

import (
    "fmt"
    "os/exec"
)

func MeasureBandwidth(target string) (string, error) {
    cmd := exec.Command("iperf3", "-c", target, "-f", "m")
    output, err := cmd.CombinedOutput()
    if err != nil {
        return "", fmt.Errorf("failed to measure bandwidth: %w", err)
    }

    return string(output), nil
}
