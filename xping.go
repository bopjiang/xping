package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var (
	count    int
	interval time.Duration
	timeout  time.Duration
	size     int
	protocol string
	port     int
)

func init() {
	flag.IntVar(&count, "c", 4, "number of ping packets to send")
	flag.DurationVar(&interval, "i", time.Second, "interval between pings")
	flag.DurationVar(&timeout, "W", time.Second, "timeout for each reply")
	flag.IntVar(&size, "s", 56, "size of ping packet to send")
	flag.StringVar(&protocol, "t", "icmp", "protocol to use (icmp or tcp)")
	flag.IntVar(&port, "p", 80, "TCP port")
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Please provide a hostname or IP address")
		os.Exit(1)
	}

	host := flag.Arg(0)

	// DNS resolution timing
	dnsStart := time.Now()
	ips, err := net.LookupIP(host)
	dnsDuration := time.Since(dnsStart)

	if err != nil {
		fmt.Printf("Could not resolve host %s: %v\n", host, err)
		os.Exit(1)
	}

	ip := ips[0]
	portString := fmt.Sprintf(":%d", port)
	fmt.Printf("Pinging %s (%s), protocol: %s%s:\n", host, ip, protocol, portString)
	fmt.Printf("DNS resolution time: %v\n", dnsDuration)

	if protocol == "icmp" {
		pingICMP(ip.String())
	} else if protocol == "tcp" {
		pingTCP(ip.String(), port)
	} else {
		fmt.Println("Unsupported protocol. Use 'icmp' or 'tcp'.")
		os.Exit(1)
	}
}

func pingICMP(ip string) {
	// ICMP ping implementation
	// ...
}

func pingTCP(ip string) {
	// TCP ping implementation
	// ...
}

// Other helper functions
// ...
