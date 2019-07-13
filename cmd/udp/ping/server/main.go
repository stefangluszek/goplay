package main

import (
	"flag"
	"fmt"
	"github.com/stefangluszek/goplay/pkg/udp/ping"
)

func main() {
	address := flag.String("a", "0:9999", "address to listen on")
	fmt.Println(ping.ListenAndServe(*address))
}
