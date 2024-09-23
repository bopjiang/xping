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
	var totalTCPTime, totalDNSTime time.Duration
	var minTCPTime, maxTCPTime, minDNSTime, maxDNSTime time.Duration

	for i := 0; i < count; i++ {
		if i > 0 {
			time.Sleep(interval)
		}

		// DNS resolution timing
		dnsStart := time.Now()
		ips, err := net.LookupIP(host)
		dnsDuration := time.Since(dnsStart)

		if err != nil {
			fmt.Printf("%d: Could not resolve host %s: %v\n", i+1, host, err)
			failed++
			continue
		}

		ip := ips[0]

		// TCP connection timing
		address := fmt.Sprintf("%s:%d", ip, port)
		tcpStart := time.Now()
		conn, err := net.DialTimeout("tcp", address, timeout)
		tcpDuration := time.Since(tcpStart)

		if err != nil {
			fmt.Printf("%d: Failed - DNS: %v, TCP: failed (%v)\n", i+1, dnsDuration, err)
			failed++
		} else {
			fmt.Printf("%d: Success - DNS: %v, TCP: %v\n", i+1, dnsDuration, tcpDuration)
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
		avgTCPTime := totalTCPTime / time.Duration(successful)
		avgDNSTime := totalDNSTime / time.Duration(count)
		
		fmt.Printf("\nTCP Connection Statistics:\n")
		fmt.Printf("min/avg/max = %v/%v/%v\n", minTCPTime, avgTCPTime, maxTCPTime)
		
		fmt.Printf("\nDNS Resolution Statistics:\n")
		fmt.Printf("min/avg/max = %v/%v/%v\n", minDNSTime, avgDNSTime, maxDNSTime)
	}
}

// Other helper functions
// ...
