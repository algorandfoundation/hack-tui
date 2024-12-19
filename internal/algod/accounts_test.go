package algod

import (
	"context"
	"github.com/algorandfoundation/algorun-tui/api"
	"github.com/algorandfoundation/algorun-tui/internal/test"
	"github.com/algorandfoundation/algorun-tui/internal/test/mock"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/securityprovider"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_AccountsFromState(t *testing.T) {

	// Setup elevated client
	apiToken, err := securityprovider.NewSecurityProviderApiKey("header", "X-Algo-API-Token", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	if err != nil {
		t.Fatal(err)
	}
	client, err := api.NewClientWithResponses("http://localhost:8080", api.WithRequestEditorFn(apiToken.Intercept))

	addresses, rewardsPool, feeSink, err := test.GetAddressesFromGenesis(context.Background(), client)

	if err != nil {
		t.Fatal(err)
	}

	var mapAccounts = make(map[string]api.Account)
	var onlineAccounts = make([]api.Account, 0)
	for _, address := range addresses {
		acct, err := GetAccount(client, address)
		if err != nil {
			t.Fatal(err)
		}

		assert.True(t, acct.Status == "Online" || acct.Status == "Offline")
		mapAccounts[address] = acct
		if acct.Status == "Online" {
			onlineAccounts = append(onlineAccounts, acct)
		}
	}

	acct, err := GetAccount(client, rewardsPool)
	if err != nil {
		t.Fatal(err)
	}
	if acct.Status != "Not Participating" {
		t.Fatalf("Expected RewardsPool to be 'Not Participating', got %s", acct.Status)
	}

	acct, err = GetAccount(client, feeSink)
	if err != nil {
		t.Fatal(err)
	}
	if acct.Status != "Not Participating" {
		t.Fatalf("Expected FeeSink to be 'Not Participating', got %s", acct.Status)
	}

	_, err = GetAccount(client, "invalid_address")
	if err == nil {
		t.Fatal("Expected error for invalid address")
	}

	// Mock StateModel
	state := &StateModel{
		Metrics: Metrics{
			Enabled:   true,
			Window:    100,
			RoundTime: time.Duration(2) * time.Second,
			TPS:       20,
			RX:        1024,
			TX:        2048,
		},
		Status: Status{
			State:       "WATCHING",
			Version:     "v0.0.0-test",
			Network:     "tuinet",
			Voting:      false,
			NeedsUpdate: false,
			LastRound:   1337,
			Client:      client,
			HttpPkg:     new(api.HttpPkg),
		},
		ParticipationKeys: mock.Keys,
		Client:            client,
		HttpPkg:           new(api.HttpPkg),
	}

	// Calculate expiration
	clock := new(mock.Clock)
	state.UpdateKeys(context.Background(), clock)

}
