package stripe

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
)

type Stripe struct {
	client *client.API
}

func New(client *client.API) *Stripe {
	return &Stripe{
		client: client,
	}
}

func (s *Stripe) GetBalance(params *stripe.BalanceParams) (*stripe.Balance, error) {
	return s.client.Balance.Get(params)
}

func (s *Stripe) GetPlan(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	return s.client.Plans.Get(id, params)
}

func (s *Stripe) GetCustomer(id string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	return s.client.Customers.Get(id, params)
}

func (s *Stripe) GetCharge(id string, params *stripe.ChargeParams) (*stripe.Charge, error) {
	return s.client.Charges.Get(id, params)
}

func (s *Stripe) GetSubscription(id string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return s.client.Subscriptions.Get(id, params)
}

func (s *Stripe) GetInvoice(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return s.client.Invoices.Get(id, params)
}

func (s *Stripe) GetAccount() (*stripe.Account, error) {
	return s.client.Account.Get()
}

func (s *Stripe) GetBankAccount(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return s.client.BankAccounts.Get(id, params)
}

func (s *Stripe) GetProduct(id string, params *stripe.ProductParams) (*stripe.Product, error) {
	return s.client.Products.Get(id, params)
}

func (s *Stripe) Recipients(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	return s.client.Recipients.Get(id, params)
}

func (s *Stripe) GetOrder(id string, params *stripe.OrderParams) (*stripe.Order, error) {
	return s.client.Orders.Get(id, params)
}

func (s *Stripe) GetPayout(id string, params *stripe.PayoutParams) (*stripe.Payout, error) {
	return s.client.Payouts.Get(id, params)
}
