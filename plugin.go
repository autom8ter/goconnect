package goconnect

import (
	"github.com/autom8ter/engine"
	"github.com/autom8ter/engine/driver"
)

type Plugin interface {
	driver.Plugin
	RegisterWithClient(g *GoConnect)
}

func (g *GoConnect) Serve(addr string, plugs ...Plugin) error {
	plugins := []driver.Plugin{}
	for _, p := range plugs {
		p.RegisterWithClient(g)
		plugins = append(plugins, p)
	}
	return engine.Serve(addr, true, plugins...)
}
