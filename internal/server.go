package internal

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/brotich/go-log-forwarder/config"
	"github.com/brotich/go-log-forwarder/internal/logfmt"
)

var formatter = logfmt.JSONFmt{}

func forwardLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()

	// http://localhost:8090?url=https://www.postb.in/1746619276054-7661341652274&fmt=logfmt

	defer r.Body.Close()
	pURL := query.Get("url")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	inBody := map[string]any{}
	err = json.Unmarshal(body, &inBody)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	outBody, err := formatter.ToLogfmt(inBody)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	req, err := http.NewRequestWithContext(ctx, r.Method, pURL, nil)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	orginalHost := req.Host
	req.Header = r.Header
	req.Host = orginalHost

	req.Body = io.NopCloser(bytes.NewBuffer(outBody))
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte(err.Error()))
		return
	}

	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(503)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)

	slog.InfoContext(ctx, "handle request", slog.String("url", pURL), slog.String("out", string(outBody)))

}

func NewServer(ctx context.Context, cfg config.Config) (*http.Server, error) {
	http.HandleFunc("/", forwardLog)

	srv := &http.Server{
		Addr: cfg.ListenAddr,
	}
	return srv, nil
}
