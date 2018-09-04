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

package kafka

import (
	"flag"
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/jaegertracing/jaeger/pkg/kafka/producer"
)

const (
	configPrefix       = "kafka"
	suffixBrokers      = ".brokers"
	suffixTopic        = ".topic"
	suffixEncoding     = ".encoding"
	suffixSaslEnabled  = ".sasl.enabled"
	suffixSaslUsername = ".sasl.username"
	suffixSaslPassword = ".sasl.password"
	suffixMetadata     = ".metadata"
	suffixTLSEnabled   = ".tls.enabled"
	//TODO Add required TLS config
	suffixVersion = ".version"

	encodingJSON  = "json"
	encodingProto = "protobuf"

	defaultBroker   = "127.0.0.1:9092"
	defaultTopic    = "jaeger-spans"
	defaultEncoding = encodingProto
)

// Options stores the configuration options for Kafka
type Options struct {
	config   producer.Configuration
	topic    string
	encoding string
}

// AddFlags adds flags for Options
func (opt *Options) AddFlags(flagSet *flag.FlagSet) {
	flagSet.String(
		configPrefix+suffixBrokers,
		defaultBroker,
		"The comma-separated list of kafka brokers. i.e. '127.0.0.1:9092,0.0.0:1234'")
	flagSet.String(
		configPrefix+suffixTopic,
		defaultTopic,
		"The name of the kafka topic")
	flagSet.String(
		configPrefix+suffixEncoding,
		defaultEncoding,
		fmt.Sprintf(`Encoding of spans ("%s" or "%s") sent to kafka.`, encodingProto, encodingJSON),
	)
	flagSet.Bool(
		configPrefix+suffixSaslEnabled,
		false,
		fmt.Sprintf("Enable SASL configuration"),
	)
	flagSet.String(
		configPrefix+suffixSaslUsername,
		"",
		fmt.Sprintf("SASL username"),
	)
	flagSet.String(
		configPrefix+suffixSaslPassword,
		"",
		fmt.Sprintf("SASL password"),
	)
	flagSet.Bool(
		configPrefix+suffixTLSEnabled,
		false,
		fmt.Sprintf("Enable TLS configuration"),
	)
}

// InitFromViper initializes Options with properties from viper
func (opt *Options) InitFromViper(v *viper.Viper) {
	auths := make([]producer.Authenticator, 0)
	saslEnabled := v.GetBool(configPrefix + suffixSaslEnabled)
	if saslEnabled {
		authConfig := producer.SASLAuthConfiguration{
			Username: v.GetString(configPrefix + suffixSaslUsername),
			Password: v.GetString(configPrefix + suffixSaslPassword),
		}
		auth, err := producer.NewSASLAuthenticator(authConfig)
		if err != nil {
			panic(fmt.Sprintf("cannot initialize new SASL authenticator: %+v", err))
		}
		auths = append(auths, auth)
	}
	tlsEnabled := v.GetBool(configPrefix + suffixTLSEnabled)
	if tlsEnabled {
		//TODO Build full TLS configuration
		authConfig := producer.TLSAuthConfiguration{
			Enabled: true,
		}
		auth, err := producer.NewTLSAuthenticator(authConfig)
		if err != nil {
			panic(fmt.Sprintf("cannot initialize new TLS authenticator: %+v", err))
		}
		auths = append(auths, auth)
	}
	fullMetadata := true
	metadata := v.Get(configPrefix + suffixMetadata)
	if metadata != nil {
		if m, ok := metadata.(bool); !ok {
			panic(fmt.Sprintf("config value %s%s must be a bool (true/false)", configPrefix, suffixMetadata))
		} else {
			fullMetadata = m
		}
	}
	opt.config = producer.Configuration{
		Brokers:        strings.Split(v.GetString(configPrefix+suffixBrokers), ","),
		Authenticators: auths,
		Metadata:       fullMetadata,
		Version:        v.GetString(configPrefix + suffixVersion),
	}
	opt.topic = v.GetString(configPrefix + suffixTopic)
	opt.encoding = v.GetString(configPrefix + suffixEncoding)
}
