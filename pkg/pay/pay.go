package pay

import "github.com/stripe/stripe-go"

type ChargeOption func(params *stripe.ChargeListParams)

func NewCharge(opts ...ChargeOption) *stripe.ChargeListParams {
	params := &stripe.ChargeListParams{}
	for _, o := range opts {
		o(params)
	}
	return params
}
