package goconnect

import (
	"github.com/autom8ter/api/go/api"
	"github.com/sfreiberg/gotwilio"
	"github.com/stripe/stripe-go"
)

func (g *GoConnect) customerKeys(m map[string]*stripe.Customer) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func (g *GoConnect) merge(ex *gotwilio.Exception, err error) error {
	if err != nil && ex != nil {
		return api.Util.WrapErr(err, string(api.Util.MarshalJSON(ex)))
	}
	if err != nil {
		return err
	}
	return nil
}
