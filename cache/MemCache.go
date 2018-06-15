package cache

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/patrickmn/go-cache"
	"time"
)

type MemCache struct {
	cache *cache.Cache
}

func (mc MemCache) GetAnswer(q dns.Question) ([]dns.RR, error) {
	if mc.cache == nil {
		mc.cache = cache.New(5*time.Minute, 15*time.Minute)
	}
	domain, ok := mc.cache.Get(q.Name)
	if ok {
		return domain.([]dns.RR), nil
	}
	return nil, fmt.Errorf("Not Found")
}
func (mc *MemCache) SaveAnswer(q dns.Question, r []dns.RR, expires bool) {
	if mc.cache == nil {
		mc.cache = cache.New(5*time.Minute, 15*time.Minute)
	}
	if !expires {
		mc.cache.Set(q.Name, r, cache.NoExpiration)
	} else {
		mc.cache.SetDefault(q.Name, r)
	}
}
