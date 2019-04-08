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

func (s *Stripe) UpdatePlan(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	return s.client.Plans.Update(id, params)
}

func (s *Stripe) DeletePlan(id string, params *stripe.PlanParams) (*stripe.Plan, error) {
	return s.client.Plans.Del(id, params)
}

func (s *Stripe) NewPlan(params *stripe.PlanParams) (*stripe.Plan, error) {
	return s.client.Plans.New(params)
}

func (s *Stripe) GetCustomer(id string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	return s.client.Customers.Get(id, params)
}

func (s *Stripe) UpdateCustomer(id string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	return s.client.Customers.Update(id, params)
}

func (s *Stripe) DeleteCustomer(id string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	return s.client.Customers.Del(id, params)
}

func (s *Stripe) NewCustomer(params *stripe.CustomerParams) (*stripe.Customer, error) {
	return s.client.Customers.New(params)
}

func (s *Stripe) GetCharge(id string, params *stripe.ChargeParams) (*stripe.Charge, error) {
	return s.client.Charges.Get(id, params)
}

func (s *Stripe) NewCharge(id string, params *stripe.ChargeParams) (*stripe.Charge, error) {
	return s.client.Charges.New(params)
}

func (s *Stripe) UpdateCharge(id string, params *stripe.ChargeParams) (*stripe.Charge, error) {
	return s.client.Charges.Update(id, params)
}

func (s *Stripe) GetSubscription(id string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return s.client.Subscriptions.Get(id, params)
}

func (s *Stripe) NewSubscription(params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return s.client.Subscriptions.New(params)
}

func (s *Stripe) UpdateSubscription(id string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return s.client.Subscriptions.Update(id, params)
}

func (s *Stripe) DeleteSubscription(id string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return s.client.Subscriptions.Update(id, params)
}

func (s *Stripe) GetInvoice(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return s.client.Invoices.Get(id, params)
}

func (s *Stripe) NewInvoice(params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return s.client.Invoices.New(params)
}

func (s *Stripe) UpdateInvoice(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return s.client.Invoices.Update(id, params)
}

func (s *Stripe) DeleteInvoice(id string, params *stripe.InvoiceParams) (*stripe.Invoice, error) {
	return s.client.Invoices.Del(id, params)
}

func (s *Stripe) GetAccount() (*stripe.Account, error) {
	return s.client.Account.Get()
}

func (s *Stripe) GetBankAccount(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return s.client.BankAccounts.Get(id, params)
}

func (s *Stripe) NewBankAccount(params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return s.client.BankAccounts.New(params)
}

func (s *Stripe) UpdateBankAccount(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return s.client.BankAccounts.Update(id, params)
}

func (s *Stripe) DeleteBankAccount(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return s.client.BankAccounts.Del(id, params)
}

func (s *Stripe) GetProduct(id string, params *stripe.ProductParams) (*stripe.Product, error) {
	return s.client.Products.Get(id, params)
}

func (s *Stripe) UpdateProduct(id string, params *stripe.ProductParams) (*stripe.Product, error) {
	return s.client.Products.Update(id, params)
}

func (s *Stripe) DeleteProduct(id string, params *stripe.ProductParams) (*stripe.Product, error) {
	return s.client.Products.Del(id, params)
}

func (s *Stripe) NewProduct(params *stripe.ProductParams) (*stripe.Product, error) {
	return s.client.Products.New(params)
}

func (s *Stripe) Recipients(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	return s.client.Recipients.Get(id, params)
}

func (s *Stripe) UpdateRecipients(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	return s.client.Recipients.Update(id, params)
}

func (s *Stripe) DeleteRecipients(id string, params *stripe.RecipientParams) (*stripe.Recipient, error) {
	return s.client.Recipients.Del(id, params)
}

func (s *Stripe) GetOrder(id string, params *stripe.OrderParams) (*stripe.Order, error) {
	return s.client.Orders.Get(id, params)
}

func (s *Stripe) UpdateOrder(id string, params *stripe.OrderUpdateParams) (*stripe.Order, error) {
	return s.client.Orders.Update(id, params)
}

func (s *Stripe) PayOrder(id string, params *stripe.OrderPayParams) (*stripe.Order, error) {
	return s.client.Orders.Pay(id, params)
}

func (s *Stripe) ReturnOrder(id string, params *stripe.OrderReturnParams) (*stripe.OrderReturn, error) {
	return s.client.Orders.Return(id, params)
}

func (s *Stripe) GetPayout(id string, params *stripe.PayoutParams) (*stripe.Payout, error) {
	return s.client.Payouts.Get(id, params)
}

func (s *Stripe) NewPayout(params *stripe.PayoutParams) (*stripe.Payout, error) {
	return s.client.Payouts.New(params)
}

func (s *Stripe) UpdatePayout(id string, params *stripe.PayoutParams) (*stripe.Payout, error) {
	return s.client.Payouts.Update(id, params)
}
