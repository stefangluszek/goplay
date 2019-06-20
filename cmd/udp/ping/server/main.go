// Simple echo UDP server
package main

import (
	"flag"
	"log"
	"net"
)

func main() {
	address := flag.String("a", "0:9999", "address to listen on")
	conn, err := net.ListenPacket("udp", *address)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println("Listening for UDP packets on:", conn.LocalAddr())

	b := make([]byte, 1500)
	for {
		_, addr, err := conn.ReadFrom(b)
		log.Println("Got UDP packet from: ", addr)
		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.WriteTo(b, addr)
		if err != nil {
			log.Fatal(err)
		}
	}
}
