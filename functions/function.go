package functions

import (
	"github.com/autom8ter/goconnect"
	"github.com/pkg/errors"
)

type Func func(g *goconnect.GoConnect) error

func Execute(g *goconnect.GoConnect, fns ...Func) error {
	var err error
	for _, f := range fns {
		if newErr := f(g); newErr != nil {
			err = errors.Wrap(err, newErr.Error())
		}
	}
	return err
}
