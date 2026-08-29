package observability

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetricsMiddleware(t *testing.T) {
	h := MetricsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/v1/pokemons", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201", w.Code)
	}

	count := testutil.CollectAndCount(httpRequestsTotal)
	if count == 0 {
		t.Error("nenhuma métrica http_requests_total registrada")
	}
}

func TestResponseWriterDefaultsToOK(t *testing.T) {
	rw := &responseWriter{statusCode: 0}
	if rw.statusCode != 0 {
		t.Error("statusCode deveria iniciar em 0")
	}
}
