package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/miekg/dns"
)

var resolveIP net.IP

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	q := r.Question[0]

	info := fmt.Sprintf("Question: Type=%s Class=%s Name=%s", dns.TypeToString[q.Qtype], dns.ClassToString[q.Qclass], q.Name)

	if q.Qtype == dns.TypeA && q.Qclass == dns.ClassINET {
		m := new(dns.Msg)
		m.SetReply(r)
		a := new(dns.A)
		a.Hdr = dns.RR_Header{Name: q.Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 600}
		a.A = resolveIP
		m.Answer = []dns.RR{a}
		w.WriteMsg(m)
		log.Printf("%s (RESOLVED)\n", info)
	} else {
		m := new(dns.Msg)
		m.SetReply(r)
		m.Rcode = dns.RcodeNameError // NXDOMAIN
		w.WriteMsg(m)
		log.Printf("%s (NXDOMAIN)\n", info)
	}
}

const EVAL = `
set -x
: === configure DNS:
sudo mkdir -p /etc/resolver; echo -e 'nameserver 127.0.0.1\nport 5300\n'|sudo tee /etc/resolver/dev

: === restart container:
docker rm -f devdns; docker run -d --name devdns -p 5300:5300/udp lalyos/devdns -addr 0.0.0.0:5300

: === test it
ping -c 3 whatever.dev
set +x
`
const USAGE = `
devdns is a DNS server that replies 127.0.0.1 always
full docs: https://github.com/lalyos/devdns

oneliner on OSX:
  eval "$(docker run --rm lalyos/devdns --eval)"
`

func main() {
	var addr = flag.String("addr", "127.0.0.1:5300", "listen address")
	var ip = flag.String("ip", "127.0.0.1", "resolve ipv4 address")
	var usage = flag.Bool("usage", false, "prints docker usage information")
	var eval = flag.Bool("eval", false, "generates script to be shell 'eval' -ed")
	flag.Parse()

	resolveIP = net.ParseIP(*ip)

	if *usage {
		fmt.Println(USAGE)
		os.Exit(0)
	}
	if *eval {
		fmt.Println(EVAL)
		os.Exit(0)
	}
	if resolveIP == nil {
		log.Fatalf("Invalid ip address: %s\n", *ip)
	}

	if resolveIP.To4() == nil {
		log.Fatalf("Invalid ipv4 address: %s\n", *ip)
	}

	server := &dns.Server{Addr: *addr, Net: "udp"}
	server.Handler = dns.HandlerFunc(handleRequest)

	log.Printf("Listening on %s, resolving to %s\n", *addr, *ip)
	log.Fatal(server.ListenAndServe())
}
