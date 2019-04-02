package goconnect_test

import (
	"github.com/autom8ter/goconnect"
	"log"
	"testing"
)

func init() {
	g = goconnect.New(nil)

}

var g *goconnect.GoConnect

func TestNewFromEnv(t *testing.T) {
	if g == nil {
		log.Fatalln("nil goconnect")
	}
}
