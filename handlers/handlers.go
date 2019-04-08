package handlers

import (
	"github.com/autom8ter/goconnect"
	"github.com/autom8ter/objectify"
)

var util = objectify.New()

func NooP() goconnect.HandlerFunc {
	return func(g *goconnect.GoConnect) error {
		util.Debug("noop handler", "key", "value")
		return nil
	}
}
