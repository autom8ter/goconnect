package config

import "net/http"

type ConfigOption func(c *Config)

func NewConfig(opts ...ConfigOption) *Config {
	c := &Config{}
	for _, o := range opts {
		o(c)
	}
	return c
}

type Twilio struct {
	Account      string
	Token        string
	ClientKey    string
	ClientSecret string
}

type SendGrid struct {
	Account      string
	Token        string
	ClientKey    string
	ClientSecret string
}

type Stripe struct {
	Account      string
	Token        string
	ClientKey    string
	ClientSecret string
}

type Config struct {
	Client   *http.Client
	Twilio   *Twilio
	SendGrid *SendGrid
	Stripe   *Stripe
}
