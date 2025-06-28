package tls

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

type TLSConfigProvider interface {
	GetCertFile() string
	GetKeyFile() string
	GetCaFile() string
	GetServerName() string
}

func BuildConfig(provider TLSConfigProvider) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(provider.GetCertFile(), provider.GetKeyFile())
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	var clientCAPool *x509.CertPool
	if provider.GetCaFile() != "" {
		caBytes, err := os.ReadFile(provider.GetCaFile())
		if err != nil {
			return nil, err
		}
		clientCAPool = x509.NewCertPool()
		if !clientCAPool.AppendCertsFromPEM(caBytes) {
			return nil, err
		}
	} else {
		clientCAPool, err = x509.SystemCertPool()
		if err != nil {
			return nil, err
		}
	}
	tlsConfig.ClientCAs = clientCAPool

	return tlsConfig, nil
}

func BuildServerConfig(provider TLSConfigProvider, clientAuth bool) (*tls.Config, error) {
	tlsConfig, err := BuildConfig(provider)
	if err != nil {
		return nil, err
	}

	if clientAuth {
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	} else {
		tlsConfig.ClientAuth = tls.NoClientCert
	}

	return tlsConfig, nil
}
