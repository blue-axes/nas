package service

import (
	"net"

	"github.com/hashicorp/mdns"
	log "github.com/sirupsen/logrus"
)

func getLocalIPv4s() []net.IP {
	var ips []net.IP
	ifaces, err := net.Interfaces()
	if err != nil {
		return ips
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					ips = append(ips, ip4)
				}
			}
		}
	}
	return ips
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

	ips := getLocalIPv4s()
	if len(ips) == 0 {
		log.Warn("mDNS: no local IPv4 address found, skipping")
		return nil
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
		return err
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		return err
	}

	svc.mdnsServer = server
	log.Infof("mDNS: broadcasting %s as %s.local:%d", cfg.MDNS.ServiceName, cfg.MDNS.Hostname, cfg.Http.ListenPort)
	return nil
}

func (svc *Service) ShutdownMDNS() {
	if svc.mdnsServer != nil {
		svc.mdnsServer.Shutdown()
		log.Info("mDNS: stopped")
	}
}
