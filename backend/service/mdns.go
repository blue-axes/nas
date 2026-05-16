package service

import (
	stdlog "log"
	"net"
	"strings"

	"github.com/hashicorp/mdns"
	log "github.com/sirupsen/logrus"
)

func ipStrings(ips []net.IP) string {
	s := make([]string, len(ips))
	for i, ip := range ips {
		s[i] = ip.String()
	}
	return strings.Join(s, ", ")
}

func (svc *Service) StartMDNS() error {
	cfg := svc.Config()
	if !cfg.MDNS.Enabled {
		return nil
	}

	hostname := cfg.MDNS.Hostname
	if hostname[len(hostname)-1] != '.' {
		hostname += "."
	}

	mdnsLogger := stdlog.New(log.StandardLogger().WriterLevel(log.DebugLevel), "", 0)

	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		if strings.HasPrefix(iface.Name, "docker") ||
			strings.HasPrefix(iface.Name, "br-") ||
			strings.HasPrefix(iface.Name, "veth") ||
			strings.HasPrefix(iface.Name, "virbr") {
			continue
		}
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 读取ip地址
		addrs, err := iface.Addrs()
		if err != nil {
			log.Errorf("get intf:%s addr error:%+v", iface.Name, err)
			continue
		}
		var (
			ips     = []net.IP{}
			hasIPv4 = false
			hasIPv6 = false
		)
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if !ipnet.IP.IsGlobalUnicast() && !ipnet.IP.IsLinkLocalUnicast() {
					continue
				}
				if ipnet.IP.To4() != nil {
					hasIPv4 = true
					ips = append(ips, ipnet.IP)
				} else if ipnet.IP.To16() != nil {
					hasIPv6 = true
					ips = append(ips, ipnet.IP)
				}
			}
		}

		if len(ips) == 0 {
			continue
		}
		service, err := mdns.NewMDNSService(
			cfg.MDNS.Info,
			cfg.MDNS.ServiceName,
			"",
			hostname,
			int(cfg.Http.ListenPort),
			ips,
			[]string{"path=/"},
		)
		if err != nil {
			log.Errorf("NewMDNSService on intf:%s ips:%s error:%+v", iface.Name, ipStrings(ips), err)
			continue
		}

		server, err := mdns.NewServer(&mdns.Config{
			Zone:              service,
			Iface:             &iface,
			LogEmptyResponses: true,
			Logger:            mdnsLogger,
		})
		if err != nil {
			log.Warnf("mDNS: failed to create server on %s: %v", iface.Name, err)
			continue
		}
		svc.mdnsServers = append(svc.mdnsServers, server)
		log.Infof("mDNS: broadcasting %s as %s  on %s with IPs [%s](v4=%v v6=%v)", cfg.MDNS.ServiceName, cfg.MDNS.Hostname, iface.Name, ipStrings(ips), hasIPv4, hasIPv6)
	}

	if len(svc.mdnsServers) == 0 {
		log.Warn("mDNS: no servers could be started")
		return nil
	}

	return nil
}

func (svc *Service) ShutdownMDNS() {
	for _, s := range svc.mdnsServers {
		s.Shutdown()
	}
	if len(svc.mdnsServers) > 0 {
		svc.mdnsServers = nil
		log.Info("mDNS: stopped")
	}
}
