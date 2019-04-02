package config

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
)

type Option func(c *Config)

func NewConfig(opts ...Option) *Config {
	c := &Config{}
	for _, o := range opts {
		o(c)
	}
	return c
}

type Twilio struct {
	Account string
	Token   string
}

type SendGrid struct {
	Account string
	Token   string
}

type Stripe struct {
	Account string
	Token   string
}

type Config struct {
	Context           context.Context
	Debug             bool
	FirebaseCredsFile string
	Twilio            *Twilio
	SendGrid          *SendGrid
	Stripe            *Stripe
	Client            *http.Client
}

func FromEnv(cli *http.Client) Option {
	return func(c *Config) {
		if !exists("credentials.json") {
			log.Fatalln("please place your firebase credentials in $PWD/credentials.json")
		}
		if os.Getenv("TWILIO_ACCOUNT") == "" {
			log.Fatalln("TWILIO_ACCOUNT is required")
		}
		if os.Getenv("TWILIO_TOKEN") == "" {
			log.Fatalln("TWILIO_TOKEN is required")
		}
		if os.Getenv("SENDGRID_TOKEN") == "" {
			log.Fatalln("SENDGRID_TOKEN is required")
		}
		if os.Getenv("STRIPE_TOKEN") == "" {
			log.Fatalln("STRIPE_TOKEN is required")
		}

		var err error
		c.Client = cli
		if err != nil {
			log.Fatalf("failed to create google default client: %s\n", err.Error())
		}
		c.Twilio.Account = os.Getenv("TWILIO_ACCOUNT")
		c.Twilio.Token = os.Getenv("TWILIO_TOKEN")
		c.SendGrid.Account = os.Getenv("SENDGRID_ACCOUNT")
		c.SendGrid.Token = os.Getenv("SENDGRID_TOKEN")
		c.Stripe.Account = os.Getenv("STRIPE_ACCOUNT")
		c.Stripe.Token = os.Getenv("STRIPE_TOKEN")
		if strings.Contains(os.Getenv("DEBUG"), "y") || strings.Contains(os.Getenv("DEBUG"), "Y") || strings.Contains(os.Getenv("DEBUG"), "t") || strings.Contains(os.Getenv("DEBUG"), "T") {
			c.Debug = true
		}
		log.Println("searching for google credentials in file: ", "credentials.json")
		c.FirebaseCredsFile = "credentials.json"
		if c.Context == nil {
			c.Context = context.TODO()
		}
	}
}

// exists checks if a file or directory exists.
func exists(path string) bool {
	if path == "" {
		return false
	}
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if !os.IsNotExist(err) {
		log.Fatalln(err.Error())
	}
	return false
}
