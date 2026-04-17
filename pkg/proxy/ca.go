package proxy

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type CAManager struct {
	caCert *x509.Certificate
	caKey  interface{}
	dir    string
	mu     sync.Mutex
	certs  map[string]*tls.Certificate
}

func NewCAManager(baseDir string) (*CAManager, error) {
	dir := filepath.Join(baseDir, ".agent-replay", "ca")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	cm := &CAManager{
		dir:   dir,
		certs: make(map[string]*tls.Certificate),
	}

	if err := cm.loadOrGenerateCA(); err != nil {
		return nil, err
	}

	return cm, nil
}

func (cm *CAManager) loadOrGenerateCA() error {
	certPath := filepath.Join(cm.dir, "ca.crt")
	keyPath := filepath.Join(cm.dir, "ca.key")

	if _, err := os.Stat(certPath); err == nil {
		catls, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err == nil {
			cm.caKey = catls.PrivateKey
			cm.caCert, _ = x509.ParseCertificate(catls.Certificate[0])
			return nil
		}
	}

	// Generate new CA
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"agrepl CA"},
			CommonName:   "agrepl Interception Authority",
		},
		NotBefore:             time.Now().Add(-24 * time.Hour),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// Save cert
	certOut, _ := os.Create(certPath)
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()

	// Save key
	keyOut, _ := os.OpenFile(keyPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	privBytes, _ := x509.MarshalECPrivateKey(priv)
	pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})
	keyOut.Close()

	cm.caKey = priv
	cm.caCert, _ = x509.ParseCertificate(derBytes)

	fmt.Printf("\033[33m[CA] Generated new Root CA at %s\033[0m\n", certPath)
	fmt.Printf("\033[33m[CA] To intercept HTTPS, you may need to trust this certificate in your system or environment.\033[0m\n")

	return nil
}

func (cm *CAManager) GetCertificate(host string) (*tls.Certificate, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// Strip port if present
	h, _, err := net.SplitHostPort(host)
	if err == nil {
		host = h
	}

	if cert, ok := cm.certs[host]; ok {
		return cert, nil
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: host,
		},
		NotBefore:    time.Now().Add(-1 * time.Hour),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{host},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, cm.caCert, &priv.PublicKey, cm.caKey)
	if err != nil {
		return nil, err
	}

	cert := &tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
	}
	cm.certs[host] = cert
	return cert, nil
}
