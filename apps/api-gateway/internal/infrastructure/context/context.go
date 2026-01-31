package context

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/newrelic/go-agent/v3/newrelic"
)

func SetAction(ctx context.Context, action string) context.Context {
	ctx = context.WithValue(ctx, ActionKey, action)

	return ctx
}

func SetTxnName(request *http.Request, urlName string, tags ...string) {
	txnTags := strings.Join(tags, " - ")

	name := fmt.Sprintf("%s - %s (%s)", urlName, txnTags, request.Method)

	txn := newrelic.FromContext(request.Context())
	if txn != nil {
		txn.SetName(name)
	}
}

func NewGoroutineContext(ctx context.Context) context.Context {
	txn := newrelic.FromContext(ctx)
	if txn != nil {
		newTx := txn.NewGoroutine()

		return newrelic.NewContext(ctx, newTx)
	}

	return ctx
}
