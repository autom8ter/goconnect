package fire

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/autom8ter/gosaas/sdk/go/proto/accounts"
	"github.com/autom8ter/gosaas/sdk/go/proto/contacts"
	"github.com/autom8ter/objectify"
	"github.com/pkg/errors"
)

var util = objectify.New()

type GoogleFireStore struct {
	fire *firestore.Client
}

func New(fire *firestore.Client) *GoogleFireStore {
	return &GoogleFireStore{
		fire: fire,
	}
}

func (g *GoogleFireStore) CreateAccount(ctx context.Context, r *accounts.CreateAccountRequest) error {
	notes := make(map[string]string)
	snap := g.fire.Collection("accounts").Doc(r.User.Email)
	acc := &accounts.Account{
		User: &contacts.Contact{
			Name:        r.User.Name,
			Email:       r.User.Email,
			Phone:       r.User.Phone,
			Address:     r.User.Address,
			Annotations: r.User.Annotations,
		},
		PasswordResetQuestion: r.ResetQuestion,
		PasswordResetAnswer:   r.ResetAnswer,
		Annotations:           notes,
	}
	pass, err := util.HashPassword(r.Password)
	if err != nil {
		return err
	}
	acc.HashedPassword = pass
	_, err = snap.Create(ctx, acc)
	if err != nil {
		return err
	}
	return nil
}

func (g *GoogleFireStore) ReadAccount(ctx context.Context, email string) (*accounts.Account, error) {
	snap, err := g.fire.Collection("accounts").Doc(email).Get(ctx)
	if err != nil {
		return nil, err
	}
	acc := &accounts.Account{}
	err = snap.DataTo(acc)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (g *GoogleFireStore) ListAccounts(ctx context.Context) ([]*accounts.Account, error) {
	coll := g.fire.Collection("accounts")
	snaps, err := coll.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}
	var accs = []*accounts.Account{}
	for _, snap := range snaps {
		thisAcc := &accounts.Account{}
		err = snap.DataTo(thisAcc)
		if err != nil {
			return accs, err
		}
	}
	return accs, nil
}

func (g *GoogleFireStore) UpdateAccount(ctx context.Context, a *accounts.Account) error {
	updates := []firestore.Update{}
	for k, v := range util.ToMap(a) {
		updates = append(updates, firestore.Update{
			Path:  k,
			Value: v,
		})
	}
	_, err := g.fire.Collection("accounts").Doc(a.User.Email).Update(ctx, updates)
	if err != nil {
		return err
	}
	return nil
}

func (g *GoogleFireStore) DeleteAccount(ctx context.Context, email string) error {
	_, err := g.fire.Collection("accounts").Doc(email).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (g *GoogleFireStore) AccountExists(ctx context.Context, email string) (bool, error) {
	acc := g.fire.Collection("accounts")
	docs := acc.Documents(ctx)
	allDocs, err := docs.GetAll()
	if err != nil {
		return false, err
	}
	for _, d := range allDocs {
		if d.Ref.ID == email {
			return true, nil
		}
	}
	return false, nil
}

func (g *GoogleFireStore) RecoverToken(ctx context.Context, email string) (string, error) {
	acc, err := g.ReadAccount(ctx, email)
	if err != nil {
		return "", err
	}
	return acc.RecoveryToken, nil
}

func (g *GoogleFireStore) ResetPassword(ctx context.Context, email, token, newPass string) error {
	acc, err := g.ReadAccount(ctx, email)
	if err != nil {
		return err
	}
	if acc.RecoveryToken != token {
		return errors.New("provided token is invalid")
	}
	hash, err := util.HashPassword(newPass)
	if err != nil {
		return err
	}
	acc.HashedPassword = hash
	err = g.UpdateAccount(ctx, acc)
	if err != nil {
		return err
	}
	return nil
}
