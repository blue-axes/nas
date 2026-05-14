package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "gen-ca":
		genCA(os.Args[2:])
	case "gen-server":
		genServer(os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Print(`Usage:
  certgen gen-ca    Generate a root CA certificate
  certgen gen-server   Generate a server certificate signed by root CA

  gen-ca flags:
    --ca-cert <path>    Output CA certificate file (default: ca.crt)
    --ca-key <path>     Output CA private key file (default: ca.key)
    --days <int>        Validity in days (default: 3650)
    --org <name>        Organization name (default: My CA)

  gen-server flags:
    --ca-cert <path>    Root CA certificate file (default: ca.crt)
    --ca-key <path>     Root CA private key file (default: ca.key)
    --cert <path>       Output server certificate file (default: server.crt)
    --key <path>        Output server private key file (default: server.key)
    --days <int>        Validity in days (default: 365)
    --host <hosts>      Comma-separated SANs (DNS:xxx,IP:x.x.x.x)
                        Required. Example: DNS:localhost,IP:127.0.0.1,DNS:nas.local
`)
}

func genCA(args []string) {
	fs := flag.NewFlagSet("gen-ca", flag.ExitOnError)
	caCert := fs.String("ca-cert", "ca.crt", "")
	caKey := fs.String("ca-key", "ca.key", "")
	days := fs.Int("days", 3650, "")
	org := fs.String("org", "My CA", "")
	fs.Parse(args)

	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate CA key: %v\n", err)
		os.Exit(1)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate serial: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:   *org,
			Organization: []string{*org},
		},
		NotBefore:             now,
		NotAfter:              now.Add(time.Duration(*days) * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create CA certificate: %v\n", err)
		os.Exit(1)
	}

	savePEM(*caCert, "CERTIFICATE", certDER)
	saveECKey(*caKey, key)
	fmt.Printf("Root CA generated:\n  Certificate: %s\n  Private Key: %s\n", *caCert, *caKey)
}

func genServer(args []string) {
	fs := flag.NewFlagSet("gen-server", flag.ExitOnError)
	caCert := fs.String("ca-cert", "ca.crt", "")
	caKey := fs.String("ca-key", "ca.key", "")
	certOutput := fs.String("cert", "server.crt", "")
	keyOutput := fs.String("key", "server.key", "")
	days := fs.Int("days", 365, "")
	host := fs.String("host", "", "")
	fs.Parse(args)

	if *host == "" {
		fmt.Fprintln(os.Stderr, "error: --host is required (e.g., DNS:localhost,IP:127.0.0.1,DNS:nas.local)")
		os.Exit(1)
	}

	caCertPEM, err := os.ReadFile(*caCert)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read CA cert: %v\n", err)
		os.Exit(1)
	}
	caKeyPEM, err := os.ReadFile(*caKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read CA key: %v\n", err)
		os.Exit(1)
	}

	caCertBlock, _ := pem.Decode(caCertPEM)
	if caCertBlock == nil {
		fmt.Fprintln(os.Stderr, "failed to decode CA certificate PEM")
		os.Exit(1)
	}
	caCertParsed, err := x509.ParseCertificate(caCertBlock.Bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse CA certificate: %v\n", err)
		os.Exit(1)
	}

	caKeyBlock, _ := pem.Decode(caKeyPEM)
	if caKeyBlock == nil {
		fmt.Fprintln(os.Stderr, "failed to decode CA private key PEM")
		os.Exit(1)
	}
	caKeyParsed, err := x509.ParseECPrivateKey(caKeyBlock.Bytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse CA private key: %v\n", err)
		os.Exit(1)
	}

	serverKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate server key: %v\n", err)
		os.Exit(1)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to generate serial: %v\n", err)
		os.Exit(1)
	}

	now := time.Now()
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName: parseFirstDNS(*host),
		},
		NotBefore: now,
		NotAfter:  now.Add(time.Duration(*days) * 24 * time.Hour),
		KeyUsage:  x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
	}

	for _, part := range strings.Split(*host, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "DNS:") {
			tmpl.DNSNames = append(tmpl.DNSNames, part[4:])
		} else if strings.HasPrefix(part, "IP:") {
			ip := net.ParseIP(part[3:])
			if ip != nil {
				tmpl.IPAddresses = append(tmpl.IPAddresses, ip)
			}
		}
	}

	certDER, err := x509.CreateCertificate(rand.Reader, tmpl, caCertParsed, &serverKey.PublicKey, caKeyParsed)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create server certificate: %v\n", err)
		os.Exit(1)
	}

	savePEM(*certOutput, "CERTIFICATE", certDER)
	saveECKey(*keyOutput, serverKey)
	fmt.Printf("Server certificate generated:\n  Certificate: %s\n  Private Key: %s\n", *certOutput, *keyOutput)
}

func parseFirstDNS(host string) string {
	for _, part := range strings.Split(host, ",") {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "DNS:") {
			return part[4:]
		}
	}
	return "localhost"
}

func savePEM(path, blockType string, der []byte) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create %s: %v\n", path, err)
		os.Exit(1)
	}
	defer f.Close()
	pem.Encode(f, &pem.Block{Type: blockType, Bytes: der})
}

func saveECKey(path string, key *ecdsa.PrivateKey) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create %s: %v\n", path, err)
		os.Exit(1)
	}
	defer f.Close()
	der, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal EC key: %v\n", err)
		os.Exit(1)
	}
	pem.Encode(f, &pem.Block{Type: "EC PRIVATE KEY", Bytes: der})
}