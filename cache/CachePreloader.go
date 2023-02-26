package cache

import (
	"github.com/miekg/dns"
	"github.com/jckimble/cachedns/cdns"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

type CachePreloader struct {
	Urls []string
}

func (cp CachePreloader) Load(cache cdns.Cache) {
	for _, url := range cp.Urls {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		for _, line := range strings.Split(strings.Trim(string(content), " \t\r\n"), "\n") {
			line = strings.Replace(strings.Trim(line, " \t"), "\t", " ", -1)
			if len(line) == 0 || line[0] == ';' || line[0] == '#' {
				continue
			}
			pieces := strings.SplitN(line, " ", 2)
			if len(pieces) > 1 && len(pieces[0]) > 0 {
				if names := strings.Fields(pieces[1]); len(names) > 0 {
					for _, name := range names {
						cp.addCache(name+".", cache)
					}
				}
			}
		}

	}
}

func (cp CachePreloader) addCache(domain string, cache cdns.Cache) {
	rr := &dns.A{
		Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
		A:   net.ParseIP("127.0.0.1").To4(),
	}
	cache.SaveAnswer(dns.Question{
		Name: domain,
	}, []dns.RR{rr}, false)
}
