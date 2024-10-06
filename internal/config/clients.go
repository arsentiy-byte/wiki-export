package config

import "wiki-export/internal/config/clients"

type Clients struct {
	Http *clients.Http `yaml:"http" env-required:"true"`
}
