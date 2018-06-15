package cdns

import (
	"github.com/miekg/dns"
)

type Resolver interface {
	Exchange(*dns.Msg) (*dns.Msg, error)
}

type Cache interface {
	GetAnswer(dns.Question) ([]dns.RR, error)
	SaveAnswer(dns.Question, []dns.RR, bool)
}
