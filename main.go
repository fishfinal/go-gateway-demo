// Copyright (c) 2026 fishfinal
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// Server encapsulates service information
type Server struct {
	hostname string
	localIP  string
	port     string
}

// NewServer creates a new server instance
func NewServer(port string) *Server {
	hostname, _ := os.Hostname()
	localIP := getLocalIP()
	return &Server{
		hostname: hostname,
		localIP:  localIP,
		port:     port,
	}
}

// getLocalIP retrieves the local non-loopback IPv4 address
func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "unknown"
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "unknown"
}

// setupRouter configures the routing
func (s *Server) setupRouter() *gin.Engine {
	r := gin.Default()

	// Health check endpoint (for Nginx/Keepalived)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":    "ok",
			"hostname":  s.hostname,
			"ip":        s.localIP,
			"timestamp": time.Now().Unix(),
		})
	})

	// Business endpoint: simulate task dispatching
	r.POST("/api/v1/task", func(c *gin.Context) {
		var req map[string]interface{}
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "invalid request"})
			return
		}

		c.JSON(200, gin.H{
			"code":    0,
			"message": "task accepted",
			"processed_by": gin.H{
				"hostname": s.hostname,
				"ip":       s.localIP,
			},
			"task": req,
		})
	})

	// Simulated slow endpoint: for timeout testing
	r.GET("/api/slow", func(c *gin.Context) {
		time.Sleep(5 * time.Second)
		c.JSON(200, gin.H{
			"message":      "slow response",
			"processed_by": s.hostname,
		})
	})

	// Simulated error endpoint: for failure detection testing
	r.GET("/api/error", func(c *gin.Context) {
		c.JSON(500, gin.H{
			"error":        "internal server error",
			"processed_by": s.hostname,
		})
	})

	// Debug endpoint: view service information
	r.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hostname": s.hostname,
			"ip":       s.localIP,
			"port":     s.port,
		})
	})

	return r
}

// Run starts the server
func (s *Server) Run() error {
	r := s.setupRouter()
	addr := ":" + s.port

	fmt.Printf("🚀 Server starting on %s\n", addr)
	fmt.Printf("📡 Hostname: %s\n", s.hostname)
	fmt.Printf("🌐 IP: %s\n", s.localIP)
	fmt.Printf("✅ Health check: http://localhost%s/health\n", addr)
	fmt.Printf("📋 Task API: POST http://localhost%s/api/v1/task\n", addr)

	return r.Run(addr)
}

func main() {
	// Read port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := NewServer(port)

	// Graceful shutdown handling
	go func() {
		if err := server.Run(); err != nil {
			fmt.Printf("Server error: %v\n", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n🛑 Shutting down server...")
}
