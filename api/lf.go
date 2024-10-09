// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/oapi-codegen/runtime"
)

const (
	Api_keyScopes = "api_key.Scopes"
)

// AccountStateDelta Application state delta.
type AccountStateDelta struct {
	Address string `json:"address"`

	// Delta Application state delta.
	Delta StateDelta `json:"delta"`
}

// ApplicationStateOperation An operation against an application's global/local/box state.
type ApplicationStateOperation struct {
	// Account For local state changes, the address of the account associated with the local state.
	Account *string `json:"account,omitempty"`

	// AppStateType Type of application state. Value `g` is **global state**, `l` is **local state**, `b` is **boxes**.
	AppStateType string `json:"app-state-type"`

	// Key The key (name) of the global/local/box state.
	Key []byte `json:"key"`

	// NewValue Represents an AVM value.
	NewValue *AvmValue `json:"new-value,omitempty"`

	// Operation Operation type. Value `w` is **write**, `d` is **delete**.
	Operation string `json:"operation"`
}

// AvmValue Represents an AVM value.
type AvmValue struct {
	// Bytes bytes value.
	Bytes *[]byte `json:"bytes,omitempty"`

	// Type value type. Value `1` refers to **bytes**, value `2` refers to **uint64**
	Type int `json:"type"`

	// Uint uint value.
	Uint *int `json:"uint,omitempty"`
}

// ErrorResponse An error response with optional data field.
type ErrorResponse struct {
	Data    *map[string]interface{} `json:"data,omitempty"`
	Message string                  `json:"message"`
}

// EvalDelta Represents a TEAL value delta.
type EvalDelta struct {
	// Action \[at\] delta action.
	Action int `json:"action"`

	// Bytes \[bs\] bytes value.
	Bytes *string `json:"bytes,omitempty"`

	// Uint \[ui\] uint value.
	Uint *int `json:"uint,omitempty"`
}

// EvalDeltaKeyValue Key-value pairs for StateDelta.
type EvalDeltaKeyValue struct {
	Key string `json:"key"`

	// Value Represents a TEAL value delta.
	Value EvalDelta `json:"value"`
}

// PendingTransactionResponse Details about a pending transaction. If the transaction was recently confirmed, includes confirmation details like the round and reward details.
type PendingTransactionResponse struct {
	// ApplicationIndex The application index if the transaction was found and it created an application.
	ApplicationIndex *int `json:"application-index,omitempty"`

	// AssetClosingAmount The number of the asset's unit that were transferred to the close-to address.
	AssetClosingAmount *int `json:"asset-closing-amount,omitempty"`

	// AssetIndex The asset index if the transaction was found and it created an asset.
	AssetIndex *int `json:"asset-index,omitempty"`

	// CloseRewards Rewards in microalgos applied to the close remainder to account.
	CloseRewards *int `json:"close-rewards,omitempty"`

	// ClosingAmount Closing amount for the transaction.
	ClosingAmount *int `json:"closing-amount,omitempty"`

	// ConfirmedRound The round where this transaction was confirmed, if present.
	ConfirmedRound *int `json:"confirmed-round,omitempty"`

	// GlobalStateDelta Application state delta.
	GlobalStateDelta *StateDelta `json:"global-state-delta,omitempty"`

	// InnerTxns Inner transactions produced by application execution.
	InnerTxns *[]PendingTransactionResponse `json:"inner-txns,omitempty"`

	// LocalStateDelta Local state key/value changes for the application being executed by this transaction.
	LocalStateDelta *[]AccountStateDelta `json:"local-state-delta,omitempty"`

	// Logs Logs for the application being executed by this transaction.
	Logs *[][]byte `json:"logs,omitempty"`

	// PoolError Indicates that the transaction was kicked out of this node's transaction pool (and specifies why that happened).  An empty string indicates the transaction wasn't kicked out of this node's txpool due to an error.
	PoolError string `json:"pool-error"`

	// ReceiverRewards Rewards in microalgos applied to the receiver account.
	ReceiverRewards *int `json:"receiver-rewards,omitempty"`

	// SenderRewards Rewards in microalgos applied to the sender account.
	SenderRewards *int `json:"sender-rewards,omitempty"`

	// Txn The raw signed transaction.
	Txn map[string]interface{} `json:"txn"`
}

