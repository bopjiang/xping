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
)

func init() {
	flag.IntVar(&count, "c", 4, "number of ping packets to send")
	flag.DurationVar(&interval, "i", time.Second, "interval between pings")
	flag.DurationVar(&timeout, "W", time.Second, "timeout for each reply")
	flag.IntVar(&port, "p", 80, "TCP port to connect to")
	flag.BoolVar(&verbose, "v", false, "verbose output")
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Please provide a hostname or IP address")
		os.Exit(1)
	}

	host := flag.Arg(0)
	pingTCP(host, port)
}

func pingTCP(host string, port int) {
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
			fmt.Printf("%d: Failed - DNS: %.3fms, TCP: N/A (DNS resolution failed: %v)\n", i+1, dnsDuration, err)
			failed++
			continue
		}

		ip := ips[0]

		// TCP connection timing
		address := fmt.Sprintf("%s:%d", ip, port)
		tcpStart := time.Now()
		conn, err := net.DialTimeout("tcp", address, timeout)
		tcpDuration := time.Since(tcpStart).Seconds() * 1000 // Convert to milliseconds

		totalDuration := dnsDuration + tcpDuration

		if err != nil {
			fmt.Printf("%d: Failed - DNS: %.3fms, TCP: failed (%.3fms, %v), Total: %.3fms\n", i+1, dnsDuration, tcpDuration, err, totalDuration)
			failed++
		} else {
			fmt.Printf("%d: Success - DNS: %.3fms, TCP: %.3fms, Total: %.3fms\n", i+1, dnsDuration, tcpDuration, totalDuration)
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
