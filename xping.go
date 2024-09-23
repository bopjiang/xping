package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

var (
	count    int
	interval time.Duration
	timeout  time.Duration
	port     int
	verbose  bool
	host     string
)

func init() {
	flag.IntVar(&count, "c", 4, "number of ping packets to send")
	flag.DurationVar(&interval, "i", time.Second, "interval between pings")
	flag.DurationVar(&timeout, "W", time.Second, "timeout for each reply")
	flag.StringVar(&host, "h", "", "hostname or IP address to ping")
	flag.IntVar(&port, "p", 80, "TCP port to connect to")
	flag.BoolVar(&verbose, "v", false, "verbose output")
}

func main() {
	flag.Parse()

	// Initial DNS lookup
	ips, err := net.LookupIP(host)
	if err != nil {
		fmt.Printf("Could not resolve host %s: %v\n", host, err)
		os.Exit(1)
	}
	ip := ips[0]

	fmt.Printf("PING %s (%s) on TCP port %d\n", host, ip, port)
	pingTCP(host, ip.String(), port)
}

func pingTCP(host, ip string, port int) {
	var successful, failed int
	var totalTCPTime, totalDNSTime, totalTime float64
	var minTCPTime, maxTCPTime, minDNSTime, maxDNSTime, minTotalTime, maxTotalTime float64

	for i := 0; i < count; i++ {
		if i > 0 {
			time.Sleep(interval)
		}

		// DNS resolution timing
		dnsStart := time.Now()
		ips, err := net.LookupIP(host)
		dnsDuration := time.Since(dnsStart).Seconds() * 1000 // Convert to milliseconds

		if err != nil {
			fmt.Printf("From %s: tcp_seq=%d DNS resolution failed: %v\n", ip, i+1, err)
			failed++
			continue
		}

		resolvedIP := ips[0]

		// TCP connection timing
		address := fmt.Sprintf("%s:%d", resolvedIP, port)
		tcpStart := time.Now()
		conn, err := net.DialTimeout("tcp", address, timeout)
		tcpDuration := time.Since(tcpStart).Seconds() * 1000 // Convert to milliseconds

		totalDuration := dnsDuration + tcpDuration

		if err != nil {
			fmt.Printf("From %s: tcp_seq=%d Port %d closed (%.3f ms)\n", resolvedIP, i+1, port, totalDuration)
			failed++
		} else {
			fmt.Printf("From %s: tcp_seq=%d Port %d open time=%.3f ms\n", resolvedIP, i+1, port, totalDuration)
			conn.Close()
			successful++

			// Update TCP statistics
			totalTCPTime += tcpDuration
			if minTCPTime == 0 || tcpDuration < minTCPTime {
				minTCPTime = tcpDuration
			}
			if tcpDuration > maxTCPTime {
				maxTCPTime = tcpDuration
			}

			// Update Total time statistics
			totalTime += totalDuration
			if minTotalTime == 0 || totalDuration < minTotalTime {
				minTotalTime = totalDuration
			}
			if totalDuration > maxTotalTime {
				maxTotalTime = totalDuration
			}
		}

		// Update DNS statistics
		totalDNSTime += dnsDuration
		if minDNSTime == 0 || dnsDuration < minDNSTime {
			minDNSTime = dnsDuration
		}
		if dnsDuration > maxDNSTime {
			maxDNSTime = dnsDuration
		}
	}

	// Print summary statistics
	fmt.Printf("\n--- %s:%d ping statistics ---\n", host, port)
	fmt.Printf("%d packets transmitted, %d successful, %d failed\n", count, successful, failed)
	
	if successful > 0 {
		avgTCPTime := totalTCPTime / float64(successful)
		avgDNSTime := totalDNSTime / float64(count)
		avgTotalTime := totalTime / float64(successful)
		
		fmt.Printf("\nDNS Resolution Statistics:\n")
		fmt.Printf("min/avg/max = %.3f/%.3f/%.3f ms\n", minDNSTime, avgDNSTime, maxDNSTime)
		
		fmt.Printf("\nTCP Connection Statistics:\n")
		fmt.Printf("min/avg/max = %.3f/%.3f/%.3f ms\n", minTCPTime, avgTCPTime, maxTCPTime)
		
		fmt.Printf("\nTotal Time Statistics:\n")
		fmt.Printf("min/avg/max = %.3f/%.3f/%.3f ms\n", minTotalTime, avgTotalTime, maxTotalTime)
	}
}

// Other helper functions
// ...
