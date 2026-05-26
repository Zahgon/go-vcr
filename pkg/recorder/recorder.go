// Copyright (c) 2015-2024 Marin Atanasov Nikolov <dnaeon@gmail.com>
// Copyright (c) 2016 David Jack <davars@gmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
// 1. Redistributions of source code must retain the above copyright
//    notice, this list of conditions and the following disclaimer
//    in this position and unchanged.
// 2. Redistributions in binary form must reproduce the above copyright
//    notice, this list of conditions and the following disclaimer in the
//    documentation and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE AUTHOR(S) ``AS IS'' AND ANY EXPRESS OR
// IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES
// OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED.
// IN NO EVENT SHALL THE AUTHOR(S) BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT
// NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF
// THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package recorder

import (
	"errors"
	"net/http"

	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
)

type MatcherFunc = cassette.MatcherFunc

// ErrNoCassetteName is an error, which is returned when the recorder was
// created without specifying a cassette name.
var ErrNoCassetteName = errors.New("no cassette name specified")

// Mode represents the mode of operation of the recorder
type Mode int

// Recorder states
const (
	// ModeRecordOnly specifies that VCR will run in recording mode
	// only. HTTP interactions will be recorded for each interaction. If the
	// cassette file is present, it will be overwritten.
	ModeRecordOnly Mode = iota

	// ModeReplayOnly specifies that VCR will only replay interactions from
	// previously recorded cassette. If an interaction is missing from the
	// cassette it will return ErrInteractionNotFound error. If the cassette
	// file is missing it will return ErrCassetteNotFound error.
	ModeReplayOnly

	// ModeReplayWithNewEpisodes starts the recorder in replay mode, where
	// existing interactions are returned from the cassette, and missing
	// ones will be recorded and added to the cassette. This mode is useful
	// in cases where you need to update an existing cassette with new
	// interactions, but don't want to wipe out previously recorded
	// interactions. If the cassette file is missing it will create a new
	// one.
	ModeReplayWithNewEpisodes

	// ModeRecordOnce will record new HTTP interactions once only. This mode
	// is useful in cases where you need to record a set of interactions
	// once only and replay only the known interactions. Unknown/missing
	// interactions will cause the recorder to return an
	// ErrInteractionNotFound error. If the cassette file is missing, it
	// will be created.
	ModeRecordOnce

	// ModePassthrough specifies that VCR will not record any interactions
	// at all. In this mode all HTTP requests will be forwarded to the
	// endpoints using the real HTTP transport. In this mode no cassette
	// will be created.
	ModePassthrough
)

// ErrInvalidMode is returned when attempting to start the recorder with invalid
// mode
var ErrInvalidMode = errors.New("invalid recorder mode")

// HookFunc represents a function, which will be invoked in different stages of
// the playback. The hook functions allow for plugging in to the playback and
// transform an interaction, if needed. For example a hook function might redact
// or remove sensitive data from a request/response before it is added to the
// in-memory cassette, or before it is saved on disk. Another use case would be
// to transform the HTTP response before it is returned to the client during
// replay mode.
type HookFunc func(i *cassette.Interaction) error

// Hook kinds
type HookKind int

const (
	// AfterCaptureHook represents a hook, which will be invoked after
	// capturing a request/response pair.
	AfterCaptureHook HookKind = iota

	// BeforeSaveHook represents a hook, which will be invoked right before
	// the cassette is saved on disk.
	BeforeSaveHook

	// BeforeResponseReplayHook represents a hook, which will be invoked
	// before replaying a previously recorded response to the client.
	BeforeResponseReplayHook

	// OnRecorderStopHook is a hook, which will be invoked when the recorder
	// is about to be stopped. This hook is useful for performing any
	// post-actions such as cleanup or reporting.
	OnRecorderStopHook
)

// Hook represents a function hook of a given kind. Depending on the hook kind,
// the function will be invoked in different stages of the playback.
type Hook struct {
	// Handler is the function which will be invoked
	Handler HookFunc

	// Kind represents the hook kind
	Kind HookKind
}

// NewHook creates a new hook.
func NewHook(handler HookFunc, kind HookKind) *Hook { _ = "STUB: not implemented"; return nil }

// PassthroughFunc is a predicate which determines whether a specific HTTP
// request is to be forwarded to the original endpoint. It should return true
// when a request needs to be passed through, and false otherwise.
type PassthroughFunc func(req *http.Request) bool

