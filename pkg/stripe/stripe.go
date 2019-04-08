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

type BalanceOption func(params *stripe.BalanceParams)

func (s *Stripe) GetBalance(opts ...BalanceOption) (*stripe.Balance, error) {
	var p = &stripe.BalanceParams{}
	for _, o := range opts {
		o(p)
	}
	return s.client.Balance.Get(p)
}

type PlanOption func(params *stripe.PlanParams)

func (s *Stripe) GetPlan(id string, opts ...PlanOption) (*stripe.Plan, error) {
	var p = &stripe.PlanParams{}
	for _, o := range opts {
		o(p)
	}
	return s.client.Plans.Get(id, p)
}

func (s *Stripe) UpdatePlan(id string, opts ...PlanOption) (*stripe.Plan, error) {
	var params = &stripe.PlanParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Plans.Update(id, params)
}

func (s *Stripe) DeletePlan(id string, opts ...PlanOption) (*stripe.Plan, error) {
	var params = &stripe.PlanParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Plans.Del(id, params)
}

func (s *Stripe) NewPlan(id string, opts ...PlanOption) (*stripe.Plan, error) {
	var params = &stripe.PlanParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Plans.New(params)
}

type CustomerOption func(params *stripe.CustomerParams)

func (s *Stripe) GetCustomer(id string, opts ...CustomerOption) (*stripe.Customer, error) {
	var params = &stripe.CustomerParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Customers.Get(id, params)
}

func (s *Stripe) UpdateCustomer(id string, opts ...CustomerOption) (*stripe.Customer, error) {
	var params = &stripe.CustomerParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Customers.Update(id, params)
}

func (s *Stripe) DeleteCustomer(id string, opts ...CustomerOption) (*stripe.Customer, error) {
	var params = &stripe.CustomerParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Customers.Del(id, params)
}

func (s *Stripe) NewCustomer(id string, opts ...CustomerOption) (*stripe.Customer, error) {
	var params = &stripe.CustomerParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Customers.New(params)
}

type ChargeOption func(params *stripe.ChargeParams)

func (s *Stripe) GetCharge(id string, opts ...ChargeOption) (*stripe.Charge, error) {
	var params = &stripe.ChargeParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Charges.Get(id, params)
}

func (s *Stripe) NewCharge(id string, opts ...ChargeOption) (*stripe.Charge, error) {
	var params = &stripe.ChargeParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Charges.New(params)
}

func (s *Stripe) UpdateCharge(id string, opts ...ChargeOption) (*stripe.Charge, error) {
	var params = &stripe.ChargeParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Charges.Update(id, params)
}

type SubscriptionOption func(params *stripe.SubscriptionParams)

func (s *Stripe) GetSubscription(id string, opts ...SubscriptionOption) (*stripe.Subscription, error) {
	var params = &stripe.SubscriptionParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Subscriptions.Get(id, params)
}

func (s *Stripe) NewSubscription(opts ...SubscriptionOption) (*stripe.Subscription, error) {
	var params = &stripe.SubscriptionParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Subscriptions.New(params)
}

func (s *Stripe) UpdateSubscription(id string, opts ...SubscriptionOption) (*stripe.Subscription, error) {
	var params = &stripe.SubscriptionParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Subscriptions.Update(id, params)
}

func (s *Stripe) DeleteSubscription(id string, opts ...SubscriptionOption) (*stripe.Subscription, error) {
	var params = &stripe.SubscriptionParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Subscriptions.Update(id, params)
}

type InvoiceOption func(params *stripe.InvoiceParams)

func (s *Stripe) GetInvoice(id string, opts ...InvoiceOption) (*stripe.Invoice, error) {
	var params = &stripe.InvoiceParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Invoices.Get(id, params)
}

func (s *Stripe) NewInvoice(opts ...InvoiceOption) (*stripe.Invoice, error) {
	var params = &stripe.InvoiceParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Invoices.New(params)
}

func (s *Stripe) UpdateInvoice(id string, opts ...InvoiceOption) (*stripe.Invoice, error) {
	var params = &stripe.InvoiceParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Invoices.Update(id, params)
}

