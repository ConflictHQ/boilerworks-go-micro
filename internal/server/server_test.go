package server_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ConflictHQ/boilerworks-go-micro/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-micro/internal/server"
)

type apiResponse struct {
	Ok      bool            `json:"ok"`
	Message string          `json:"message,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
	Errors  []string        `json:"errors,omitempty"`
}

var (
	testRouter http.Handler
	testQ      *queries.Queries
	testPool   *pgxpool.Pool
	adminKey   = "test-admin-key-with-all-scopes"
	readKey    = "test-read-only-key"
	noScopeKey = "test-no-scope-key"
)

func TestMain(m *testing.M) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5438/boilerworks?sslmode=disable"
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to db: %v\n", err)
		os.Exit(1)
	}
	testPool = pool
	testQ = queries.New(pool)

	// Clean up test data
	_, _ = pool.Exec(ctx, "DELETE FROM events")
	_, _ = pool.Exec(ctx, "DELETE FROM api_keys")

	// Create test API keys
	createTestKey(ctx, adminKey, []string{"*"})
	createTestKey(ctx, readKey, []string{"events.read"})
	createTestKey(ctx, noScopeKey, []string{})

	testRouter = server.New(testQ)

	code := m.Run()

	// Cleanup
	_, _ = pool.Exec(ctx, "DELETE FROM events")
	_, _ = pool.Exec(ctx, "DELETE FROM api_keys")
	pool.Close()
	os.Exit(code)
}

func createTestKey(ctx context.Context, rawKey string, scopes []string) {
	hash := sha256.Sum256([]byte(rawKey))
	keyHash := hex.EncodeToString(hash[:])
	_, _ = testQ.CreateApiKey(ctx, queries.CreateApiKeyParams{
		Name:    "test-" + rawKey[:8],
		KeyHash: keyHash,
		Scopes:  scopes,
	})
}

func doRequest(method, path string, body interface{}, apiKey string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if apiKey != "" {
		req.Header.Set("X-API-Key", apiKey)
	}
	w := httptest.NewRecorder()
	testRouter.ServeHTTP(w, req)
	return w
}

func parseResponse(t *testing.T, w *httptest.ResponseRecorder) apiResponse {
	t.Helper()
	var resp apiResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v\nbody: %s", err, w.Body.String())
	}
	return resp
}

// --- Health ---

func TestHealthCheck(t *testing.T) {
	w := doRequest("GET", "/health", nil, "")
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestHealthCheckNoAuth(t *testing.T) {
	w := doRequest("GET", "/health", nil, "")
	if w.Code != http.StatusOK {
		t.Errorf("health should not require auth, got %d", w.Code)
	}
}

// --- Auth ---

func TestAuthMissingKey(t *testing.T) {
	w := doRequest("GET", "/events", nil, "")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthInvalidKey(t *testing.T) {
	w := doRequest("GET", "/events", nil, "totally-bogus-key")
	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthValidKey(t *testing.T) {
	w := doRequest("GET", "/events", nil, adminKey)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// --- Scope enforcement ---

func TestScopeReadOnlyCannotWrite(t *testing.T) {
	body := map[string]interface{}{"type": "test.event", "payload": map[string]string{"k": "v"}}
	w := doRequest("POST", "/events", body, readKey)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestScopeReadOnlyCanRead(t *testing.T) {
	w := doRequest("GET", "/events", nil, readKey)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestScopeNoScopeBlocked(t *testing.T) {
	w := doRequest("GET", "/events", nil, noScopeKey)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}

func TestScopeWildcardGrantsAll(t *testing.T) {
	w := doRequest("GET", "/api-keys", nil, adminKey)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200 for wildcard scope, got %d", w.Code)
	}
}

// --- Events CRUD ---

func TestCreateEvent(t *testing.T) {
	body := map[string]interface{}{"type": "order.created", "payload": map[string]string{"id": "123"}}
	w := doRequest("POST", "/events", body, adminKey)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseResponse(t, w)
	if !resp.Ok {
		t.Errorf("expected ok=true")
	}
}

func TestCreateEventMissingType(t *testing.T) {
	body := map[string]interface{}{"payload": map[string]string{"id": "123"}}
	w := doRequest("POST", "/events", body, adminKey)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestListEvents(t *testing.T) {
	w := doRequest("GET", "/events", nil, adminKey)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	resp := parseResponse(t, w)
	if !resp.Ok {
		t.Errorf("expected ok=true")
	}
}

func TestListEventsByType(t *testing.T) {
	w := doRequest("GET", "/events?type=order.created", nil, adminKey)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetEvent(t *testing.T) {
	// Create an event first
	body := map[string]interface{}{"type": "get.test", "payload": map[string]string{"k": "v"}}
	createW := doRequest("POST", "/events", body, adminKey)
	if createW.Code != http.StatusCreated {
		t.Fatalf("setup: expected 201, got %d", createW.Code)
	}

	var createResp struct {
		Ok   bool `json:"ok"`
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(createW.Body.Bytes(), &createResp)

	w := doRequest("GET", "/events/"+createResp.Data.ID, nil, adminKey)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetEventNotFound(t *testing.T) {
	w := doRequest("GET", "/events/00000000-0000-0000-0000-000000000000", nil, adminKey)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestSoftDeleteEvent(t *testing.T) {
	// Create an event
	body := map[string]interface{}{"type": "delete.test", "payload": map[string]string{}}
	createW := doRequest("POST", "/events", body, adminKey)
	var createResp struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	_ = json.Unmarshal(createW.Body.Bytes(), &createResp)

	// Delete it
	w := doRequest("DELETE", "/events/"+createResp.Data.ID, nil, adminKey)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	// Should be gone from GET
	getW := doRequest("GET", "/events/"+createResp.Data.ID, nil, adminKey)
	if getW.Code != http.StatusNotFound {
		t.Errorf("expected 404 after soft delete, got %d", getW.Code)
	}
}

// --- API Keys ---

func TestCreateApiKey(t *testing.T) {
	body := map[string]interface{}{"name": "test-created-key", "scopes": []string{"events.read"}}
	w := doRequest("POST", "/api-keys", body, adminKey)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	resp := parseResponse(t, w)
	if !resp.Ok {
		t.Errorf("expected ok=true")
	}
}

func TestListApiKeys(t *testing.T) {
	w := doRequest("GET", "/api-keys", nil, adminKey)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestCreateApiKeyMissingName(t *testing.T) {
	body := map[string]interface{}{"scopes": []string{"events.read"}}
	w := doRequest("POST", "/api-keys", body, adminKey)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestApiKeyManageRequiresScope(t *testing.T) {
	w := doRequest("GET", "/api-keys", nil, readKey)
	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403, got %d", w.Code)
	}
}
