package server

import (
	"cloud.google.com/go/firestore"
	strage "cloud.google.com/go/storage"
	"context"
	"github.com/autom8ter/engine"
	"github.com/autom8ter/engine/config"
	"github.com/autom8ter/engine/driver"
	"github.com/autom8ter/goconnect"
	"github.com/autom8ter/goconnect/pkg/accounts"
	"github.com/autom8ter/goconnect/pkg/errors"
	"github.com/autom8ter/goconnect/pkg/store"
	accountspb "github.com/autom8ter/gosaas/sdk/go/proto/accounts"
	"github.com/autom8ter/gosaas/sdk/go/proto/contacts"
	"github.com/autom8ter/gosaas/sdk/go/proto/storage"
	"github.com/autom8ter/objectify"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sfreiberg/gotwilio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"os"
)

func init() {
	util = objectify.New()
}

// Global variables//////////////
var ctx context.Context
var util *objectify.Handler
var fire *firestore.Client
var strg *strage.Client
var grid *sendgrid.Client
var twil *gotwilio.Twilio

type env struct {
	Vars []string `json:"vars"`
}

/////////////////////////////////

func New(c context.Context, credsFile, project, bucket string) *Server {
	ctx = context.WithValue(c, "env", util.ToMap(&env{os.Environ()}))
	conn := goconnect.NewFromFileEnv(credsFile)
	err := conn.Init()
	if err != nil {
		util.FatalErr(err, "failed to initialize goconnect instance")
	}
	fire, err = conn.GCP.Firestore(ctx, project)
	if err != nil {
		util.FatalErr(err, "failed to initialize firestore client")
	}
	strg, err = conn.GCP.Storage(ctx)
	if err != nil {
		util.FatalErr(err, "failed to initialize strorage client")
	}
	handle := strg.Bucket(bucket)
	serv := &Server{
		strg: store.NewGoogleCloudStorage(strg, handle, bucket),
	}
	serv.PluginFunc = func(s *grpc.Server) {
		accountspb.RegisterAccountServiceServer(s, serv)
		contacts.RegisterSMSServiceServer(s, serv)
		contacts.RegisterEmailServiceServer(s, serv)
		contacts.RegisterCallServiceServer(s, serv)
		storage.RegisterStorageServiceServer(s, serv)
	}
	return serv
}

type Server struct {
	strg *store.GoogleCloudStorage
	conn *goconnect.GoConnect
	driver.PluginFunc
}

func (s *Server) AccountToken(ctx context.Context, r *accountspb.RecoveryTokenRequest) (*accountspb.RecoveryTokenResponse, error) {
	resp, err := accounts.RecoverToken(ctx, r.Email, fire)
	if err != nil {
		return nil, errors.StatusStack(err, codes.Internal, "server.UpdateAccount: %s", "failed to update account")
	}
	return &accountspb.RecoveryTokenResponse{
		Token: resp,
	}, nil

}

func (s *Server) ResetPassword(ctx context.Context, r *accountspb.ResetPasswordRequest) (*accountspb.Empty, error) {
	err := accounts.ResetPassword(ctx, r.Email, r.Token, r.NewPassword, fire)
	if err != nil {
		return nil, errors.StatusStack(err, codes.Internal, "server.UpdateAccount: %s", "failed to update account")
	}
	return &accountspb.Empty{}, nil
}

func (s *Server) StoreObject(ctx context.Context, r *storage.StoreRequest) (*storage.StorageObject, error) {
	err := s.strg.Store(
		ctx,
		r.Filename,
		r.Data,
		map[string]string{},
	)

	if err != nil {
		return nil, errors.StatusStack(err, codes.Internal, "server.UpdateAccount: %s", "failed to update account")
	}

	return &storage.StorageObject{
		Filename: r.Filename,
		Url:      s.strg.PublicURL(r.Filename),
	}, nil
}

func (s *Server) DeleteObject(ctx context.Context, r *storage.DeleteRequest) (*storage.DeleteResponse, error) {
	err := s.strg.Delete(ctx, r.Filename)
	if err != nil {
		return nil, errors.StatusStack(err, codes.Internal, "server.UpdateAccount: %s", "failed to update account")
	}

	return &storage.DeleteResponse{
		Filename: r.Filename,
	}, nil
}

