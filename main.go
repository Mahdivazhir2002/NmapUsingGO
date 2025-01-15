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
)

var openPorts []int
var mu sync.Mutex
const maxConcurrency = 100

func scanPort(ip string, port int, wg *sync.WaitGroup, sem chan struct{}) {
	defer wg.Done()
	sem <- struct{}{}
	defer func() { <-sem }()

	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", addr)
	if err == nil {
		conn.Close()
		mu.Lock()
		openPorts = append(openPorts, port)
		fmt.Printf("%d\t%s\n", port, "open")
		mu.Unlock()
	}
}


//stmp


func stmp(ip int){
	ipaddress :=fmt.Sprintf("127.0.0.1:%d", ip)
	fmt.Println(ipaddress)
	conn, err := net.Dial("tcp", ipaddress)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("ERROR reading from server:", err)
		os.Exit(1)
	}
	fmt.Println("Server banner STMP:", response)
	fmt.Fprintf(conn, "EHLO test.com\r\n")
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from server STMP:", err)
		os.Exit(1)
	}

	fmt.Println("Response after EHLO:", response)
	for {
		response, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Println(response)
		if response == "250 CHUNKING\r\n" {
			break
		}
		
	}
	


}


//ftp



func ftp(ip int) {

	ipaddress := fmt.Sprintf("127.0.0.1:%d", ip)
	fmt.Println("Connecting to:", ipaddress)

	
	conn, err := net.Dial("tcp", ipaddress)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("ERROT reading from server ftp:", err)
		os.Exit(1)
	}
	fmt.Println("Server banner FTP:", response)

	
	fmt.Fprintf(conn, "USER *****\r\n")
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from server ftp:", err)
		os.Exit(1)
	}
	fmt.Println("Response to USER command:", response)

	
	fmt.Fprintf(conn, "PASS ****\r\n")
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from server ftp:", err)
		os.Exit(1)
	}
	fmt.Println("Response to PASS command:", response)

	
	fmt.Fprintf(conn, "SYST\r\n")
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from server:", err)
		os.Exit(1)
	}
	fmt.Println("Response to SYST command:", response)

	fmt.Fprintf(conn, "FEAT\r\n")
	response, err = reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading from server:", err)
		os.Exit(1)
	}
	fmt.Println("Response to FEAT command:", response)


	for {
		response, err = reader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Println(response)
		if response == "211 End\r\n" {
			break
				}
	}
}



func main() {
	var portRangeFlag string
	flag.StringVar(&portRangeFlag, "p", "", "port range to scan")
	flag.Parse()

	if ip := flag.Arg(0); ip != "" {
		fmt.Printf("nmap scan report for %s : [%s]\n", ip, portRangeFlag)
		fmt.Printf("PORT\tSTATE\n")

		portRange := strings.Split(portRangeFlag, "-")
		if len(portRange) != 2 {
			fmt.Println("Invalid port range format")
			os.Exit(1)
		}

		startPort, err := strconv.Atoi(portRange[0])
		if err != nil {
			fmt.Println("Invalid start port")
			os.Exit(1)
		}

		endPort, err := strconv.Atoi(portRange[1])
		if err != nil {
			fmt.Println("Invalid end port")
			os.Exit(1)
		}

		if endPort > 65535 {
			fmt.Println("Error: end port must be between 0 and 65535")
			os.Exit(1)
		}

		var wg sync.WaitGroup
		sem := make(chan struct{}, maxConcurrency) 
		for port := startPort; port <= endPort; port++ {
			wg.Add(1)
			go scanPort(ip, port, &wg, sem)
		}
		wg.Wait()
		fmt.Println("Scan complete")
		// for i := 0; i < len(openPorts); i++ { 
		// 	ftp(openPorts[i])
		// 	stmp(openPorts[i])
		for i := 0; i < len(openPorts); i++ { 
			if openPorts[i]==21{
				ftp(openPorts[i])
			}else if openPorts[i] == 25 || openPorts[i] == 587 {
				stmp(openPorts[i])
			}
			
		} 
	} else {
		fmt.Println("Error IP address ")
		os.Exit(1)
	}
	
}
