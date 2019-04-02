package customer

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
)

type Option func(params *stripe.CustomerParams)

func New(opts ... Option) (*stripe.Customer, error) {
	c := &stripe.CustomerParams{}
	for _, o := range opts {
		o(c)
	}
	return customer.New(c)
}