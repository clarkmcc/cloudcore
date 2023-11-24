//go:build !dev

package envtls

import (
	"crypto/tls"
	"fmt"
	"os"
)

func TLSConfig() (*tls.Config, error) {
	c := os.Getenv("TLS_CERTIFICATE")
	if len(c) == 0 {
		return nil, fmt.Errorf("missing tls certificate")
	}
	k := os.Getenv("TLS_PRIVATE_KEY")
	if len(k) == 0 {
		return nil, fmt.Errorf("missing tls certificate")
	}
	cert, err := tls.X509KeyPair([]byte(c), []byte(k))
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}
