package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	timeout         time.Duration
	parallelScans   int
	showClosedPorts bool
	showBanners     bool
	showOpenPorts   bool
)

func scanPort(host string, port int, wg *sync.WaitGroup) {
	defer wg.Done()

	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)

	if err != nil {
		if showClosedPorts {
			fmt.Printf("Port %d: Closed\n", port)
		}
		return
	}

	defer conn.Close()

	if showBanners {
		banner, err := retrieveBanner(conn)
		if err == nil {
			fmt.Printf("Port %d: Open - Banner: %s\n", port, banner)
		} else {
			fmt.Printf("Port %d: Open\n", port)
		}
	} else {
		if showOpenPorts {
			fmt.Printf("Port %d: Open\n", port)
		} else if showClosedPorts {
			fmt.Printf("Port %d: Closed\n", port)
		}
	}
}

func retrieveBanner(conn net.Conn) (string, error) {
	conn.SetReadDeadline(time.Now().Add(timeout))
	reader := bufio.NewReader(conn)
	banner, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(banner), nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <host> <start-port> <end-port>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	var timeoutString string
	flag.StringVar(&timeoutString, "timeout", "3s", "Connection timeout duration")
	flag.IntVar(&parallelScans, "parallel", 100, "Number of parallel scans")
	flag.BoolVar(&showClosedPorts, "show-closed", false, "Show closed ports in the output")
	flag.BoolVar(&showBanners, "show-banners", false, "Show banners from open ports")
	flag.BoolVar(&showOpenPorts, "show-open", false, "Show open ports in the output")

	flag.Parse()

	var err error
	timeout, err = time.ParseDuration(timeoutString)
	if err != nil {
		fmt.Println("Invalid timeout duration")
		os.Exit(1)
	}

	if flag.NArg() != 3 {
		flag.Usage()
		os.Exit(1)
	}

	host := flag.Arg(0)
	startPort, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		fmt.Println("Invalid start port")
		os.Exit(1)
	}

	endPort, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		fmt.Println("Invalid end port")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, parallelScans)

	for port := startPort; port <= endPort; port++ {
		semaphore <- struct{}{}
		wg.Add(1)
		go func(port int) {
			defer func() { <-semaphore }()
			scanPort(host, port, &wg)
		}(port)
	}

	wg.Wait()
}
