package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/seriouspoop/gopush/model"
)

type Credentials struct {
	Username string
	Token    string
}

type Config struct {
	Auth struct {
		BitBucket *Credentials
		GitHub    *Credentials
		GitLab    *Credentials
	}

	DefaultRemote string
	BranchPrefix  string
}

func (c *Config) ProviderAuth(p model.Provider) *Credentials {
	providerToAuth := map[model.Provider]*Credentials{
		model.BITBUCKET: c.Auth.BitBucket,
		model.GITHUB:    c.Auth.GitHub,
		model.GITLAB:    c.Auth.GitLab,
	}
	return providerToAuth[p]
}

func Read(filename, path string) (*Config, error) {
	b, err := os.ReadFile(filepath.Join(path, filename))
	if err != nil {
		return nil, err
	}
	c := &Config{}
	err = toml.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Config) Write(filename, path string) error {
	b, err := toml.Marshal(c)
	if err != nil {
		return err
	}
	writeByte := fmt.Sprintf("# AUTO-GENERATED FILE BY GOPUSH\n# DO NOT EDIT\n%s", string(b))
	return os.WriteFile(filepath.Join(path, filename), []byte(writeByte), 0700)
}