// ScratchChange A write operation into a scratch slot.
type ScratchChange struct {
	// NewValue Represents an AVM value.
	NewValue AvmValue `json:"new-value"`

	// Slot The scratch slot written.
	Slot int `json:"slot"`
}

// SimulationOpcodeTraceUnit The set of trace information and effect from evaluating a single opcode.
type SimulationOpcodeTraceUnit struct {
	// Pc The program counter of the current opcode being evaluated.
	Pc int `json:"pc"`

	// ScratchChanges The writes into scratch slots.
	ScratchChanges *[]ScratchChange `json:"scratch-changes,omitempty"`

	// SpawnedInners The indexes of the traces for inner transactions spawned by this opcode, if any.
	SpawnedInners *[]int `json:"spawned-inners,omitempty"`

	// StackAdditions The values added by this opcode to the stack.
	StackAdditions *[]AvmValue `json:"stack-additions,omitempty"`

	// StackPopCount The number of deleted stack values by this opcode.
	StackPopCount *int `json:"stack-pop-count,omitempty"`

	// StateChanges The operations against the current application's states.
	StateChanges *[]ApplicationStateOperation `json:"state-changes,omitempty"`
}

// SimulationTransactionExecTrace The execution trace of calling an app or a logic sig, containing the inner app call trace in a recursive way.
type SimulationTransactionExecTrace struct {
	// ApprovalProgramHash SHA512_256 hash digest of the approval program executed in transaction.
	ApprovalProgramHash *[]byte `json:"approval-program-hash,omitempty"`

	// ApprovalProgramTrace Program trace that contains a trace of opcode effects in an approval program.
	ApprovalProgramTrace *[]SimulationOpcodeTraceUnit `json:"approval-program-trace,omitempty"`

	// ClearStateProgramHash SHA512_256 hash digest of the clear state program executed in transaction.
	ClearStateProgramHash *[]byte `json:"clear-state-program-hash,omitempty"`

	// ClearStateProgramTrace Program trace that contains a trace of opcode effects in a clear state program.
	ClearStateProgramTrace *[]SimulationOpcodeTraceUnit `json:"clear-state-program-trace,omitempty"`

	// ClearStateRollback If true, indicates that the clear state program failed and any persistent state changes it produced should be reverted once the program exits.
	ClearStateRollback *bool `json:"clear-state-rollback,omitempty"`

	// ClearStateRollbackError The error message explaining why the clear state program failed. This field will only be populated if clear-state-rollback is true and the failure was due to an execution error.
	ClearStateRollbackError *string `json:"clear-state-rollback-error,omitempty"`

	// InnerTrace An array of SimulationTransactionExecTrace representing the execution trace of any inner transactions executed.
	InnerTrace *[]SimulationTransactionExecTrace `json:"inner-trace,omitempty"`

	// LogicSigHash SHA512_256 hash digest of the logic sig executed in transaction.
	LogicSigHash *[]byte `json:"logic-sig-hash,omitempty"`

	// LogicSigTrace Program trace that contains a trace of opcode effects in a logic sig.
	LogicSigTrace *[]SimulationOpcodeTraceUnit `json:"logic-sig-trace,omitempty"`
}

// StateDelta Application state delta.
type StateDelta = []EvalDeltaKeyValue

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn func(ctx context.Context, req *http.Request) error

// Doer performs HTTP requests.
//
// The standard http.Client implements this interface.
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Algod which conforms to the OpenAPI3 specification for this service.
type Algod struct {
	// The endpoint of the server conforming to this interface, with scheme,
	// https://api.deepmap.com for example. This can contain a path relative
	// to the server, such as https://api.deepmap.com/dev-test, and all the
	// paths in the swagger spec will be appended to the server.
	Server string

	// Doer for performing requests, typically a *http.Client with any
	// customized settings, such as certificate chains.
	Client HttpRequestDoer

	// A list of callbacks for modifying requests which are generated before sending over
	// the network.
	RequestEditors []RequestEditorFn
}

