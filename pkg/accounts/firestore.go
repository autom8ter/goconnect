package accounts

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/autom8ter/gosaas/pkg/util"
	"github.com/autom8ter/gosaas/sdk/go/proto/accounts"
	"github.com/autom8ter/gosaas/sdk/go/proto/contacts"
	"github.com/pkg/errors"
)

func Create(ctx context.Context, r *accounts.CreateAccountRequest, client *firestore.Client) error {
	notes := make(map[string]string)
	snap := client.Collection("accounts").Doc(r.User.Email)
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

func Read(ctx context.Context, email string, client *firestore.Client) (*accounts.Account, error) {
	snap, err := client.Collection("accounts").Doc(email).Get(ctx)
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

func List(ctx context.Context, client *firestore.Client) ([]*accounts.Account, error) {
	coll := client.Collection("accounts")
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

func Update(ctx context.Context, a *accounts.Account, client *firestore.Client) error {
	updates := []firestore.Update{}
	for k, v := range util.AsMap(a) {
		updates = append(updates, firestore.Update{
			Path:  k,
			Value: v,
		})
	}
	_, err := client.Collection("accounts").Doc(a.User.Email).Update(ctx, updates)
	if err != nil {
		return err
	}
	return nil
}

func Delete(ctx context.Context, email string, client *firestore.Client) error {
	_, err := client.Collection("accounts").Doc(email).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func Exists(ctx context.Context, email string, client *firestore.Client) (bool, error) {
	acc := client.Collection("accounts")
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

func RecoverToken(ctx context.Context, email string, client *firestore.Client) (string, error) {
	acc, err := Read(ctx, email, client)
	if err != nil {
		return "", err
	}
	return acc.RecoveryToken, nil
}

func ResetPassword(ctx context.Context, email, token, newPass string, client *firestore.Client) error {
	acc, err := Read(ctx, email, client)
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
	err = Update(ctx, acc, client)
	if err != nil {
		return err
	}
	return nil
}
