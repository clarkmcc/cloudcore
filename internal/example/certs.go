package example

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
)

// Requires mkcert to be installed (https://github.com/FiloSottile/mkcert).
//go:generate mkcert example.com

//go:embed example.com.pem
var certificate []byte

//go:embed example.com-key.pem
var key []byte

func TLSConfig() *tls.Config {
	cert, err := tls.X509KeyPair(certificate, key)
	if err != nil {
		panic(err)
	}
	pool, err := x509.SystemCertPool()
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
}
