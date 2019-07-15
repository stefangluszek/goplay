package main

import (
	"flag"
	"fmt"
	"github.com/stefangluszek/goplay/pkg/udp/ping"
	"os"
	"time"
)

func main() {
	var host string
	host, err := os.Hostname()
	if err != nil {
		host = "Unknown"
	}
	address := flag.String("a", "127.0.0.1:9999", "address to ping")
	interval := flag.Duration("i", time.Second, "ping interval")
	size := flag.Int("s", 100, "UDP datagram size")
	alias := flag.String("alias", fmt.Sprintf("%s:%d", host, *size), "Client alias")
	flag.Parse()
	ping.Ping(*address, *alias, *interval, *size)
}
