package ping

import (
	"log"
	"net"
)

func ListenAndServe(address string) error {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Println("Listening for UDP packets on:", conn.LocalAddr())
	b := make([]byte, 1500)
	for {
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			log.Println(err)
		}
		_, err = conn.WriteTo(b[:n], addr)
		if err != nil {
			log.Println(err)
		}
		log.Println("Got UDP packet from: ", addr)
	}
	return nil
}
