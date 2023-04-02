package resolvers

import (
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/jckimble/cachedns/cdns"
	"github.com/miekg/dns"
	"net"
)

type DockerContainerResolver struct {
	cdns.Resolver

	containers map[string]labels
}

func (dcr *DockerContainerResolver) Watch() error {
	cli, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}
	listener := make(chan *docker.APIEvents)
	if err := cli.AddEventListener(listener); err != nil {
		return err
	}
	go func() {
		for {
			select {
			case msg, ok := <-listener:
				if !ok {
					break
				}
				if (msg.Action == "start" || msg.Action == "health_status: healthy") && msg.Type == "container" {
					container, err := cli.InspectContainer(msg.Actor.ID)
					if err != nil {
						panic(err)
					}
					if container.Config.Healthcheck != nil && msg.Action == "start" {
						continue
					}
					labels := dcr.parseLabels(container.Config.Labels)
					if !labels.Enabled {
						continue
					}
					if labels.Network != "" {
						if net, ok := container.NetworkSettings.Networks[labels.Network]; ok {
							labels.IPAddress = net.IPAddress
						}
					}
					if labels.IPAddress == "" && len(container.NetworkSettings.Networks) == 1 {
						for _, net := range container.NetworkSettings.Networks {
							labels.IPAddress = net.IPAddress
						}
					}
					dcr.addServer(labels)
				} else if (msg.Action == "stop" || msg.Action == "health_status: unhealthy") && msg.Type == "container" {
					container, err := cli.InspectContainer(msg.Actor.ID)
					if err != nil {
						panic(err)
					}
					labels := dcr.parseLabels(container.Config.Labels)
					if !labels.Enabled {
						continue
					}
					if labels.Network != "" {
						if net, ok := container.NetworkSettings.Networks[labels.Network]; ok {
							labels.IPAddress = net.IPAddress
						}
					}
					if labels.IPAddress == "" && len(container.NetworkSettings.Networks) == 1 {
						for _, net := range container.NetworkSettings.Networks {
							labels.IPAddress = net.IPAddress
						}
					}
					dcr.removeServer(labels)
				}
			}
		}
	}()
	containers, err := cli.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		return err
	}
	for _, container := range containers {
		labels := dcr.parseLabels(container.Labels)
		if !labels.Enabled {
			continue
		}
		if labels.Network != "" {
			if net, ok := container.Networks.Networks[labels.Network]; ok {
				labels.IPAddress = net.IPAddress
			}
		}
		if labels.IPAddress == "" && len(container.Networks.Networks) == 1 {
			for _, net := range container.Networks.Networks {
				labels.IPAddress = net.IPAddress
			}
		}
		dcr.addServer(labels)
	}
	return nil
}

type labels struct {
	Enabled  bool
	Hostname []string
	Network  string

	IPAddress string
}

func (dcr *DockerContainerResolver) addServer(l labels) {
	if dcr.containers == nil {
		dcr.containers = map[string]labels{}
	}
	for _, host := range l.Hostname {
		dcr.containers[host] = l
	}
}
func (dcr *DockerContainerResolver) removeServer(l labels) {
	if dcr.containers == nil {
		return
	}
	for _, host := range l.Hostname {
		if c, ok := dcr.containers[host]; ok {
			if c.IPAddress == l.IPAddress {
				delete(dcr.containers, host)
			}
		}
	}
}

func (DockerContainerResolver) parseLabels(l map[string]string) labels {
	lr := labels{Hostname: []string{}}
	for k, v := range l {
		switch k {
		case "cdns.docker.enabled":
			lr.Enabled = v == "true"
		case "cdns.docker.hostname":
			lr.Hostname = append(lr.Hostname, v)
		case "cdns.docker.network":
			lr.Network = v
		}
	}
	return lr
}

func (dcr DockerContainerResolver) Exchange(r *dns.Msg) (*dns.Msg, error) {
	if dcr.containers != nil {
		if l, ok := dcr.containers[r.Question[0].Name]; ok {
			m := new(dns.Msg)
			m.SetReply(r)
			rr := &dns.A{
				Hdr: dns.RR_Header{Name: r.Question[0].Name, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 0},
				A:   net.ParseIP(l.IPAddress).To4(),
			}
			if r.Question[0].Qtype == dns.TypeA {
				m.Answer = []dns.RR{rr}
			}
		}
	}
	if dcr.Resolver != nil {
		return dcr.Resolver.Exchange(r)
	}
	return nil, fmt.Errorf("No Resolver For DNS")
}