// ClientOption allows setting custom parameters during construction
type ClientOption func(*Algod) error

// Creates a new Algod, with reasonable defaults
func NewClient(server string, opts ...ClientOption) (*Algod, error) {
	// create a client with sane default values
	client := Algod{
		Server: server,
	}
	// mutate client and add all optional params
	for _, o := range opts {
		if err := o(&client); err != nil {
			return nil, err
		}
	}
	// ensure the server URL always has a trailing slash
	if !strings.HasSuffix(client.Server, "/") {
		client.Server += "/"
	}
	// create httpClient, if not already present
	if client.Client == nil {
		client.Client = &http.Client{}
	}
	return &client, nil
}

// WithHTTPClient allows overriding the default Doer, which is
// automatically created using http.Client. This is useful for tests.
func WithHTTPClient(doer HttpRequestDoer) ClientOption {
	return func(c *Algod) error {
		c.Client = doer
		return nil
	}
}

// WithRequestEditorFn allows setting up a callback function, which will be
// called right before sending the request. This can be used to mutate the request.
func WithRequestEditorFn(fn RequestEditorFn) ClientOption {
	return func(c *Algod) error {
		c.RequestEditors = append(c.RequestEditors, fn)
		return nil
	}
}

// The interface specification for the client above.
type ClientInterface interface {
	// Metrics request
	Metrics(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// GetStatus request
	GetStatus(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error)

	// WaitForBlock request
	WaitForBlock(ctx context.Context, round int, reqEditors ...RequestEditorFn) (*http.Response, error)
}

func (c *Algod) Metrics(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewMetricsRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Algod) GetStatus(ctx context.Context, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewGetStatusRequest(c.Server)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

func (c *Algod) WaitForBlock(ctx context.Context, round int, reqEditors ...RequestEditorFn) (*http.Response, error) {
	req, err := NewWaitForBlockRequest(c.Server, round)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	if err := c.applyEditors(ctx, req, reqEditors); err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// NewMetricsRequest generates requests for Metrics
func NewMetricsRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/metrics")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewGetStatusRequest generates requests for GetStatus
func NewGetStatusRequest(server string) (*http.Request, error) {
	var err error

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v2/status")
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// NewWaitForBlockRequest generates requests for WaitForBlock
func NewWaitForBlockRequest(server string, round int) (*http.Request, error) {
	var err error

	var pathParam0 string

	pathParam0, err = runtime.StyleParamWithLocation("simple", false, "round", runtime.ParamLocationPath, round)
	if err != nil {
		return nil, err
	}

	serverURL, err := url.Parse(server)
	if err != nil {
		return nil, err
	}

	operationPath := fmt.Sprintf("/v2/status/wait-for-block-after/%s", pathParam0)
	if operationPath[0] == '/' {
		operationPath = "." + operationPath
	}

	queryURL, err := serverURL.Parse(operationPath)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", queryURL.String(), nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Algod) applyEditors(ctx context.Context, req *http.Request, additionalEditors []RequestEditorFn) error {
	for _, r := range c.RequestEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	for _, r := range additionalEditors {
		if err := r(ctx, req); err != nil {
			return err
		}
	}
	return nil
}

// ClientWithResponses builds on ClientInterface to offer response payloads
type ClientWithResponses struct {
	ClientInterface
}

// NewClientWithResponses creates a new ClientWithResponses, which wraps
// Client with return type handling
func NewClientWithResponses(server string, opts ...ClientOption) (*ClientWithResponses, error) {
	client, err := NewClient(server, opts...)
	if err != nil {
		return nil, err
	}
	return &ClientWithResponses{client}, nil
}

// WithBaseURL overrides the baseURL.
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Algod) error {
		newBaseURL, err := url.Parse(baseURL)
		if err != nil {
			return err
		}
		c.Server = newBaseURL.String()
		return nil
	}
}

// ClientWithResponsesInterface is the interface specification for the client with responses above.
type ClientWithResponsesInterface interface {
	// MetricsWithResponse request
	MetricsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*MetricsResponse, error)

	// GetStatusWithResponse request
	GetStatusWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetStatusResponse, error)

	// WaitForBlockWithResponse request
	WaitForBlockWithResponse(ctx context.Context, round int, reqEditors ...RequestEditorFn) (*WaitForBlockResponse, error)
}

