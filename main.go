package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/miekg/dns"
	"github.com/tenfyzhong/cityhash"
)

var (
	h      bool
	p      string
	b      string
	prefix string
)

func parseQuery(m *dns.Msg) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			//log.Printf("Query for %s\n", q.Name)
			ip := int2ip(ip2int(net.ParseIP(prefix+".0.0.0")) + cityhash.CityHash32([]byte(q.Name))>>8).String()
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)
	}

	w.WriteMsg(m)
}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

func init() {
	flag.BoolVar(&h, "h", false, "show this help")
	flag.StringVar(&prefix, "prefix", "11", "Set the `prefix` of the fakeIP (default is 11 for 11.x.x.x)")
	flag.StringVar(&b, "b", "", "Set the listening address (default is dual-stack)")
	flag.StringVar(&p, "p", "53", "Set the listening `port` (default is 53)")
}

func main() {
	// 解析命令行参数
	flag.Parse()

	if h {
		flag.Usage()
		os.Exit(0)
	}

	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	port, errPort := strconv.Atoi(p)
	if errPort != nil {
		panic(errPort)
	}
	if port > 65535 || port < 0 {
		panic("Error: the port range must be between 0 and 65535")
	}

	ipPrefix, errIP := strconv.Atoi(prefix)
	if errIP != nil {
		panic(errPort)
	}
	if ipPrefix > 255 || ipPrefix < 0 {
		panic("Error: the prefix range must be between 0 and 255")
	}
	// start server

	server := &dns.Server{Addr: b + ":" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