// ErrUnsafeRequestMethod is returned when the [Recorder] was configured to
// block unsafe methods, and an attempt to use such was invoked. Safe Methods
// are defined as part of RFC 9110, section 9.2.1.
var ErrUnsafeRequestMethod = errors.New("request uses an unsafe method")

type blockUnsafeMethodsRoundTripper struct {
	RoundTripper http.RoundTripper
}

func (r *blockUnsafeMethodsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// Recorder represents a type used to record and replay client and server
// interactions.
type Recorder struct {
	// Cassette used by the recorder
	cassette *cassette.Cassette

	// cassetteName is the name of the cassette to be used by the recorder.
	cassetteName string

	// mode is the mode of the recorder
	mode Mode

	// RealTransport is the underlying http.RoundTripper to make
	// the real requests
	realTransport http.RoundTripper

	// blockUnsafeMethods specifies whether to block requests when making
	// HTTP requests which are not safe. The "Safe Methods" are defined as
	// part of RFC 9110, section 9.2.1, and SHOULD NOT have side effects on
	// the server.
	blockUnsafeMethods bool

	// skipRequestLatency specifies whether to simulate the latency of the
	// recorded interaction.
	skipRequestLatency bool

	// Passthrough handlers
	passthroughs []PassthroughFunc

	// hooks is a list of hooks, which are invoked in different
	// stages of the playback.
	hooks []*Hook

	// matcher is the [MatcherFunc] predicate used to match HTTP requests
	// against recorded interactions.
	matcher MatcherFunc

	// replayableInteractions specifies whether to allow interactions to be
	// replayed multiple times.
	replayableInteractions bool

	// fs specifies custom filesystem ([cassette.FS]) implementation.
	fs cassette.FS

	marshalFunc cassette.MarshalFunc
}

// Option is a function which configures the [Recorder].
type Option func(r *Recorder)

// WithMode is an [Option], which configures the [Recorder] to run in the
// specified mode.
func WithMode(mode Mode) Option { _ = "STUB: not implemented"; return *new(Option) }

// WithRealTransport is an [Option], which configures the [Recorder] to use the
// specified [http.RoundTripper] when making actual HTTP requests.
func WithRealTransport(rt http.RoundTripper) Option { _ = "STUB: not implemented"; return *new(Option) }

// WithBlockUnsafeMethods is an [Option], which configures the [Recorder] to
// block HTTP requests, which are not considered "Safe Methods", according to
// RFC 9110, section 9.2.1.
func WithBlockUnsafeMethods(val bool) Option { _ = "STUB: not implemented"; return *new(Option) }

// WithSkipRequestLatency is an [Option], which configures the [Recorder] whether
// to simulate the latency of the recorded interaction. When set to false it
// will block for the period of time taken by the original request to simulate
// the latency between the recorder and the remote endpoints.
func WithSkipRequestLatency(val bool) Option { _ = "STUB: not implemented"; return *new(Option) }

// WithPassthrough is an [Option], which configures the [Recorder] to
// passthrough requests for requests which satisfy the provided
// [PassthroughFunc] predicate.
func WithPassthrough(passfunc PassthroughFunc) Option {
	_ = "STUB: not implemented"
	return *new(Option)
}

// WithHook is an [Option], which configures the [Recorder] to invoke the
// provided hook at the specified playback stage.
func WithHook(handler HookFunc, kind HookKind) Option {
	_ = "STUB: not implemented"
	return *new(Option)
}

// WithMatchers is an [Option], which configures the [Recorder] to use the
// provided [MatcherFunc] predicate when matching HTTP requests against record
// interactions.
func WithMatcher(matcher MatcherFunc) Option { _ = "STUB: not implemented"; return *new(Option) }

// WithReplayableInteractions is an [Option], which configures the [Recorder] to
// allow replaying interactions multiple times. This is useful in situations
// when you need to hit the same endpoint multiple times and want to replay the
// interaction from the cassette each time.
func WithReplayableInteractions(val bool) Option { _ = "STUB: not implemented"; return *new(Option) }

// WithFS is an [Option], which configures the [Recorder] to use
// custom filesystem ([cassette.FS]) implementation. This allows the [Recorder] to use any
// FS-compatible backend (e.g., local disk, in-memory, or mock) for reading and writing files.
func WithFS(fs cassette.FS) Option { _ = "STUB: not implemented"; return *new(Option) }

// WithMarshalFunc is an [Option], which configures the [Recorder] to use
// custom YAML marshal func. This allows customization of the YAML encoding
// process, such as setting string literal style, etc.
func WithMarshalFunc(marshalFunc cassette.MarshalFunc) Option {
	_ = "STUB: not implemented"
	return *new(Option)
}

// New creates a new [Recorder] and configures it using the provided options.
func New(cassetteName string, opts ...Option) (*Recorder, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// Configure the cassette based on the recorder configuration

// getCassette creates a new [*cassette.Cassette], or loads an already existing
// one depending on the mode of the recorder.
func (rec *Recorder) getCassette() (*cassette.Cassette, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// Create or the cassette depending on the mode we are operating in.

// getRoundTripper returns the [http.RoundTripper] used by the recorder.
func (rec *Recorder) getRoundTripper() http.RoundTripper {
	_ = "STUB: not implemented"
	return *new(http.RoundTripper)
}

// requestHandler proxies requests to their original destination
// If serverResponse is provided, this is used for the recording instead of using RoundTrip
func (rec *Recorder) requestHandler(r *http.Request, serverResponse *http.Response) (*cassette.Interaction, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// Interaction found, return it

// Interaction not found, we have a new episode

// Any other error is an error

// We've got an existing cassette, return what we've got

// Passthrough requests always hit the original endpoint

// When running with replayable interactions look for existing
// interaction first, so we avoid hitting multiple times the
// same endpoint.

// Interaction found, return it

// Interaction not found, we have to record it

// Any other error is an error

// Anything else hits the original endpoint

// Copy the original request, so we can read the form values

// Record the request body so we can add it to the cassette

// when serverResponse is provided by middleware, it has to be read in order
// for reqBody buffer to be populated

// Perform request to it's original destination and record the interactions
// If serverResponse is provided, use it instead

// Add interaction to the cassette

// Apply after-capture hooks before we add the interaction to
// the in-memory cassette.

// Stop is used to stop the recorder and save any recorded
// interactions if running in one of the recording modes. When
// running in ModePassthrough no cassette will be saved on disk.
func (rec *Recorder) Stop() error { _ = "STUB: not implemented"; return nil }

// Nothing to do for ModeReplayOnly and ModePassthrough here

// Apply on-recorder-stop hooks

// persisteCassette persists the cassette on disk for future re-use
func (rec *Recorder) persistCassette() error {
	_ = "STUB: not implemented"
	// Apply any before-save hooks
	return nil
}

// applyHooks applies the registered hooks of the given kind with the
// specified interaction
func (rec *Recorder) applyHooks(i *cassette.Interaction, kind HookKind) error {
	_ = "STUB: not implemented"
	return nil
}

// RoundTrip implements the [http.RoundTripper] interface
func (rec *Recorder) RoundTrip(req *http.Request) (*http.Response, error) {
	_ = "STUB: not implemented"
	return nil, nil
}

// executeAndRecord is used internally by the HTTPMiddleware to allow recording a response on the server side
func (rec *Recorder) executeAndRecord(req *http.Request, serverResponse *http.Response) (*http.Response, error) {
	_ = "STUB: not implemented"
	// Passthrough mode, use real transport
	return nil, nil
}

// Apply passthrough handler functions

// Apply before-response-replay hooks

// Apply the duration defined in the interaction

// Mode returns recorder state
func (rec *Recorder) Mode() Mode {
	_ = "STUB: not implemented"

	// GetDefaultClient returns an HTTP client with a pre-configured
	// transport
	return *new(Mode)
}

func (rec *Recorder) GetDefaultClient() *http.Client { _ = "STUB: not implemented"; return nil }

// IsNewCassette returns true, if the recorder was started with a
// new/empty cassette. Returns false, if it was started using an
// existing cassette, which was loaded.
func (rec *Recorder) IsNewCassette() bool { _ = "STUB: not implemented"; return false }

// IsRecording returns true, if the recorder is recording
// interactions, returns false otherwise. Note, that in some modes
// (e.g. ModeReplayWithNewEpisodes and ModeRecordOnce) the recorder
// might be recording new interactions. For example in ModeRecordOnce,
// we are replaying interactions only if there was an existing
// cassette, and we are recording it, if the cassette is a new one.
// ModeReplayWithNewEpisodes would replay interactions, if they are
// present in the cassette, but will also record new ones, if they are
// not part of the cassette already. In these cases the recorder is
// considered to be recording for these modes.
func (rec *Recorder) IsRecording() bool { _ = "STUB: not implemented"; return false }
