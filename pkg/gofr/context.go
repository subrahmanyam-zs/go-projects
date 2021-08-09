package gofr

import (
	ctx "context"
	"net/http"
	"strings"
	"sync"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"developer.zopsmart.com/go/gofr/pkg/datastore/pubsub"
	"developer.zopsmart.com/go/gofr/pkg/gofr/request"
	"developer.zopsmart.com/go/gofr/pkg/gofr/responder"
	"developer.zopsmart.com/go/gofr/pkg/log"
	"developer.zopsmart.com/go/gofr/pkg/middleware/oauth"
)

type Context struct {
	ctx.Context
	*Gofr

	Logger log.Logger
	resp   responder.Responder
	req    request.Request

	WebSocketConnection *websocket.Conn
	ServerPush          http.Pusher
	ServerFlush         http.Flusher
}

func NewContext(w responder.Responder, r request.Request, k *Gofr) *Context {
	var cID string
	if r != nil {
		cID = r.Header("X-Correlation-Id")
	}

	return &Context{
		req:    r,
		resp:   w,
		Gofr:   k,
		Logger: log.NewCorrelationLogger(cID),
	}
}

func (c *Context) reset(w responder.Responder, r request.Request) {
	c.req = r
	c.resp = w
	c.Context = nil
	c.Logger = nil
}

// Trace returns an open telemetry span. We have to always close the span after corresponding work is done.
func (c *Context) Trace(name string) trace.Span {
	tr := otel.GetTracerProvider().Tracer("gofr-context")
	_, span := tr.Start(c.Context, name)
	return span
}

// Request returns the underlying HTTP request
func (c *Context) Request() *http.Request {
	return c.req.Request()
}

func (c *Context) Param(key string) string {
	return c.req.Param(key)
}

func (c *Context) Params() map[string]string {
	return c.req.Params()
}

func (c *Context) PathParam(key string) string {
	return c.req.PathParam(key)
}

func (c *Context) Bind(i interface{}) error {
	return c.req.Bind(i)
}

func (c *Context) BindStrict(i interface{}) error {
	return c.req.BindStrict(i)
}

func (c *Context) Header(key string) string {
	return c.req.Header(key)
}

// Log logs the key-value pair into the logs
func (c *Context) Log(key string, value interface{}) {
	// This section takes care of middleware logging
	if key == "correlationID" { // This condition will not allow the user to unset the CorrelationID.
		return
	}

	r := c.Request()
	appLogData, ok := r.Context().Value(appData).(*sync.Map)

	if !ok {
		c.Logger.Warn("couldn't log appData")
		return
	}

	appLogData.Store(key, value)
	*r = *r.WithContext(ctx.WithValue(r.Context(), appData, appLogData))

	// This section takes care of all the individual context loggers
	c.Logger.AddData(key, value)
}

// SetPathParams sets the URL path variables to the given value. These can be accessed
// by c.PathParam(key). This method should only be used for testing purposes.
func (c *Context) SetPathParams(pathParams map[string]string) {
	r := c.req.Request()

	r = mux.SetURLVars(r, pathParams)

	c.req = request.NewHTTPRequest(r)
}

func (c *Context) getMapClaims() jwt.MapClaims {
	claims, _ := c.Context.Value(oauth.JWTContextKey("claims")).(jwt.MapClaims)
	return claims
}

func (c *Context) ValidateClaimSub(subject string) bool {
	claims := c.getMapClaims()

	sub, ok := claims["sub"]
	if ok && sub == subject {
		return true
	}

	return false
}

func (c *Context) ValidateClaimsPFCX(pfcx string) bool {
	claims := c.getMapClaims()

	pfcxValue, ok := claims["pfcx"]
	if ok && pfcxValue == pfcx {
		return true
	}

	return false
}

func (c *Context) ValidateClaimsScope(scope string) bool {
	claims := c.getMapClaims()

	scopes, ok := claims["scope"]

	if !ok {
		return false
	}

	scopesArr := strings.Split(scopes.(string), " ")

	for i := range scopesArr {
		if scopesArr[i] == scope {
			return true
		}
	}

	return false
}

/*
	PublishEventWithOptions publishes message to the pubsub(kafka) configured.

		Ability to provide additional options as described in PublishOptions struct

		returns error if publish encounters a failure
*/
func (c *Context) PublishEventWithOptions(key string, value interface{}, headers map[string]string, options *pubsub.PublishOptions) error {
	return c.PubSub.PublishEventWithOptions(key, value, headers, options)
}

/*
	PublishEvent publishes message to the pubsub(kafka) configured.

		Information like topic is read from config, timestamp is set to current time
		other fields like offset and partition are set to it's default value
		if desire to overwrite these fields, refer PublishEventWithOptions() method above

		returns error if publish encounters a failure
*/
func (c *Context) PublishEvent(key string, value interface{}, headers map[string]string) error {
	return c.PubSub.PublishEvent(key, value, headers)
}

/*
	Subscribe read messages from the pubsub(kafka) configured.

		If multiple topics are provided in the environment or
		in kafka config while creating the consumer, reads messages from multiple topics
		reads only one message at a time. If desire to read multiple messages
		call Subscribe in a for loop

		returns error if subscribe encounters a failure
		on success returns the message received in the Message struct format
*/
func (c *Context) Subscribe(target interface{}) (*pubsub.Message, error) {
	message, err := c.PubSub.Subscribe()
	if err != nil {
		return message, err
	}

	return message, c.PubSub.Bind([]byte(message.Value), &target)
}

/*
		SubscribeWithCommit read messages from the pubsub(kafka) configured.

			calls the CommitFunc after subscribing message from kafka and based on
	        the return values decides whether to commit message and consume another message
*/
func (c *Context) SubscribeWithCommit(f pubsub.CommitFunc) (*pubsub.Message, error) {
	return c.PubSub.SubscribeWithCommit(f)
}
