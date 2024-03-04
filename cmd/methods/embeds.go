// file: ./cmd/methods/createfunctions.go
package methods

import (
	_ "embed"
)

var (

	//go:embed "docker-compose/docker-compose.yaml"
	dockercompose []byte

	//go:embed "configurations/env.env"
	configurations []byte
)

func GetConfigurationsEmbed() []byte {
	return configurations
}

func GetDockerComposeEmbed() []byte {
	return dockercompose
}
