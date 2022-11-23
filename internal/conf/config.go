// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package conf

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var ImageReplacer = struct {
	TlsCrtFile     string            `yaml:"tls_crt_file"`
	TlsKeyFile     string            `yaml:"tls_key_file"`
	DockerRegistry string            `yaml:"docker_registry"`
	ReplacePolicy  map[string]string `yaml:"replace_policy"`
}{}

func Init(configFile string) error {
	fileBytes, err := os.ReadFile(configFile)
	if err != nil {
		return errors.Wrap(err, "read file")
	}

	if err := yaml.Unmarshal(fileBytes, &ImageReplacer); err != nil {
		return errors.Wrap(err, "unmarshal")
	}
	return nil
}
