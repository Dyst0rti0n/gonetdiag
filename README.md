# GoNetDiag

GoNetDiag is a powerful network diagnostic tool written in Go. It provides various functionalities to analyze network performance, including ping, traceroute, bandwidth measurement, latency analysis, packet loss detection, and comprehensive network reports.

## Features

- **Ping**: Test the reachability of a host and measure round-trip time.
- **Traceroute**: Trace the route packets take to a network host.
- **Bandwidth**: Measure upload and download bandwidth to a target.
- **Latency**: Analyze the latency to a target.
- **Packet Loss**: Detect packet loss to a target.
- **Report**: Generate a comprehensive network diagnostic report.
- **Interactive Mode**: Easily input commands and parameters interactively.

## Installation

1. Ensure you have Go installed. If not, download and install it from [golang.org](https://golang.org/dl/).
2. Clone the repository:
   ```sh
   git clone https://github.com/Dyst0rti0n/gonetdiag.git
   ```
3. Navigate to the project directory:
   ```sh
   cd gonetdiag
   ```
4. Build the project:
   ```sh
   go build -o gonetdiag cmd/main.go
   ```

## Usage

Run the `gonetdiag` executable with the desired command and options.

### Ping

Ping a target to test reachability and measure round-trip time.
```sh
./gonetdiag ping [target] [flags]
```
Example:
```sh
./gonetdiag ping 8.8.8.8 --count 10 --timeout 2s
```

### Traceroute

Trace the route packets take to a network host.
```sh
./gonetdiag traceroute [target]
```
Example:
```sh
./gonetdiag traceroute 8.8.8.8
```

### Bandwidth

Measure upload or download bandwidth to a target.
```sh
./gonetdiag bandwidth [target] [flags]
```
Flags:
- `--protocol`: Specify the protocol (`http` or `https`).

Example:
```sh
./gonetdiag bandwidth 8.8.8.8 --protocol http
```

### Latency

Analyze the latency to a target.
```sh
./gonetdiag latency [target] [flags]
```
Example:
```sh
./gonetdiag latency 8.8.8.8 --count 10 --timeout 3s
```

### Packet Loss

Detect packet loss to a target.
```sh
./gonetdiag packetloss [target] [flags]
```
Example:
```sh
./gonetdiag packetloss 8.8.8.8 --count 20 --timeout 5s
```

### Report

Generate a comprehensive network diagnostic report for a target.
```sh
./gonetdiag report [target]
```
Example:
```sh
./gonetdiag report 8.8.8.8
```

### Interactive Mode

Run the tool in interactive mode for easier input.
```sh
./gonetdiag interactive
```

## Examples

- **Ping Example**:
  ```sh
  ./gonetdiag ping 8.8.8.8 --count 5 --timeout 1s
  ```

- **Traceroute Example**:
  ```sh
  ./gonetdiag traceroute 8.8.8.8
  ```

- **Bandwidth Measurement Example**:
  ```sh
  ./gonetdiag bandwidth 8.8.8.8 --protocol https
  ```

- **Latency Analysis Example**:
  ```sh
  ./gonetdiag latency 8.8.8.8 --count 10 --timeout 2s
  ```

- **Packet Loss Detection Example**:
  ```sh
  ./gonetdiag packetloss 8.8.8.8 --count 15 --timeout 3s
  ```

- **Generate Report Example**:
  ```sh
  ./gonetdiag report 8.8.8.8
  ```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