type MetricsResponse struct {
	Body         []byte
	HTTPResponse *http.Response
}

// Status returns HTTPResponse.Status
func (r MetricsResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r MetricsResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type GetStatusResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		// Catchpoint The current catchpoint that is being caught up to
		Catchpoint *string `json:"catchpoint,omitempty"`

		// CatchpointAcquiredBlocks The number of blocks that have already been obtained by the node as part of the catchup
		CatchpointAcquiredBlocks *int `json:"catchpoint-acquired-blocks,omitempty"`

		// CatchpointProcessedAccounts The number of accounts from the current catchpoint that have been processed so far as part of the catchup
		CatchpointProcessedAccounts *int `json:"catchpoint-processed-accounts,omitempty"`

		// CatchpointProcessedKvs The number of key-values (KVs) from the current catchpoint that have been processed so far as part of the catchup
		CatchpointProcessedKvs *int `json:"catchpoint-processed-kvs,omitempty"`

		// CatchpointTotalAccounts The total number of accounts included in the current catchpoint
		CatchpointTotalAccounts *int `json:"catchpoint-total-accounts,omitempty"`

		// CatchpointTotalBlocks The total number of blocks that are required to complete the current catchpoint catchup
		CatchpointTotalBlocks *int `json:"catchpoint-total-blocks,omitempty"`

		// CatchpointTotalKvs The total number of key-values (KVs) included in the current catchpoint
		CatchpointTotalKvs *int `json:"catchpoint-total-kvs,omitempty"`

		// CatchpointVerifiedAccounts The number of accounts from the current catchpoint that have been verified so far as part of the catchup
		CatchpointVerifiedAccounts *int `json:"catchpoint-verified-accounts,omitempty"`

		// CatchpointVerifiedKvs The number of key-values (KVs) from the current catchpoint that have been verified so far as part of the catchup
		CatchpointVerifiedKvs *int `json:"catchpoint-verified-kvs,omitempty"`

		// CatchupTime CatchupTime in nanoseconds
		CatchupTime int `json:"catchup-time"`

		// LastCatchpoint The last catchpoint seen by the node
		LastCatchpoint *string `json:"last-catchpoint,omitempty"`

		// LastRound LastRound indicates the last round seen
		LastRound int `json:"last-round"`

		// LastVersion LastVersion indicates the last consensus version supported
		LastVersion string `json:"last-version"`

		// NextVersion NextVersion of consensus protocol to use
		NextVersion string `json:"next-version"`

		// NextVersionRound NextVersionRound is the round at which the next consensus version will apply
		NextVersionRound int `json:"next-version-round"`

		// NextVersionSupported NextVersionSupported indicates whether the next consensus version is supported by this node
		NextVersionSupported bool `json:"next-version-supported"`

		// StoppedAtUnsupportedRound StoppedAtUnsupportedRound indicates that the node does not support the new rounds and has stopped making progress
		StoppedAtUnsupportedRound bool `json:"stopped-at-unsupported-round"`

		// TimeSinceLastRound TimeSinceLastRound in nanoseconds
		TimeSinceLastRound int `json:"time-since-last-round"`

		// UpgradeDelay Upgrade delay
		UpgradeDelay *int `json:"upgrade-delay,omitempty"`

		// UpgradeNextProtocolVoteBefore Next protocol round
		UpgradeNextProtocolVoteBefore *int `json:"upgrade-next-protocol-vote-before,omitempty"`

		// UpgradeNoVotes No votes cast for consensus upgrade
		UpgradeNoVotes *int `json:"upgrade-no-votes,omitempty"`

		// UpgradeNodeVote This node's upgrade vote
		UpgradeNodeVote *bool `json:"upgrade-node-vote,omitempty"`

		// UpgradeVoteRounds Total voting rounds for current upgrade
		UpgradeVoteRounds *int `json:"upgrade-vote-rounds,omitempty"`

		// UpgradeVotes Total votes cast for consensus upgrade
		UpgradeVotes *int `json:"upgrade-votes,omitempty"`

		// UpgradeVotesRequired Yes votes required for consensus upgrade
		UpgradeVotesRequired *int `json:"upgrade-votes-required,omitempty"`

		// UpgradeYesVotes Yes votes cast for consensus upgrade
		UpgradeYesVotes *int `json:"upgrade-yes-votes,omitempty"`
	}
	JSON401 *ErrorResponse
	JSON500 *string
}

