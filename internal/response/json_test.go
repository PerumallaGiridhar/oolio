package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSON_SetsStatusAndContentType(t *testing.T) {
	rr := httptest.NewRecorder()

	body := map[string]string{"message": "ok"}
	JSON(rr, http.StatusCreated, body)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", ct)
	}

	var decoded map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	if decoded["message"] != "ok" {
		t.Fatalf("expected message 'ok', got %#v", decoded)
	}
}

func TestJSONErrorResponse_ShapesBody(t *testing.T) {
	rr := httptest.NewRecorder()

	JSONErrorResponse(rr, http.StatusBadRequest, "something went wrong")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	var decoded map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	if decoded["error"] != "something went wrong" {
		t.Fatalf("expected error message %q, got %#v", "something went wrong", decoded["error"])
	}
}

func TestJSONValidationErrorResponse_ShapesBody(t *testing.T) {
	rr := httptest.NewRecorder()

	fields := map[string]string{
		"email": "invalid email",
		"name":  "required",
	}

	JSONValidationErrorResponse(rr, fields)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected status %d, got %d", http.StatusUnprocessableEntity, rr.Code)
	}

	var decoded map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}

	if decoded["error"] != "validation_failed" {
		t.Fatalf("expected error 'validation_failed', got %#v", decoded["error"])
	}

	fieldsVal, ok := decoded["fields"].(map[string]any)
	if !ok {
		t.Fatalf("expected 'fields' to be an object, got %#v", decoded["fields"])
	}

	if fieldsVal["email"] != "invalid email" || fieldsVal["name"] != "required" {
		t.Fatalf("unexpected fields payload: %#v", fieldsVal)
	}
}

func TestJSONResponse_DelegatesToJSON(t *testing.T) {
	rr := httptest.NewRecorder()

	data := map[string]string{"k": "v"}
	JSONResponse(rr, http.StatusOK, data)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var decoded map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("response body is not valid JSON: %v", err)
	}
	if decoded["k"] != "v" {
		t.Fatalf("expected body %v, got %v", data, decoded)
	}
}
