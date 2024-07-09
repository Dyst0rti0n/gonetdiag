package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Dyst0rti0n/gonetdiag/internal/bandwidth"
	"github.com/Dyst0rti0n/gonetdiag/internal/latency"
	"github.com/Dyst0rti0n/gonetdiag/internal/packetloss"
	"github.com/Dyst0rti0n/gonetdiag/internal/ping"
	"github.com/Dyst0rti0n/gonetdiag/internal/report"
	"github.com/Dyst0rti0n/gonetdiag/internal/traceroute"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("web/templates/*")
	router.Static("/static", "./web/static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "layout.html", nil)
	})

	router.GET("/ping/:target", handlePing)
	router.GET("/traceroute/:target", handleTraceroute)
	router.GET("/bandwidth/:target", handleBandwidth)
	router.GET("/latency/:target", handleLatency)
	router.GET("/packetloss/:target", handlePacketLoss)
	router.GET("/report/:target", handleReport)

	router.GET("/ws", func(c *gin.Context) {
		handleWebSocket(c.Writer, c.Request)
	})

	return router
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			var msg string
			if err := websocket.Message.Receive(ws, &msg); err != nil {
				log.Println("Error reading message:", err)
				break
			}

			log.Println("Received message:", msg)

			var request map[string]string
			if err := json.Unmarshal([]byte(msg), &request); err != nil {
				log.Println("Error unmarshalling message:", err)
				break
			}

			action := request["action"]
			target := request["target"]

			switch action {
			case "ping":
				go handlePingWebSocket(ws, target)
			case "traceroute":
				go handleTracerouteWebSocket(ws, target)
			case "bandwidth":
				go handleBandwidthWebSocket(ws, target)
			case "latency":
				go handleLatencyWebSocket(ws, target)
			case "packetloss":
				go handlePacketLossWebSocket(ws, target)
			case "report":
				go handleReportWebSocket(ws, target)
			default:
				log.Println("Unknown action:", action)
			}
		}
	}).ServeHTTP(w, r)
}

func handlePingWebSocket(ws *websocket.Conn, target string) {
	count := 4
	timeout := 5 * time.Second
	result, err := ping.Ping(target, count, timeout)
	if err != nil {
		sendError(ws, err)
		return
	}
	sendResult(ws, "Ping Result", result)
}

func handleTracerouteWebSocket(ws *websocket.Conn, target string) {
	result, err := traceroute.TraceRoute(target)
	if err != nil {
		sendError(ws, err)
		return
	}
	sendResult(ws, "Traceroute Result", result)
}

func handleBandwidthWebSocket(ws *websocket.Conn, target string) {
	uploadResult, err := bandwidth.MeasureUploadBandwidth(target)
	if err != nil {
		sendError(ws, err)
		return
	}
	sendResult(ws, "Upload Bandwidth Result", uploadResult)

	downloadResult, err := bandwidth.MeasureDownloadBandwidth(target, "http")
	if err != nil {
		sendError(ws, err)
		return
	}
	sendResult(ws, "Download Bandwidth Result", downloadResult)
}

func handleLatencyWebSocket(ws *websocket.Conn, target string) {
	count := 4
	timeout := 5 * time.Second
	result, err := latency.AnalyzeLatency(target, count, timeout)
	if err != nil {
		sendError(ws, err)
		return
	}
	sendResult(ws, "Latency Result", result)
}

func handlePacketLossWebSocket(ws *websocket.Conn, target string) {
	count := 4
	timeout := 5 * time.Second
	result, err := packetloss.DetectPacketLoss(target, count, timeout)
	if err != nil {
		sendError(ws, err)
		return
	}
	sendResult(ws, "Packet Loss Result", result)
}

func handleReportWebSocket(ws *websocket.Conn, target string) {
	var pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult string
	var pingErr, traceErr, bandwidthErr, latencyErr, packetLossErr error

	pingResult, pingErr = ping.Ping(target, 4, 5*time.Second)
	traceResult, traceErr = traceroute.TraceRoute(target)
	uploadResult, uploadErr := bandwidth.MeasureUploadBandwidth(target)
	if uploadErr != nil {
		bandwidthErr = uploadErr
	} else {
		downloadResult, downloadErr := bandwidth.MeasureDownloadBandwidth(target, "http")
		if downloadErr != nil {
			bandwidthErr = downloadErr
		} else {
			bandwidthResult = fmt.Sprintf("%s\n%s", uploadResult, downloadResult)
		}
	}
	latencyResult, latencyErr = latency.AnalyzeLatency(target, 4, 5*time.Second)
	packetLossResult, packetLossErr = packetloss.DetectPacketLoss(target, 4, 5*time.Second)

	if pingErr != nil || traceErr != nil || bandwidthErr != nil || latencyErr != nil || packetLossErr != nil {
		sendError(ws, fmt.Errorf("Error in generating report: pingErr=%v, traceErr=%v, bandwidthErr=%v, latencyErr=%v, packetLossErr=%v",
			pingErr, traceErr, bandwidthErr, latencyErr, packetLossErr))
		return
	}

	err := report.GenerateReport(target, pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult)
	if err != nil {
		sendError(ws, err)
		return
	}

	sendResult(ws, "Report generated successfully!", "")
}

