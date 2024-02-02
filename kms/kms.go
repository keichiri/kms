package kms

import (
	"divine/kms/keys"
	"github.com/pkg/errors"
	"os/exec"
	"strconv"
	"strings"
)

type KMS struct {
	keyManager *keys.KeysManager
}

func NewKMS() *KMS {
	keyManager := keys.NewKeysManager()
	return &KMS{keyManager: keyManager}
}

func (k *KMS) IssueRewards(amounts []uint64, userEmails []string) error {
	users := make([]*User, 0)

	for _, userEmail := range userEmails {
		user, err := k.GetUser(userEmail)
		if err != nil {
			return errors.Wrapf(err, "Failed to retrieve user %s", userEmail)
		}

		if user == nil {
			user, err = k.CreateUser(userEmail)
			if err != nil {
				return errors.Wrapf(err, "Failed to create user %s", userEmail)
			}
		}

		users = append(users, user)
	}

	return k.issueRewards(users, amounts)
}

func (k *KMS) issueRewards(users []*User, amounts []uint64) error {
	userAddresses := make([]string, 0)
	for _, user := range users {
		userAddresses = append(userAddresses, user.Address)
	}
	userAddressesStr := strings.Join(userAddresses, ",")

	amountStrings := make([]string, 0)
	for _, amount := range amounts {
		amountStrings = append(amountStrings, strconv.FormatUint(amount, 10))
	}
	amountsStr := strings.Join(amountStrings, ",")

	cmd := exec.Command("divined", "tx", "rewards", "issue", userAddressesStr, amountsStr, "-y", "--from", "alice", "--chain-id", "divine")
	if _, err := cmd.Output(); err != nil {
		return errors.Wrap(err, "Failed to execute CLI command")
	}

	return nil
}

func (k *KMS) GetUser(userEmail string) (*User, error) {
	return nil, nil
}

func (k *KMS) CreateUser(userEmail string) (*User, error) {
	kp, err := k.keyManager.CreateNewKeypair(userEmail)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to create new keypair")
	}

	return &User{
		Email:   userEmail,
		Address: kp.Address,
	}, nil
}

type User struct {
	Email   string
	Address string
}
