package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/algorandfoundation/hack-tui/api"
)

// Account represents a user's account, including address, status, balance, and number of keys.
type Account struct {
	// Account Address is the algorand encoded address
	Address string
	// Status is the Online/Offline/"NotParticipating" status of the account
	Status string
	// Balance is the current holdings in ALGO for the address.
	// the balance should be tracked infrequently and use an appropriate distance from the
	// LastModified value.
	Balance int
	// A count of how many participation Keys exist on this node for this Account
	Keys int
	// Expires is the date the participation key will expire
	Expires time.Time
	// The LastModified round, this only pertains to keys that can be updated
	LastModified int
}

// GetAccount get online status of Account
func GetAccount(client *api.ClientWithResponses, address string) (api.Account, error) {
	var format api.AccountInformationParamsFormat = "json"
	r, err := client.AccountInformationWithResponse(
		context.Background(),
		address,
		&api.AccountInformationParams{
			Format: &format,
		})

	var accountInfo api.Account
	if err != nil {
		return accountInfo, err
	}

	if r.StatusCode() != 200 {
		return accountInfo, errors.New(fmt.Sprintf("Failed to get account information. Received error code: %d", r.StatusCode()))
	}

	return *r.JSON200, nil
}

// AccountsFromParticipationKeys maps an array of api.ParticipationKey to a keyed map of Account
func AccountsFromState(state *StateModel, client *api.ClientWithResponses) map[string]Account {
	values := make(map[string]Account)
	if state == nil || state.ParticipationKeys == nil {
		return values
	}
	for _, key := range *state.ParticipationKeys {
		val, ok := values[key.Address]
		if !ok {

			account, err := GetAccount(client, key.Address)

			// TODO: handle error
			if err != nil {
				// TODO: Logging
				panic(err)
			}

			var expires = time.Now()
			if key.EffectiveLastValid != nil {
				now := time.Now()
				roundDiff := max(0, *key.EffectiveLastValid-int(state.Status.LastRound))
				distance := int(state.Metrics.RoundTime) * roundDiff
				expires = now.Add(time.Duration(distance))
			}

			values[key.Address] = Account{
				Address: key.Address,
				Status:  account.Status,
				Balance: account.Amount / 1000000,
				Expires: expires,
				Keys:    1,
			}
		} else {
			val.Keys++
			values[key.Address] = val
		}
	}

	return values
}