// Status returns HTTPResponse.Status
func (r GetStatusResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r GetStatusResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

type WaitForBlockResponse struct {
	Body         []byte
	HTTPResponse *http.Response
	JSON200      *struct {
		// Catchpoint The current catchpoint that is being caught up to
		Catchpoint *string `json:"catchpoint,omitempty"`

		// CatchpointAcquiredBlocks The number of blocks that have already been obtained by the node as part of the catchup
		CatchpointAcquiredBlocks *int `json:"catchpoint-acquired-blocks,omitempty"`

		// CatchpointProcessedAccounts The number of accounts from the current catchpoint that have been processed so far as part of the catchup
		CatchpointProcessedAccounts *int `json:"catchpoint-processed-accounts,omitempty"`

		// CatchpointProcessedKvs The number of key-values (KVs) from the current catchpoint that have been processed so far as part of the catchup
		CatchpointProcessedKvs *int `json:"catchpoint-processed-kvs,omitempty"`

		// CatchpointTotalAccounts The total number of accounts included in the current catchpoint
		CatchpointTotalAccounts *int `json:"catchpoint-total-accounts,omitempty"`

		// CatchpointTotalBlocks The total number of blocks that are required to complete the current catchpoint catchup
		CatchpointTotalBlocks *int `json:"catchpoint-total-blocks,omitempty"`

		// CatchpointTotalKvs The total number of key-values (KVs) included in the current catchpoint
		CatchpointTotalKvs *int `json:"catchpoint-total-kvs,omitempty"`

		// CatchpointVerifiedAccounts The number of accounts from the current catchpoint that have been verified so far as part of the catchup
		CatchpointVerifiedAccounts *int `json:"catchpoint-verified-accounts,omitempty"`

		// CatchpointVerifiedKvs The number of key-values (KVs) from the current catchpoint that have been verified so far as part of the catchup
		CatchpointVerifiedKvs *int `json:"catchpoint-verified-kvs,omitempty"`

		// CatchupTime CatchupTime in nanoseconds
		CatchupTime int `json:"catchup-time"`

		// LastCatchpoint The last catchpoint seen by the node
		LastCatchpoint *string `json:"last-catchpoint,omitempty"`

		// LastRound LastRound indicates the last round seen
		LastRound int `json:"last-round"`

		// LastVersion LastVersion indicates the last consensus version supported
		LastVersion string `json:"last-version"`

		// NextVersion NextVersion of consensus protocol to use
		NextVersion string `json:"next-version"`

		// NextVersionRound NextVersionRound is the round at which the next consensus version will apply
		NextVersionRound int `json:"next-version-round"`

		// NextVersionSupported NextVersionSupported indicates whether the next consensus version is supported by this node
		NextVersionSupported bool `json:"next-version-supported"`

		// StoppedAtUnsupportedRound StoppedAtUnsupportedRound indicates that the node does not support the new rounds and has stopped making progress
		StoppedAtUnsupportedRound bool `json:"stopped-at-unsupported-round"`

		// TimeSinceLastRound TimeSinceLastRound in nanoseconds
		TimeSinceLastRound int `json:"time-since-last-round"`

		// UpgradeDelay Upgrade delay
		UpgradeDelay *int `json:"upgrade-delay,omitempty"`

		// UpgradeNextProtocolVoteBefore Next protocol round
		UpgradeNextProtocolVoteBefore *int `json:"upgrade-next-protocol-vote-before,omitempty"`

		// UpgradeNoVotes No votes cast for consensus upgrade
		UpgradeNoVotes *int `json:"upgrade-no-votes,omitempty"`

		// UpgradeNodeVote This node's upgrade vote
		UpgradeNodeVote *bool `json:"upgrade-node-vote,omitempty"`

		// UpgradeVoteRounds Total voting rounds for current upgrade
		UpgradeVoteRounds *int `json:"upgrade-vote-rounds,omitempty"`

		// UpgradeVotes Total votes cast for consensus upgrade
		UpgradeVotes *int `json:"upgrade-votes,omitempty"`

		// UpgradeVotesRequired Yes votes required for consensus upgrade
		UpgradeVotesRequired *int `json:"upgrade-votes-required,omitempty"`

		// UpgradeYesVotes Yes votes cast for consensus upgrade
		UpgradeYesVotes *int `json:"upgrade-yes-votes,omitempty"`
	}
	JSON400 *ErrorResponse
	JSON401 *ErrorResponse
	JSON500 *ErrorResponse
	JSON503 *ErrorResponse
}

// Status returns HTTPResponse.Status
func (r WaitForBlockResponse) Status() string {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.Status
	}
	return http.StatusText(0)
}

