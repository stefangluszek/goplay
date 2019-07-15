package main

import (
	"flag"
	"fmt"
	"github.com/stefangluszek/goplay/pkg/udp/ping"
)

func main() {
	address := flag.String("a", ":9999", "address to listen on")
	flag.Parse()
	fmt.Println(ping.ListenAndServe(*address))
}
