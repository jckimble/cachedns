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
	viper.SetDefault("resolvers.filtering.filters", map[string]string{})
	viper.SetDefault("resolvers.forwarding.servers", []string{"208.67.220.220", "208.67.222.222"})
	viper.SetDefault("port", ":5353")

	viper.SetDefault("resolvers.filtering.enabled", true)
	viper.SetDefault("resolvers.forwarding.enabled", true)
	viper.SetDefault("resolvers.docker.enabled", false)

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
	var resolver cdns.Resolver
	if viper.GetBool("resolvers.forwarding.enabled") {
		resolver = resolvers.ForwardingResolver{
			Servers: viper.GetStringSlice("resolvers.forwarding.servers"),
		}
	}
	if viper.GetBool("resolvers.filtering.enabled") {
		resolver = resolvers.FilteringResolver{
			Filters:  viper.GetStringMapString("resolvers.filtering.filters"),
			Resolver: resolver,
		}
	}
	if viper.GetBool("resolvers.docker.enabled") {
		resolver = resolvers.DockerContainerResolver{
			Resolver: resolver,
		}
	}
	cd := cdns.CacheDns{
		Port:     viper.GetString("port"),
		Cache:    memcache,
		Resolver: resolver,
	}
	cd.Run()
}
