package ssc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Trust raccoglie gli elementi di certificazione della connessione
type Trust struct {
	Cert       tls.Certificate
	Pool       *x509.CertPool
	cert       string
	key        string
	ClientAuth tls.ClientAuthType
}

// NewTrust crea un nuovo trust a partire dai certificati passati come argomento
func NewTrust(certificate, key, ca string, auth tls.ClientAuthType) (*Trust, error) {
	root, err := os.ReadFile(ca)
	if err != nil {
		return nil, err
	}
	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	pool.AppendCertsFromPEM(root)
	c, err := tls.LoadX509KeyPair(certificate, key)
	if err != nil {
		return nil, err
	}
	return &Trust{
		Pool:       pool,
		Cert:       c,
		ClientAuth: auth,
		cert:       certificate,
		key:        key,
	}, nil
}

// Client crea client HTTP con certificati trusted
func (t *Trust) Client(timeout time.Duration) *http.Client {
	return &http.Client{
		Transport: t.transport(),
		Timeout:   timeout,
	}
}

// StartServer avvia un server HTTPS configurato per gestire le richieste definite dal Trust
func (t *Trust) StartServer(h http.Handler, port int) error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
		Handler:      h,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			ClientCAs:  t.Pool,
			ClientAuth: t.ClientAuth,
		},
	}
	return server.ListenAndServeTLS(t.cert, t.key)
}

// transport crea un Transport certificato
func (t *Trust) transport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:   tls.VersionTLS12,
			Certificates: []tls.Certificate{t.Cert},
			RootCAs:      t.Pool,
		},
		DisableKeepAlives: false,
	}
}