func (s *Stripe) DeleteInvoice(id string, opts ...InvoiceOption) (*stripe.Invoice, error) {
	var params = &stripe.InvoiceParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Invoices.Del(id, params)
}

func (s *Stripe) GetAccount() (*stripe.Account, error) {
	return s.client.Account.Get()
}

type BankAccountOption func(params *stripe.BankAccountParams)

func (s *Stripe) GetBankAccount(id string, opts ...BankAccountOption) (*stripe.BankAccount, error) {
	var params = &stripe.BankAccountParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.BankAccounts.Get(id, params)
}

func (s *Stripe) NewBankAccount(opts ...BankAccountOption) (*stripe.BankAccount, error) {
	var params = &stripe.BankAccountParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.BankAccounts.New(params)
}

func (s *Stripe) UpdateBankAccount(id string, opts ...BankAccountOption) (*stripe.BankAccount, error) {
	var params = &stripe.BankAccountParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.BankAccounts.Update(id, params)
}

func (s *Stripe) DeleteBankAccount(id string, opts ...BankAccountOption) (*stripe.BankAccount, error) {
	var params = &stripe.BankAccountParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.BankAccounts.Del(id, params)
}

type ProductOption func(params *stripe.ProductParams)

func (s *Stripe) GetProduct(id string, opts ...ProductOption) (*stripe.Product, error) {
	var params = &stripe.ProductParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Products.Get(id, params)
}

func (s *Stripe) UpdateProduct(id string, opts ...ProductOption) (*stripe.Product, error) {
	var params = &stripe.ProductParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Products.Update(id, params)
}

func (s *Stripe) DeleteProduct(id string, opts ...ProductOption) (*stripe.Product, error) {
	var params = &stripe.ProductParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Products.Del(id, params)
}

func (s *Stripe) NewProduct(opts ...ProductOption) (*stripe.Product, error) {
	var params = &stripe.ProductParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Products.New(params)
}

type ReciptientOption func(params *stripe.RecipientParams)

func (s *Stripe) GetRecipient(id string, opts ...ReciptientOption) (*stripe.Recipient, error) {
	var params = &stripe.RecipientParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Recipients.Get(id, params)
}

func (s *Stripe) UpdateRecipients(id string, opts ...ReciptientOption) (*stripe.Recipient, error) {
	var params = &stripe.RecipientParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Recipients.Update(id, params)
}

func (s *Stripe) DeleteRecipients(id string, opts ...ReciptientOption) (*stripe.Recipient, error) {
	var params = &stripe.RecipientParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Recipients.Del(id, params)
}

type OrderOption func(params *stripe.OrderParams)

func (s *Stripe) GetOrder(id string, opts ...OrderOption) (*stripe.Order, error) {
	var params = &stripe.OrderParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Orders.Get(id, params)
}

type OrderUpdateOption func(params *stripe.OrderUpdateParams)

func (s *Stripe) UpdateOrder(id string, opts ...OrderUpdateOption) (*stripe.Order, error) {
	var params = &stripe.OrderUpdateParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Orders.Update(id, params)
}

type OrderPayOption func(params *stripe.OrderPayParams)

func (s *Stripe) PayOrder(id string, opts ...OrderPayOption) (*stripe.Order, error) {
	var params = &stripe.OrderPayParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Orders.Pay(id, params)
}

type OrderReturnOption func(params *stripe.OrderReturnParams)

func (s *Stripe) ReturnOrder(id string, opts ...OrderReturnOption) (*stripe.OrderReturn, error) {
	var params = &stripe.OrderReturnParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Orders.Return(id, params)
}

type PayoutOption func(params *stripe.PayoutParams)

func (s *Stripe) GetPayout(id string, opts ...PayoutOption) (*stripe.Payout, error) {
	var params = &stripe.PayoutParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Payouts.Get(id, params)
}

func (s *Stripe) NewPayout(opts ...PayoutOption) (*stripe.Payout, error) {
	var params = &stripe.PayoutParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Payouts.New(params)
}

func (s *Stripe) UpdatePayout(id string, opts ...PayoutOption) (*stripe.Payout, error) {
	var params = &stripe.PayoutParams{}
	for _, o := range opts {
		o(params)
	}
	return s.client.Payouts.Update(id, params)
}
