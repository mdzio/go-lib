package httputil

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"net"
	"time"
)

// CertGenerator generates CA and server certificates for serving HTTPS.
type CertGenerator struct {
	Hosts          []string
	Organization   string
	NotBefore      time.Time
	NotAfter       time.Time
	CACertFile     string
	CAKeyFile      string
	ServerCertFile string
	ServerKeyFile  string
}

func writeKey(fileName string, key *ecdsa.PrivateKey) error {
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return fmt.Errorf("Unable to marshal ECDSA private key: %v", err)
	}
	keyBytes := new(bytes.Buffer)
	if err = pem.Encode(keyBytes, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}); err != nil {
		return fmt.Errorf("Failed to encode private key: %v", err)
	}
	if err = ioutil.WriteFile(fileName, keyBytes.Bytes(), 0600); err != nil {
		return fmt.Errorf("Failed to write %s: %v", fileName, err)
	}
	return nil
}

func writeCert(fileName string, cert []byte) error {
	certBytes := new(bytes.Buffer)
	if err := pem.Encode(certBytes, &pem.Block{Type: "CERTIFICATE", Bytes: cert}); err != nil {
		return fmt.Errorf("Failed to encode certificate: %v", err)
	}
	if err := ioutil.WriteFile(fileName, certBytes.Bytes(), 0666); err != nil {
		return fmt.Errorf("Failed to write %s: %v", fileName, err)
	}
	return nil
}

// Generate creates the certificate files.
func (g *CertGenerator) Generate() error {
	// generate CA private key (use ECDSA curve P256)
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("Failed to generate private CA key: %v", err)
	}
	if err = writeKey(g.CAKeyFile, caKey); err != nil {
		return err
	}

	// generate CA certificate
	snLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	sn, err := rand.Int(rand.Reader, snLimit)
	if err != nil {
		return fmt.Errorf("Failed to generate serial number: %v", err)
	}
	caTmpl := x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization: []string{g.Organization},
			CommonName:   g.Organization + " CA",
		},
		NotBefore:             g.NotBefore,
		NotAfter:              g.NotAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  true,
		BasicConstraintsValid: true,
	}
	caCert, err := x509.CreateCertificate(rand.Reader, &caTmpl, &caTmpl, &caKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("Failed to create CA certificate: %v", err)
	}
	if err = writeCert(g.CACertFile, caCert); err != nil {
		return err
	}

	// generate server private key (use ECDSA curve P256)
	serverKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("Failed to generate private server key: %v", err)
	}
	if err = writeKey(g.ServerKeyFile, serverKey); err != nil {
		return err
	}

	// generate server certificate
	sn, err = rand.Int(rand.Reader, snLimit)
	if err != nil {
		return fmt.Errorf("Failed to generate serial number: %v", err)
	}
	serverTmpl := x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Organization: []string{g.Organization},
			CommonName:   g.Organization + " Server",
		},
		NotBefore:             g.NotBefore,
		NotAfter:              g.NotAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IsCA:                  false,
		BasicConstraintsValid: true,
	}
	for _, h := range g.Hosts {
		if ip := net.ParseIP(h); ip != nil {
			serverTmpl.IPAddresses = append(serverTmpl.IPAddresses, ip)
		} else {
			serverTmpl.DNSNames = append(serverTmpl.DNSNames, h)
		}
	}
	serverCert, err := x509.CreateCertificate(rand.Reader, &serverTmpl, &caTmpl, &serverKey.PublicKey, caKey)
	if err != nil {
		return fmt.Errorf("Failed to create server certificate: %v", err)
	}
	if err = writeCert(g.ServerCertFile, serverCert); err != nil {
		return err
	}
	return nil
}
