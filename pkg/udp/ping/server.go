package ping

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"text/tabwriter"
	"time"
)

var (
	clients map[uint32]*Client
	mu      sync.Mutex
)

const interval = time.Second * 10

var reportTicker *time.Ticker

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

func (c *Client) String() string {
	c.mu.Lock()
	counter := float64(c.counter)
	highestSeq := float64(c.highestSeq)
	c.mu.Unlock()
	var host string
	host, _, err := net.SplitHostPort(c.addr.String())
	if err != nil {
		host = "Unknown"
	}
	loss := (highestSeq - counter) / highestSeq
	return fmt.Sprintf("%s\t%.3f%%\t", host, loss)
}

// Ping format:
// CLIENT_ID(4) SEQ(4) LEN(2) DATA(LEN)
// NOTE: SEQ is 1-based
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
}

func printReport() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	for _ = range reportTicker.C {
		mu.Lock()
		_clients := clients
		mu.Unlock()
		if len(_clients) > 0 {
			fmt.Fprintln(w, "HOST\tLOSS\t")
		}
		for _, c := range _clients {
			fmt.Fprintln(w, c)
		}
		w.Flush()
	}
}

func ListenAndServe(address string) error {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		return err
	}
	clients = make(map[uint32]*Client)
	defer conn.Close()
	log.Println("Listening for UDP packets on:", conn.LocalAddr())
	reportTicker = time.NewTicker(interval)
	defer reportTicker.Stop()
	go printReport()
	for {
		b := make([]byte, 1024*64)
		n, addr, err := conn.ReadFrom(b)
		if err != nil {
			log.Println(err)
			continue
		}
		go checkPing(conn, addr, b[:n])
	}
	return nil
}
