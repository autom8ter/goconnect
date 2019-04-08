package client

import (
	"context"
	"github.com/autom8ter/gosaas/sdk/go/proto/accounts"
	"github.com/autom8ter/gosaas/sdk/go/proto/contacts"
	"github.com/autom8ter/gosaas/sdk/go/proto/storage"
	"google.golang.org/grpc"
	"log"
)

type Client struct {
	acc   accounts.AccountServiceClient
	storj storage.StorageServiceClient
	sms   contacts.SMSServiceClient
	call  contacts.CallServiceClient
	mail  contacts.EmailServiceClient
}

func New(ctx context.Context, addr string, opts ...grpc.DialOption) *Client {
	conn, err := grpc.DialContext(ctx, addr, opts...)
	if err != nil {
		log.Fatalln(err.Error())
	}
	return &Client{
		acc:   accounts.NewAccountServiceClient(conn),
		storj: storage.NewStorageServiceClient(conn),
		sms:   contacts.NewSMSServiceClient(conn),
		call:  contacts.NewCallServiceClient(conn),
		mail:  contacts.NewEmailServiceClient(conn),
	}
}

func (c *Client) Accounts() accounts.AccountServiceClient {
	return c.acc
}

func (c *Client) Storage() storage.StorageServiceClient {
	return c.storj
}

func (c *Client) SMS() contacts.SMSServiceClient {
	return c.sms
}

func (c *Client) Call() contacts.CallServiceClient {
	return c.call
}

func (c *Client) Email() contacts.EmailServiceClient {
	return c.mail
}
