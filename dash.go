package goconnect

import (
	"github.com/stripe/stripe-go"
	"io"
)

type CustomerInfo struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Email       string            `json:"email"`
	Phone       string            `json:"phone"`
	Plans       []*Plan           `json:"plan"`
	Annotations map[string]string `json:"annotations"`
}

type Plan struct {
	Name    string `json:"id"`
	Amount  int64  `json:"amount"`
	Active  bool   `json:"active"`
	Service string `json:"service"`
}

func (g *GoConnect) Compile(c *CustomerInfo, hTML string, w io.Writer) error {
	return g.util.RenderHTML(hTML, c, w)
}

func (g *GoConnect) ToCustomer(c *CustomerInfo) (*stripe.Customer, error) {

	switch g.cfg.Index {
	case EMAIL:
		cust, ok := g.GetCustomer(c.Email)
		if !ok {
			return nil, NOEXIST(c.Email)
		}
		return cust, nil
	case PHONE:
		cust, ok := g.GetCustomer(c.Phone)
		if !ok {
			return nil, NOEXIST(c.Phone)
		}
		return cust, nil
	default:
		cust, ok := g.GetCustomer(c.Id)
		if !ok {
			return nil, NOEXIST(c.Id)
		}
		return cust, nil
	}
}

func (g *GoConnect) ToCustomerInfo(customerKey string) (*CustomerInfo, error) {
	cust, ok := g.GetCustomer(customerKey)
	if !ok {
		return nil, NOEXIST(customerKey)
	}
	plans := []*Plan{}
	for _, data := range cust.Subscriptions.Data {
		plans = append(plans, &Plan{
			Name:    data.Plan.Nickname,
			Amount:  data.Plan.Amount,
			Active:  data.Plan.Active,
			Service: data.Plan.Product.Name,
		})
	}

	return &CustomerInfo{
		Id:          cust.ID,
		Name:        cust.Shipping.Name,
		Email:       cust.Email,
		Phone:       cust.Shipping.Phone,
		Plans:       plans,
		Annotations: nil,
	}, nil
}