func (s *Server) SendEmail(ctx context.Context, e *contacts.Email) (*contacts.EmailResponse, error) {
	from := mail.NewEmail(e.From.Name, e.From.Email)
	to := mail.NewEmail(e.To.Name, e.To.Email)
	message := mail.NewSingleEmail(from, e.Subject, to, e.PlainText, e.HtmlAlternate)

	_, err := grid.Send(message)
	if err != nil {
		return nil, errors.StatusStack(err, codes.Internal, "server.UpdateAccount: %s", "failed to update account")
	}
	return &contacts.EmailResponse{}, nil
}

func (s *Server) SendMMS(ctx context.Context, r *contacts.MMS) (*contacts.SMSResponse, error) {
	_, ex, err := twil.SendMMS(r.Sms.From.Phone, r.Sms.To.Phone, r.Sms.Body, r.MediaUrl, "", "")
	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "server.SendMMS: failed to send mms")
	}
	return &contacts.SMSResponse{}, nil
}

func (s *Server) SendSMS(ctx context.Context, smS *contacts.SMS) (*contacts.SMSResponse, error) {
	_, ex, err := twil.SendSMS(smS.From.Phone, smS.To.Phone, smS.Body, "", "")
	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "server.SendSMS: failed to send sms")
	}
	return &contacts.SMSResponse{}, nil
}

func (s *Server) SendCall(ctx context.Context, c *contacts.Call) (*contacts.CallResponse, error) {
	_, ex, err := twil.CallWithUrlCallbacks(c.From.Phone, c.To.Phone, gotwilio.NewCallbackParameters(c.Callback))
	if err != nil {
		return nil, errors.TwilioError(err, codes.Internal, ex, "server.SendCall: failed to send call")
	}
	return &contacts.CallResponse{}, nil
}

func (s *Server) ListAccounts(ctx context.Context, r *accountspb.Empty) (*accountspb.ListAccountsResponse, error) {
	accs, err := accounts.List(ctx, fire)
	if err != nil {
		return nil, errors.StatusStack(err, codes.Internal, "server.ListAccounts: %s", "failed to update account")
	}
	return &accountspb.ListAccountsResponse{
		Accounts: accs,
	}, nil
}

func (s *Server) CreateAccount(ctx context.Context, r *accountspb.CreateAccountRequest) (*accountspb.Empty, error) {
	err := accounts.Create(ctx, r, fire)
	if err != nil {
		return nil, errors.StatusStack(err, codes.Internal, "server.CreateAccount: %s", "failed to update account")
	}
	return &accountspb.Empty{}, nil
}

func (s *Server) ReadAccount(ctx context.Context, r *accountspb.ReadAccountRequest) (*accountspb.Account, error) {
	return accounts.Read(ctx, r.Email, fire)
}

func (s *Server) UpdateAccount(ctx context.Context, a *accountspb.Account) (*accountspb.Empty, error) {
	err := accounts.Update(ctx, a, fire)
	if err != nil {
		return nil, errors.StatusStack(err, codes.Internal, "server.UpdateAccount: %s", "failed to update account")
	}
	return &accountspb.Empty{}, nil
}

func (s *Server) DeleteAccount(ctx context.Context, r *accountspb.DeleteAccountRequest) (*accountspb.Empty, error) {
	err := accounts.Delete(ctx, r.Email, fire)
	if err != nil {
		return nil, errors.StatusStack(err, codes.Internal, "server.DeleteAccount: %s", "failed to update account")
	}
	return &accountspb.Empty{}, nil
}

func (s *Server) AccountExists(ctx context.Context, r *accountspb.ExistsRequest) (*accountspb.ExistsResponse, error) {
	resp, err := accounts.Exists(ctx, r.Email, fire)
	if err != nil {
		return nil, errors.StatusStack(err, codes.Internal, "server.AccountExists: %s", "failed to update account")
	}
	return &accountspb.ExistsResponse{
		Exists: resp,
	}, nil
}

func (s *Server) Serve(network, addr string, debug bool) error {
	return s.
}
