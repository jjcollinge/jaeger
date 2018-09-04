// Copyright (c) 2018 The Jaeger Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package producer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/Shopify/sarama"
)

// TLS Conifg
type TLS struct {
	Enabled                bool
	CertPath               string
	KeyPath                string
	CaPath                 string
	EnableHostVerification bool
}

// SASLAuthenticator holds the username and password used for SASL authentication
type SASLAuthenticator struct {
	Enabled  bool   `yaml:"enabled"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Builder builds a new kafka producer
type Builder interface {
	NewProducer() (sarama.AsyncProducer, error)
}

// Configuration describes the configuration properties needed to create a Kafka producer
type Configuration struct {
	Brokers      []string
	SASL         SASLAuthenticator
	TLS          TLS
	MetadataFull bool
	Version      string
}

// NewProducer creates a new asynchronous kafka producer
func (c *Configuration) NewProducer() (sarama.AsyncProducer, error) {
	saramaConfig, err := NewConfig(c)
	if err != nil {
		return nil, err
	}
	return sarama.NewAsyncProducer(c.Brokers, saramaConfig)
}

// NewConfig creates a new sarama configuration
func NewConfig(c *Configuration) (*sarama.Config, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	if c.SASL.Enabled {
		cfg.Net.SASL.Enable = true
		cfg.Net.SASL.User = c.SASL.Username
		cfg.Net.SASL.Password = c.SASL.Password
	}
	if c.TLS.Enabled {
		cfg.Net.TLS.Enable = true
		tlsConfig, err := NewTLSConfig(c.TLS.CertPath, c.TLS.KeyPath, c.TLS.CaPath, !c.TLS.EnableHostVerification)
		if err != nil {
			return nil, fmt.Errorf("error creating TLS config: %+v", err)
		}
		cfg.Net.TLS.Config = tlsConfig
	}
	if c.Version == "" {
		cfg.Version = sarama.MinVersion
	} else {
		v, err := sarama.ParseKafkaVersion(c.Version)
		if err != nil {
			return nil, fmt.Errorf("error parsing kafka version: %+v", err)
		}
		cfg.Version = v
	}
	cfg.Metadata.Full = c.MetadataFull
	return cfg, nil
}

// NewTLSConfig constructs a new TLS config
func NewTLSConfig(clientCertFile, clientKeyFile, caCertFile string, insecureSkipVerify bool) (*tls.Config, error) {
	var err error

	tlsConfig := tls.Config{}
	clientAuth := tls.NoClientCert

	caCertPool := x509.NewCertPool()
	if caCertFile != "" {
		caCert, err := ioutil.ReadFile(caCertFile) // #nosec
		if err != nil {
			return &tlsConfig, err
		}
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caCertPool
	}

	cert := tls.Certificate{}
	if clientCertFile != "" && clientKeyFile != "" {
		cert, err = tls.LoadX509KeyPair(clientCertFile, clientKeyFile)
		if err != nil {
			return &tlsConfig, err
		}
	}

	/* #nosec */
	tlsConfig = tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: insecureSkipVerify,
		ClientAuth:         clientAuth,
	}
	return &tlsConfig, nil
}
