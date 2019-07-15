package ping

import (
	"bytes"
	crand "crypto/rand"
	"encoding/binary"
	"log"
	"math/rand"
	"net"
	"time"
)

func Ping(address string, alias string, interval time.Duration, size int) {
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	dst, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	var seq uint32
	id := rand.Uint32()
	buf := new(bytes.Buffer)
	d := make([]byte, size-12-len(alias))

	_, err = crand.Read(d)
	if err != nil {
		log.Fatal(err)
	}

	tick := time.Tick(interval)
	for _ = range tick {
		seq = seq + 1
		log.Println("Ping:", seq)
		var data = []interface{}{
			pingHeader{Id: id, Seq: seq, AliasLen: uint16(len(alias)), Len: 1},
			[]byte(alias),
			d,
		}

		for _, v := range data {
			err = binary.Write(buf, binary.BigEndian, v)
			if err != nil {
				log.Println(err)
				break
			}
		}

		_, err = conn.WriteTo(buf.Bytes(), dst)
		if err != nil {
			log.Println(err)
			return
		}
		buf.Reset()
	}
}
