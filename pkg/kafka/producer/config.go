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
	"fmt"

	"github.com/Shopify/sarama"
)

// Builder builds a new kafka producer
type Builder interface {
	NewProducer() (sarama.AsyncProducer, error)
}

// Configuration describes the configuration properties needed to create a Kafka producer
type Configuration struct {
	Brokers        []string
	Authenticators []Authenticator
	Metadata       bool
	Version        string
}

// NewProducer creates a new asynchronous kafka producer
func (c *Configuration) NewProducer() (sarama.AsyncProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	var err error
	// Last write wins on conflict
	for _, auth := range c.Authenticators {
		saramaConfig, err = auth.Authenticate(saramaConfig)
		if err != nil {
			return nil, err
		}
	}
	saramaConfig.Metadata.Full = c.Metadata
	if c.Version == "" {
		saramaConfig.Version = sarama.MinVersion
	} else {
		v, err := sarama.ParseKafkaVersion(c.Version)
		if err != nil {
			return nil, err
		}
		saramaConfig.Version = v
	}
	return sarama.NewAsyncProducer(c.Brokers, saramaConfig)
}

type AuthConfiguration interface{}

type SASLAuthConfiguration struct {
	Username string
	Password string
}

type TLSAuthConfiguration struct {
	Enabled bool
	Config  *tls.Config
}

type Authenticator interface {
	Authenticate(config *sarama.Config) (*sarama.Config, error)
}

type SASLAuthenticator struct {
	config *SASLAuthConfiguration
}

func NewSASLAuthenticator(auth AuthConfiguration) (*SASLAuthenticator, error) {
	c, ok := auth.(SASLAuthConfiguration)
	if !ok {
		return nil, fmt.Errorf("cannot type assert AuthConfiguration into SASLAuthConfiguration")
	}
	return &SASLAuthenticator{
		config: &c,
	}, nil
}

func (s *SASLAuthenticator) Authenticate(config *sarama.Config) (*sarama.Config, error) {
	if s.config == nil || config == nil {
		return nil, fmt.Errorf("error")
	}
	config.Net.SASL.Enable = true
	config.Net.SASL.User = s.config.Username
	config.Net.SASL.Password = s.config.Password
	return config, nil
}

type TLSAuthenticator struct {
	config *TLSAuthConfiguration
}

func NewTLSAuthenticator(auth AuthConfiguration) (*TLSAuthenticator, error) {
	c, ok := auth.(TLSAuthConfiguration)
	if !ok {
		return nil, fmt.Errorf("cannot type assert AuthConfiguration into TLSAuthConfiguration")
	}
	return &TLSAuthenticator{
		config: &c,
	}, nil
}

func (t *TLSAuthenticator) Authenticate(config *sarama.Config) (*sarama.Config, error) {
	if t.config == nil || config == nil {
		return nil, fmt.Errorf("error")
	}
	config.Net.TLS.Enable = t.config.Enabled
	if t.config.Config != nil {
		config.Net.TLS.Config = t.config.Config
	}
	return config, nil
}
