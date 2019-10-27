package main

import (
	"log"
	"net"
	"os"
	"strings"

	"github.com/miekg/dns"
	flag "github.com/spf13/pflag"
)

var domainsToAddresses map[string]net.IP = map[string]net.IP{}

type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := domainsToAddresses[domain]
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   address,
			})

			log.Printf("A response: %s => %s \n", domain, address.String())
		}
	}
	w.WriteMsg(&msg)
}

func main() {
	listen := flag.String("listen", ":8053", "which UDP port to listen")
	aList := flag.StringSlice("a-record", nil, "name:IP pairs")
	if aList == nil {
		log.Fatalln("At least one A record should be given")
	}

	flag.Parse()

	for _, option := range *aList {
		nameIP := strings.Split(option, ":")
		if len(nameIP) != 2 {
			log.Printf("a-record option '%s' is not in name:IP form, skipping", option)
			continue
		}

		name, ipStr := nameIP[0], nameIP[1]
		if name == "" {
			log.Printf("empty name for option '%s'", option)
			continue
		}

		if name[len(name)-1] != '.' {
			name = name + "."
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			log.Printf("bad IP address for option '%s'", option)
			continue
		}

		domainsToAddresses[name] = ip
	}

	srv := &dns.Server{Addr: *listen, Net: "udp"}
	srv.Handler = &handler{}

	log.Printf("%s listening at %s...\n", os.Args[0], *listen)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
	}
}
