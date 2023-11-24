//go:build dev

package envtls

import (
	"crypto/tls"
	"github.com/clarkmcc/cloudcore/internal/envtls/example"
)

func TLSConfig() (*tls.Config, error) {
	return example.TLSConfig(), nil
}
