package internal

import "github.com/algorandfoundation/hack-tui/api"

// Account represents a user's account, including address, status, balance, and number of keys.
type Account struct {
	// Account Address is the algorand enoded address
	Address string
	// Status is general information about the account
	Status string
	// Balance is the current holdings in ALGO for the address.
	// the balance should be tracked infrequently and use an appropriate distance from the
	// LastModified value.
	Balance int
	// A count of how many participation Keys exist on this node for this Account
	Keys int
	// The LastModified round, this only pertains to keys that can be updated
	LastModified int
}

// Map and array of api.ParticipationKey to a keyed map of Account
func AccountsFromParticipationKeys(keys *[]api.ParticipationKey) map[string]Account {
	values := make(map[string]Account)

	for _, key := range *keys {
		val, ok := values[key.Address]
		if !ok {
			values[key.Address] = Account{
				Address: key.Address,
				Status:  "NA",
				Balance: 0,
				Keys:    1,
			}
		} else {
			val.Keys++
			values[key.Address] = val
		}
	}

	return values
}
