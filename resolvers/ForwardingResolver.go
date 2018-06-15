package resolvers

import (
	"fmt"
	"github.com/miekg/dns"
	"strings"
	"time"
)

type ForwardingResolver struct {
	Servers []string
}

func (fr ForwardingResolver) Exchange(r *dns.Msg) (*dns.Msg, error) {
	ch := make(chan *dns.Msg)
	for _, server := range fr.Servers {
		if !strings.Contains(server, ":") {
			server = server + ":53"
		}
		go func(dnsServer string) {
			c := new(dns.Client)
			msg, _, err := c.Exchange(r, dnsServer)
			if err == nil {
				ch <- msg
			} else {
				fmt.Printf("%s\n", err)
			}
		}(server)
	}
	select {
	case msg := <-ch:
		return msg, nil
	case <-time.After(3 * time.Second):
		return nil, fmt.Errorf("DNS Timeout")
	}
}
