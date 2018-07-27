package cdns

import (
	"github.com/miekg/dns"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type CacheDns struct {
	Port     string
	Cache    Cache
	Resolver Resolver
}

func (cd CacheDns) HandleRequest(w dns.ResponseWriter, r *dns.Msg) {
	answer, err := cd.Cache.GetAnswer(r.Question[0])
	m := new(dns.Msg)
	m.SetReply(r)
	if err != nil {
		if cd.Resolver != nil {
			msg, err := cd.Resolver.Exchange(r)
			if err != nil {
				log.Printf("Error: %s\n", err)
				return
			}
			m.Answer = msg.Answer
			cd.Cache.SaveAnswer(r.Question[0], msg.Answer, true)
		}
	} else {
		m.Authoritative = true
		m.Answer = answer
	}
	w.WriteMsg(m)
}

func (cd CacheDns) Run() {
	dns.HandleFunc(".", cd.HandleRequest)
	go func() {
		srv := &dns.Server{Addr: cd.Port, Net: "udp"}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to set listener %s\n", err)
		}
	}()
	go func() {
		srv := &dns.Server{Addr: cd.Port, Net: "tcp"}
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatalf("Failed to set listener %s\n", err)
		}
	}()
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case s := <-sig:
			log.Fatalf("Signal (%d) received, stopping\n", s)
		}
	}
}
