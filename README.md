# goconnect
--
    import "github.com/autom8ter/goconnect"


## Usage

#### type GoConnect

```go
type GoConnect struct {
}
```


#### func  New

```go
func New(opts ...config.ConfigOption) *GoConnect
```

#### func (*GoConnect) Call

```go
func (g *GoConnect) Call(to, from, callback string) (*gotwilio.VoiceResponse, error)
```

#### func (*GoConnect) CallWithApp

```go
func (g *GoConnect) CallWithApp(to, from, appSid string) (*gotwilio.VoiceResponse, error)
```

#### func (*GoConnect) ChargeCustomer

```go
func (g *GoConnect) ChargeCustomer(opts ...pay.ChargeOption) ([]*stripe.Charge, error)
```

#### func (*GoConnect) CreateVideoRoom

```go
func (g *GoConnect) CreateVideoRoom() (*gotwilio.VideoResponse, error)
```

#### func (*GoConnect) Email

```go
func (g *GoConnect) Email(opts ...email.EmailOption) (*rest.Response, error)
```

#### func (*GoConnect) JSONString

```go
func (g *GoConnect) JSONString(obj interface{}) string
```

#### func (*GoConnect) MMS

```go
func (g *GoConnect) MMS(to, from, body, mediaURL string, callback, app string) (*gotwilio.SmsResponse, error)
```

#### func (*GoConnect) NewCustomer

```go
func (g *GoConnect) NewCustomer(opts ...customer.Option) (*stripe.Customer, error)
```

#### func (*GoConnect) SMS

```go
func (g *GoConnect) SMS(to, from, body, callback, app string) (*gotwilio.SmsResponse, error)
```

#### func (*GoConnect) SMSCopilot

```go
func (g *GoConnect) SMSCopilot(to, service, body, callback, app string) (*gotwilio.SmsResponse, error)
```
