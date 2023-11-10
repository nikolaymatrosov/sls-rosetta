package examples

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

type Language struct {
	Name  string `yaml:"name" json:"name"`
	Title string `yaml:"title" json:"title"`
}

type Deploy struct {
	Type      string   `yaml:"type" json:"type"`
	Exclusive []string `yaml:"exclusive" json:"exclusive,omitempty"`
}

const TerraformDeployType = "terraform"
const YcCliDeployType = "yccli"

type Example struct {
	Name        string   `yaml:"name" json:"name,omitempty"`
	Title       string   `yaml:"title" json:"title,omitempty"`
	Description string   `yaml:"description" json:"description,omitempty"`
	Deploy      []Deploy `yaml:"deploy" json:"deploy,omitempty"`
}

type Config struct {
	Repo      string               `yaml:"repo" json:"repo"`
	Languages []Language           `yaml:"languages" json:"languages"`
	Examples  map[string][]Example `yaml:"examples" json:"examples"`
}

func NewConfig() *Config {
	return &Config{}
}

func DefaultConfigUrl() string {
	return "https://raw.githubusercontent.com/nikolaymatrosov/sls-rosetta/main/examples/examples.yaml"
}

func (c *Config) Fetch(configUrl string) error {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		ForceContentType("application/json").
		Get(configUrl)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal(resp.Body(), c)

	return nil
}

func FromContext(ctx context.Context) (*Config, error) {
	if c, ok := ctx.Value("config").(*Config); ok {
		return c, nil
	}
	return nil, fmt.Errorf("config not found in context")
}

func ContextWith(ctx context.Context, c *Config) context.Context {
	return context.WithValue(ctx, "config", c)
}
