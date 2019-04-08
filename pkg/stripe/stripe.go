package stripe

import (
	"github.com/stripe/stripe-go/balance"
	"github.com/stripe/stripe-go/charge"
	"github.com/stripe/stripe-go/client"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/invoice"
	"github.com/stripe/stripe-go/plan"
	"github.com/stripe/stripe-go/product"
	"github.com/stripe/stripe-go/sub"
)

type Stripe struct {
	client *client.API
}

func New(client *client.API) *Stripe {
	return &Stripe{
		client: client,
	}
}

func (s *Stripe) Charges() *charge.Client {
	return s.client.Charges
}

func (s *Stripe) Invoices() *invoice.Client {
	return s.client.Invoices
}

func (s *Stripe) Customer() *customer.Client {
	return s.client.Customers
}

func (s *Stripe) Subscriptions() *sub.Client {
	return s.client.Subscriptions
}

func (s *Stripe) Products() *product.Client {
	return s.client.Products
}

func (s *Stripe) Plans() *plan.Client {
	return s.client.Plans
}

func (s *Stripe) Balance() *balance.Client {
	return s.client.Balance
}