func sendError(ws *websocket.Conn, err error) {
	message := map[string]string{"error": err.Error()}
	jsonMessage, _ := json.Marshal(message)
	websocket.Message.Send(ws, string(jsonMessage))
}

func sendResult(ws *websocket.Conn, resultType, result string) {
	message := map[string]string{"type": resultType, "value": result}
	jsonMessage, _ := json.Marshal(message)
	websocket.Message.Send(ws, string(jsonMessage))
}

func handlePing(c *gin.Context) {
	target := c.Param("target")

	go func() {
		ws, err := websocket.Dial("ws://localhost:8080/ws", "", "http://localhost/")
		if err != nil {
			log.Println("Error connecting to websocket:", err)
			return
		}
		defer ws.Close()

		count := 4
		timeout := 5 * time.Second
		result, err := ping.Ping(target, count, timeout)
		if err != nil {
			websocket.Message.Send(ws, "Error: "+err.Error())
			return
		}

		websocket.Message.Send(ws, result)
	}()
}

func handleTraceroute(c *gin.Context) {
	target := c.Param("target")

	go func() {
		ws, err := websocket.Dial("ws://localhost:8080/ws", "", "http://localhost/")
		if err != nil {
			log.Println("Error connecting to websocket:", err)
			return
		}
		defer ws.Close()

		result, err := traceroute.TraceRoute(target)
		if err != nil {
			websocket.Message.Send(ws, "Error: "+err.Error())
			return
		}

		websocket.Message.Send(ws, result)
	}()
}

func handleBandwidth(c *gin.Context) {
	target := c.Param("target")

	go func() {
		ws, err := websocket.Dial("ws://localhost:8080/ws", "", "http://localhost/")
		if err != nil {
			log.Println("Error connecting to websocket:", err)
			return
		}
		defer ws.Close()

		uploadResult, err := bandwidth.MeasureUploadBandwidth(target)
		if err != nil {
			websocket.Message.Send(ws, "Upload Error: "+err.Error())
			return
		}

		websocket.Message.Send(ws, uploadResult)

		downloadResult, err := bandwidth.MeasureDownloadBandwidth(target, "http")
		if err != nil {
			websocket.Message.Send(ws, "Download Error: "+err.Error())
			return
		}

		websocket.Message.Send(ws, downloadResult)
	}()
}

func handleLatency(c *gin.Context) {
	target := c.Param("target")

	go func() {
		ws, err := websocket.Dial("ws://localhost:8080/ws", "", "http://localhost/")
		if err != nil {
			log.Println("Error connecting to websocket:", err)
			return
		}
		defer ws.Close()

		count := 4
		timeout := 5 * time.Second
		result, err := latency.AnalyzeLatency(target, count, timeout)
		if err != nil {
			websocket.Message.Send(ws, "Error: "+err.Error())
			return
		}

		websocket.Message.Send(ws, result)
	}()
}

func handlePacketLoss(c *gin.Context) {
	target := c.Param("target")

	go func() {
		ws, err := websocket.Dial("ws://localhost:8080/ws", "", "http://localhost/")
		if err != nil {
			log.Println("Error connecting to websocket:", err)
			return
		}
		defer ws.Close()

		count := 4
		timeout := 5 * time.Second
		result, err := packetloss.DetectPacketLoss(target, count, timeout)
		if err != nil {
			websocket.Message.Send(ws, "Error: "+err.Error())
			return
		}

		websocket.Message.Send(ws, result)
	}()
}

func handleReport(c *gin.Context) {
	target := c.Param("target")

	go func() {
		ws, err := websocket.Dial("ws://localhost:8080/ws", "", "http://localhost/")
		if err != nil {
			log.Println("Error connecting to websocket:", err)
			return
		}
		defer ws.Close()

		var pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult string
		var pingErr, traceErr, bandwidthErr, latencyErr, packetLossErr error

		pingResult, pingErr = ping.Ping(target, 4, 5*time.Second)
		traceResult, traceErr = traceroute.TraceRoute(target)
		uploadResult, uploadErr := bandwidth.MeasureUploadBandwidth(target)
		if uploadErr != nil {
			bandwidthErr = uploadErr
		} else {
			downloadResult, downloadErr := bandwidth.MeasureDownloadBandwidth(target, "http")
			if downloadErr != nil {
				bandwidthErr = downloadErr
			} else {
				bandwidthResult = fmt.Sprintf("%s\n%s", uploadResult, downloadResult)
			}
		}
		latencyResult, latencyErr = latency.AnalyzeLatency(target, 4, 5*time.Second)
		packetLossResult, packetLossErr = packetloss.DetectPacketLoss(target, 4, 5*time.Second)

		if pingErr != nil || traceErr != nil || bandwidthErr != nil || latencyErr != nil || packetLossErr != nil {
			websocket.Message.Send(ws, fmt.Sprintf("Error in generating report: pingErr=%v, traceErr=%v, bandwidthErr=%v, latencyErr=%v, packetLossErr=%v",
				pingErr, traceErr, bandwidthErr, latencyErr, packetLossErr))
			return
		}

		err = report.GenerateReport(target, pingResult, traceResult, bandwidthResult, latencyResult, packetLossResult)
		if err != nil {
			websocket.Message.Send(ws, "Report generation error: "+err.Error())
			return
		}

		websocket.Message.Send(ws, "Report generated successfully!")
	}()
}