// StatusCode returns HTTPResponse.StatusCode
func (r WaitForBlockResponse) StatusCode() int {
	if r.HTTPResponse != nil {
		return r.HTTPResponse.StatusCode
	}
	return 0
}

// MetricsWithResponse request returning *MetricsResponse
func (c *ClientWithResponses) MetricsWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*MetricsResponse, error) {
	rsp, err := c.Metrics(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseMetricsResponse(rsp)
}

// GetStatusWithResponse request returning *GetStatusResponse
func (c *ClientWithResponses) GetStatusWithResponse(ctx context.Context, reqEditors ...RequestEditorFn) (*GetStatusResponse, error) {
	rsp, err := c.GetStatus(ctx, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseGetStatusResponse(rsp)
}

// WaitForBlockWithResponse request returning *WaitForBlockResponse
func (c *ClientWithResponses) WaitForBlockWithResponse(ctx context.Context, round int, reqEditors ...RequestEditorFn) (*WaitForBlockResponse, error) {
	rsp, err := c.WaitForBlock(ctx, round, reqEditors...)
	if err != nil {
		return nil, err
	}
	return ParseWaitForBlockResponse(rsp)
}

// ParseMetricsResponse parses an HTTP response from a MetricsWithResponse call
func ParseMetricsResponse(rsp *http.Response) (*MetricsResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &MetricsResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	return response, nil
}

// ParseGetStatusResponse parses an HTTP response from a GetStatusWithResponse call
func ParseGetStatusResponse(rsp *http.Response) (*GetStatusResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &GetStatusResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			// Catchpoint The current catchpoint that is being caught up to
			Catchpoint *string `json:"catchpoint,omitempty"`

			// CatchpointAcquiredBlocks The number of blocks that have already been obtained by the node as part of the catchup
			CatchpointAcquiredBlocks *int `json:"catchpoint-acquired-blocks,omitempty"`

			// CatchpointProcessedAccounts The number of accounts from the current catchpoint that have been processed so far as part of the catchup
			CatchpointProcessedAccounts *int `json:"catchpoint-processed-accounts,omitempty"`

			// CatchpointProcessedKvs The number of key-values (KVs) from the current catchpoint that have been processed so far as part of the catchup
			CatchpointProcessedKvs *int `json:"catchpoint-processed-kvs,omitempty"`

			// CatchpointTotalAccounts The total number of accounts included in the current catchpoint
			CatchpointTotalAccounts *int `json:"catchpoint-total-accounts,omitempty"`

			// CatchpointTotalBlocks The total number of blocks that are required to complete the current catchpoint catchup
			CatchpointTotalBlocks *int `json:"catchpoint-total-blocks,omitempty"`

			// CatchpointTotalKvs The total number of key-values (KVs) included in the current catchpoint
			CatchpointTotalKvs *int `json:"catchpoint-total-kvs,omitempty"`

			// CatchpointVerifiedAccounts The number of accounts from the current catchpoint that have been verified so far as part of the catchup
			CatchpointVerifiedAccounts *int `json:"catchpoint-verified-accounts,omitempty"`

			// CatchpointVerifiedKvs The number of key-values (KVs) from the current catchpoint that have been verified so far as part of the catchup
			CatchpointVerifiedKvs *int `json:"catchpoint-verified-kvs,omitempty"`

			// CatchupTime CatchupTime in nanoseconds
			CatchupTime int `json:"catchup-time"`

			// LastCatchpoint The last catchpoint seen by the node
			LastCatchpoint *string `json:"last-catchpoint,omitempty"`

			// LastRound LastRound indicates the last round seen
			LastRound int `json:"last-round"`

			// LastVersion LastVersion indicates the last consensus version supported
			LastVersion string `json:"last-version"`

			// NextVersion NextVersion of consensus protocol to use
			NextVersion string `json:"next-version"`

			// NextVersionRound NextVersionRound is the round at which the next consensus version will apply
			NextVersionRound int `json:"next-version-round"`

			// NextVersionSupported NextVersionSupported indicates whether the next consensus version is supported by this node
			NextVersionSupported bool `json:"next-version-supported"`

			// StoppedAtUnsupportedRound StoppedAtUnsupportedRound indicates that the node does not support the new rounds and has stopped making progress
			StoppedAtUnsupportedRound bool `json:"stopped-at-unsupported-round"`

			// TimeSinceLastRound TimeSinceLastRound in nanoseconds
			TimeSinceLastRound int `json:"time-since-last-round"`

			// UpgradeDelay Upgrade delay
			UpgradeDelay *int `json:"upgrade-delay,omitempty"`

			// UpgradeNextProtocolVoteBefore Next protocol round
			UpgradeNextProtocolVoteBefore *int `json:"upgrade-next-protocol-vote-before,omitempty"`

			// UpgradeNoVotes No votes cast for consensus upgrade
			UpgradeNoVotes *int `json:"upgrade-no-votes,omitempty"`

			// UpgradeNodeVote This node's upgrade vote
			UpgradeNodeVote *bool `json:"upgrade-node-vote,omitempty"`

			// UpgradeVoteRounds Total voting rounds for current upgrade
			UpgradeVoteRounds *int `json:"upgrade-vote-rounds,omitempty"`

			// UpgradeVotes Total votes cast for consensus upgrade
			UpgradeVotes *int `json:"upgrade-votes,omitempty"`

			// UpgradeVotesRequired Yes votes required for consensus upgrade
			UpgradeVotesRequired *int `json:"upgrade-votes-required,omitempty"`

			// UpgradeYesVotes Yes votes cast for consensus upgrade
			UpgradeYesVotes *int `json:"upgrade-yes-votes,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest string
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	}

	return response, nil
}

// ParseWaitForBlockResponse parses an HTTP response from a WaitForBlockWithResponse call
func ParseWaitForBlockResponse(rsp *http.Response) (*WaitForBlockResponse, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	response := &WaitForBlockResponse{
		Body:         bodyBytes,
		HTTPResponse: rsp,
	}

	switch {
	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 200:
		var dest struct {
			// Catchpoint The current catchpoint that is being caught up to
			Catchpoint *string `json:"catchpoint,omitempty"`

			// CatchpointAcquiredBlocks The number of blocks that have already been obtained by the node as part of the catchup
			CatchpointAcquiredBlocks *int `json:"catchpoint-acquired-blocks,omitempty"`

			// CatchpointProcessedAccounts The number of accounts from the current catchpoint that have been processed so far as part of the catchup
			CatchpointProcessedAccounts *int `json:"catchpoint-processed-accounts,omitempty"`

			// CatchpointProcessedKvs The number of key-values (KVs) from the current catchpoint that have been processed so far as part of the catchup
			CatchpointProcessedKvs *int `json:"catchpoint-processed-kvs,omitempty"`

			// CatchpointTotalAccounts The total number of accounts included in the current catchpoint
			CatchpointTotalAccounts *int `json:"catchpoint-total-accounts,omitempty"`

			// CatchpointTotalBlocks The total number of blocks that are required to complete the current catchpoint catchup
			CatchpointTotalBlocks *int `json:"catchpoint-total-blocks,omitempty"`

			// CatchpointTotalKvs The total number of key-values (KVs) included in the current catchpoint
			CatchpointTotalKvs *int `json:"catchpoint-total-kvs,omitempty"`

			// CatchpointVerifiedAccounts The number of accounts from the current catchpoint that have been verified so far as part of the catchup
			CatchpointVerifiedAccounts *int `json:"catchpoint-verified-accounts,omitempty"`

			// CatchpointVerifiedKvs The number of key-values (KVs) from the current catchpoint that have been verified so far as part of the catchup
			CatchpointVerifiedKvs *int `json:"catchpoint-verified-kvs,omitempty"`

			// CatchupTime CatchupTime in nanoseconds
			CatchupTime int `json:"catchup-time"`

			// LastCatchpoint The last catchpoint seen by the node
			LastCatchpoint *string `json:"last-catchpoint,omitempty"`

			// LastRound LastRound indicates the last round seen
			LastRound int `json:"last-round"`

			// LastVersion LastVersion indicates the last consensus version supported
			LastVersion string `json:"last-version"`

			// NextVersion NextVersion of consensus protocol to use
			NextVersion string `json:"next-version"`

			// NextVersionRound NextVersionRound is the round at which the next consensus version will apply
			NextVersionRound int `json:"next-version-round"`

			// NextVersionSupported NextVersionSupported indicates whether the next consensus version is supported by this node
			NextVersionSupported bool `json:"next-version-supported"`

			// StoppedAtUnsupportedRound StoppedAtUnsupportedRound indicates that the node does not support the new rounds and has stopped making progress
			StoppedAtUnsupportedRound bool `json:"stopped-at-unsupported-round"`

			// TimeSinceLastRound TimeSinceLastRound in nanoseconds
			TimeSinceLastRound int `json:"time-since-last-round"`

			// UpgradeDelay Upgrade delay
			UpgradeDelay *int `json:"upgrade-delay,omitempty"`

			// UpgradeNextProtocolVoteBefore Next protocol round
			UpgradeNextProtocolVoteBefore *int `json:"upgrade-next-protocol-vote-before,omitempty"`

			// UpgradeNoVotes No votes cast for consensus upgrade
			UpgradeNoVotes *int `json:"upgrade-no-votes,omitempty"`

			// UpgradeNodeVote This node's upgrade vote
			UpgradeNodeVote *bool `json:"upgrade-node-vote,omitempty"`

			// UpgradeVoteRounds Total voting rounds for current upgrade
			UpgradeVoteRounds *int `json:"upgrade-vote-rounds,omitempty"`

			// UpgradeVotes Total votes cast for consensus upgrade
			UpgradeVotes *int `json:"upgrade-votes,omitempty"`

			// UpgradeVotesRequired Yes votes required for consensus upgrade
			UpgradeVotesRequired *int `json:"upgrade-votes-required,omitempty"`

			// UpgradeYesVotes Yes votes cast for consensus upgrade
			UpgradeYesVotes *int `json:"upgrade-yes-votes,omitempty"`
		}
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON200 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 400:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON400 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 401:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON401 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 500:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON500 = &dest

	case strings.Contains(rsp.Header.Get("Content-Type"), "json") && rsp.StatusCode == 503:
		var dest ErrorResponse
		if err := json.Unmarshal(bodyBytes, &dest); err != nil {
			return nil, err
		}
		response.JSON503 = &dest

	}

	return response, nil
}
