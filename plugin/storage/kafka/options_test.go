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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jaegertracing/jaeger/pkg/config"
)

func TestOptionsWithFlags(t *testing.T) {
	opts := &Options{}
	v, command := config.Viperize(opts.AddFlags)
	command.ParseFlags([]string{
		"--kafka.topic=topic1",
		"--kafka.brokers=127.0.0.1:9092,0.0.0:1234",
		"--kafka.encoding=protobuf",
		"--kafka.sasl=true",
		"--kafka.sasl.username=username",
		"--kafka.sasl.password=password",
		"--kafka.metadata.full=false",
		"--kafka.tls=true",
		"--kafka.tls.cert=cert",
		"--kafka.tls.key=key",
		"--kafka.tls.ca=ca",
		"--kafka.tls.verify-host=true",
		"--kafka.version=1.0.2"})
	opts.InitFromViper(v)

	assert.Equal(t, "topic1", opts.topic)
	assert.Equal(t, []string{"127.0.0.1:9092", "0.0.0:1234"}, opts.config.Brokers)
	assert.Equal(t, "protobuf", opts.encoding)
	assert.Equal(t, true, opts.config.SASL.Enabled)
	assert.Equal(t, "username", opts.config.SASL.Username)
	assert.Equal(t, "password", opts.config.SASL.Password)
	assert.Equal(t, false, opts.config.MetadataFull)
	assert.Equal(t, true, opts.config.TLS.Enabled)
	assert.Equal(t, "cert", opts.config.TLS.CertPath)
	assert.Equal(t, "key", opts.config.TLS.KeyPath)
	assert.Equal(t, "ca", opts.config.TLS.CaPath)
	assert.Equal(t, true, opts.config.TLS.EnableHostVerification)
	assert.Equal(t, "1.0.2", opts.config.Version)
}

func TestFlagDefaults(t *testing.T) {
	opts := &Options{}
	v, command := config.Viperize(opts.AddFlags)
	command.ParseFlags([]string{})
	opts.InitFromViper(v)

	assert.Equal(t, defaultTopic, opts.topic)
	assert.Equal(t, []string{defaultBroker}, opts.config.Brokers)
	assert.Equal(t, defaultEncoding, opts.encoding)
	assert.Equal(t, false, opts.config.SASL.Enabled)
	assert.Equal(t, "", opts.config.SASL.Username)
	assert.Equal(t, "", opts.config.SASL.Password)
	assert.Equal(t, true, opts.config.MetadataFull)
	assert.Equal(t, false, opts.config.TLS.Enabled)
	assert.Equal(t, "", opts.config.TLS.CertPath)
	assert.Equal(t, "", opts.config.TLS.KeyPath)
	assert.Equal(t, "", opts.config.TLS.CaPath)
	assert.Equal(t, false, opts.config.TLS.EnableHostVerification)
	assert.Equal(t, "", opts.config.Version)
}
