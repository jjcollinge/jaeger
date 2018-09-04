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
	suffixSasl         = ".sasl"
	suffixSaslUsername = ".sasl.username"
	suffixSaslPassword = ".sasl.password"
	suffixMetadataFull = ".metadata.full"
	suffixTLS          = ".tls"
	suffixCert         = ".tls.cert"
	suffixKey          = ".tls.key"
	suffixCA           = ".tls.ca"
	suffixVerifyHost   = ".tls.verify-host"
	suffixVersion      = ".version"

	encodingJSON  = "json"
	encodingProto = "protobuf"

	defaultBroker       = "127.0.0.1:9092"
	defaultTopic        = "jaeger-spans"
	defaultEncoding     = encodingProto
	defaultMetadataFull = true
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
		configPrefix+suffixSasl,
		opt.config.SASL.Enabled,
		"Enable SASL",
	)
	flagSet.String(
		configPrefix+suffixSaslUsername,
		opt.config.SASL.Username,
		"SASL username",
	)
	flagSet.String(
		configPrefix+suffixSaslPassword,
		opt.config.SASL.Password,
		"SASL password",
	)
	flagSet.Bool(
		configPrefix+suffixMetadataFull,
		defaultMetadataFull,
		"Full metadata for all topics (default=true)",
	)
	flagSet.Bool(
		configPrefix+suffixTLS,
		opt.config.TLS.Enabled,
		"Enable TLS",
	)
	flagSet.String(
		configPrefix+suffixCert,
		opt.config.TLS.CertPath,
		"Path to TLS certificate file",
	)
	flagSet.String(
		configPrefix+suffixKey,
		opt.config.TLS.KeyPath,
		"Path to TLS key file",
	)
	flagSet.String(
		configPrefix+suffixCA,
		opt.config.TLS.CaPath,
		"Path to TLS CA file",
	)
	flagSet.Bool(
		configPrefix+suffixVerifyHost,
		opt.config.TLS.EnableHostVerification,
		"Enable (or disable) host key verification",
	)
	flagSet.String(
		configPrefix+suffixVersion,
		opt.config.Version,
		"Kafka API version",
	)
}

// InitFromViper initializes Options with properties from viper
func (opt *Options) InitFromViper(v *viper.Viper) {
	var cfg producer.Configuration

	cfg.SASL.Enabled = v.GetBool(configPrefix + suffixSasl)
	cfg.SASL.Username = v.GetString(configPrefix + suffixSaslUsername)
	cfg.SASL.Password = v.GetString(configPrefix + suffixSaslPassword)
	cfg.MetadataFull = v.GetBool(configPrefix + suffixMetadataFull)
	cfg.TLS.Enabled = v.GetBool(configPrefix + suffixTLS)
	cfg.TLS.CertPath = v.GetString(configPrefix + suffixCert)
	cfg.TLS.KeyPath = v.GetString(configPrefix + suffixKey)
	cfg.TLS.CaPath = v.GetString(configPrefix + suffixCA)
	cfg.TLS.EnableHostVerification = v.GetBool(configPrefix + suffixVerifyHost)
	cfg.Brokers = strings.Split(v.GetString(configPrefix+suffixBrokers), ",")
	cfg.Version = v.GetString(configPrefix + suffixVersion)

	opt.config = cfg
	opt.topic = v.GetString(configPrefix + suffixTopic)
	opt.encoding = v.GetString(configPrefix + suffixEncoding)
}
