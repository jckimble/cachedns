package main

import (
	"gitlab.com/paranoidsecurity/cachedns/cache"
	"gitlab.com/paranoidsecurity/cachedns/cdns"
	"gitlab.com/paranoidsecurity/cachedns/resolvers"

	"github.com/spf13/viper"

	"log"
	"time"
)

func init() {
	viper.SetDefault("preload", []string{})
	viper.SetDefault("filters", map[string]string{})
	viper.SetDefault("resolvers", []string{"208.67.220.220", "208.67.222.222"})
	viper.SetDefault("port", ":5353")

	viper.AddConfigPath("/etc")
	viper.AddConfigPath(".")
	viper.SetConfigName("cdns")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("Defaults Used: %s\n", err)
	}
}

func main() {
	memcache := &cache.MemCache{}
	preload := &cache.CachePreloader{
		Urls: viper.GetStringSlice("preload"),
	}
	go func() {
		for {
			preload.Load(memcache)
			timer := time.NewTimer(24 * time.Hour)
			<-timer.C
		}
	}()
	cd := cdns.CacheDns{
		Port:  viper.GetString("port"),
		Cache: memcache,
		Resolver: resolvers.FilteringResolver{
			Filters: viper.GetStringMapString("filters"),
			ForwardingResolver: resolvers.ForwardingResolver{
				viper.GetStringSlice("resolvers"),
			},
		},
	}
	cd.Run()
}
