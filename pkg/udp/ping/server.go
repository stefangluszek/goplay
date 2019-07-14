package ping

import (
	"bytes"
	"encoding/binary"
	"log"
	"net"
	"sync"
)

var (
	clients map[uint32]*Client
	mu      sync.Mutex
)

type Client struct {
	id         uint32
	addr       net.Addr
	counter    int
	highestSeq uint32
	mu         sync.Mutex
}

type pingHeader struct {
	Id  uint32
	Seq uint32
	Len uint16
}

func getClient(h pingHeader, addr net.Addr) *Client {
	mu.Lock()
	defer mu.Unlock()

	if c, ok := clients[h.Id]; ok {
		return c
	}
	clients[h.Id] = &Client{addr: addr, id: h.Id}
	return clients[h.Id]
}

func (c *Client) ping(h pingHeader) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if h.Seq > c.highestSeq {
		c.highestSeq = h.Seq
	}
	c.counter++
}

// Ping format:
// CLIENT_ID(4) SEQ(4) LEN(2) DATA(LEN)
func checkPing(conn net.PacketConn, addr net.Addr, b []byte) {

	var header pingHeader
	r := bytes.NewReader(b)

	if err := binary.Read(r, binary.BigEndian, &header); err != nil {
		log.Println("Failed to parse ping header:", err)
		return
	}
	data := make([]byte, header.Len)

	if err := binary.Read(r, binary.BigEndian, &data); err != nil {
		log.Println("Failed to parse ping data:", err)
		return
	}

	c := getClient(header, addr)
	c.ping(header)
	log.Println(c)
}

func ListenAndServe(address string) error {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	clients = make(map[uint32]*Client)
	defer conn.Close()
	log.Println("Listening for UDP packets on:", conn.LocalAddr())
	for {
		b := make([]byte, 1024*64)
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println("Got UDP packet from: ", addr)
		go checkPing(conn, addr, b[:n])
	}
	return nil
}
