package gqltest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
)

type Transport struct {
	Handler http.HandlerFunc
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	t.Handler(rec, req)
	return rec.Result(), nil
}

type Handler struct {
	QueryName  string
	MutateName string

	Handle http.HandlerFunc
}

type Option struct {
	Handlers []Handler
}

func New(t *testing.T, opts ...func(*Option)) (api.GQLClient, error) {
	option := Option{}
	for _, opt := range opts {
		opt(&option)
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Query string `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for _, h := range option.Handlers {
			if h.QueryName != "" && !strings.HasPrefix(body.Query, "query "+h.QueryName) {
				continue
			}
			if h.MutateName != "" && !strings.HasPrefix(body.Query, "mutation "+h.MutateName) {
				continue
			}
			h.Handle(w, r)
			return
		}
		http.Error(w, "no matched handler for query: "+body.Query, http.StatusNotFound)
	}
	return gh.GQLClient(&api.ClientOptions{
		Transport: &Transport{Handler: handler},
	})
}

func WithQueryOK(queryName string, body string) func(*Option) {
	h := Handler{
		QueryName: queryName,
		Handle: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body))
		},
	}
	return func(o *Option) {
		o.Handlers = append(o.Handlers, h)
	}
}

func WithMutateOK(mutateName string, body string) func(*Option) {
	h := Handler{
		MutateName: mutateName,
		Handle: func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(body))
		},
	}
	return func(o *Option) {
		o.Handlers = append(o.Handlers, h)
	}
}
