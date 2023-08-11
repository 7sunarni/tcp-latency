package main

import (
	"errors"
	"flag"
	"log"
	"net"
	"time"
)

func main() {
	dialTimeout := flag.Int64("timeout", 1, "tcp dial timeout")
	address := flag.String("address", "", "targe tcp address, e.g. 1.2.3.4:5678")
	count := flag.Int64("count", 10, "dial count")
	flag.Parse()
	if *address == "" {
		log.Fatal("address is empty, using --address assign tcp address")
	}
	l, d, err := TcpLatency(*address, *dialTimeout, *count)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("latency %v ms, dail fail %v%%", l, d*100)
}

func timeout(address string, timeout int64) (int64, error) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Second)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return time.Since(start).Milliseconds(), nil
}

func TcpLatency(address string, dialTimeout int64, count int64) (latency float64, dialFailed float64, err error) {
	drop := int64(0)
	sum := int64(0)

	for i := int64(0); i < count; i++ {
		t, err := timeout(address, dialTimeout)
		defer time.Sleep(10 * time.Millisecond)
		if err != nil {
			drop = drop + 1
			continue
		}
		sum = sum + t
	}

	if count == drop {
		return 0, 0, errors.New("dail failed")
	}
	return float64(sum) / float64((count - drop)), float64(drop) / float64(sum), nil

}
