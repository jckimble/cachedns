package resolvers

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/jckimble/cachedns/cdns"
	"net"
	"regexp"
)

type FilteringResolver struct {
	cdns.Resolver
	Filters map[string]string
}

func (fr FilteringResolver) Exchange(r *dns.Msg) (*dns.Msg, error) {
	for filter, ip := range fr.Filters {
		match, err := regexp.MatchString(filter, r.Question[0].Name)
		if err != nil {
			return nil, err
		}
		if match {
			m := new(dns.Msg)
			m.SetReply(r)
			rr := &dns.A{
				Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
				A:   net.ParseIP(ip).To4(),
			}
			if r.Question[0].Qtype == dns.TypeA {
				m.Answer = []dns.RR{rr}
			}
			return m, nil
		}
	}
	if fr.Resolver != nil {
		return fr.Resolver.Exchange(r)
	}
	return nil, fmt.Errorf("No Resolver For DNS")
}
