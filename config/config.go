// Copyright (c) 2020 IoTeX Foundation
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package config

import (
	"github.com/pkg/errors"
	uconfig "go.uber.org/config"
)

type (
	NetworkIdentifier struct {
		Blockchain   string `yaml:"blockchain"`
		Network      string `yaml:"network"`
		EvmNetworkID uint32 `yaml:"evmNetworkID"`
	}
	Currency struct {
		Symbol   string `yaml:"symbol"`
		Decimals int32  `yaml:"decimals"`
	}
	Server struct {
		Port           string `yaml:"port"`
		Endpoint       string `yaml:"endpoint"`
		SecureEndpoint bool   `yaml:"secureEndpoint"`
		RosettaVersion string `yaml:"rosettaVersion"`
	}
	Config struct {
		NetworkIdentifier NetworkIdentifier `yaml:"network_identifier"`
		Currency          Currency          `yaml:"currency"`
		Server            Server            `yaml:"server"`
		KeepNoneTxAction  bool              `yaml:"keepNoneTxAction"`
	}
)

func New(path string) (cfg *Config, err error) {
	opts := []uconfig.YAMLOption{uconfig.File(path)}
	yaml, err := uconfig.NewYAML(opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init config")
	}
	cfg = &Config{}
	if err := yaml.Get(uconfig.Root).Populate(cfg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal YAML config to struct")
	}
	return
}
