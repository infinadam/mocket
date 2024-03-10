package router

import "testing"

// should correctly parse the request verb
func TestHTTPActionRequestVerb(t *testing.T) {
	result, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "get"
		}
	}`)

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	if result.Request.Verb != "get" {
		t.Errorf("expected \"get\", received %q", result.Request.Verb)
	}
}

// should parse HTTP verbs to a standard case
func TestHTTPActionRequestVerbCase(t *testing.T) {
	result, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "GET"
		}
	}`)

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	if result.Request.Verb != "get" {
		t.Errorf("expected \"get\", received %q", result.Request.Verb)
	}
}

// should return an error when receiving an unrecognized verb
func TestHTTPActionRequestUnknownVerb(t *testing.T) {
	_, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "unknown"
		}
	}`)

	if err == nil {
		t.Errorf("expected error, but received none")
	}
}

// should correctly parse the URL path
func TestHTTPActionRequestPath(t *testing.T) {
	result, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "get",
			"url": "/my/test/path"
		}
	}`)

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	if result.Request.Path == nil {
		t.Fatal("path is nil")
	}

	expected := []string{"my", "test", "path"}
	if len(result.Request.Path) != len(expected) {
		t.Fatal("path is wrong length")
	}

	for i, s := range expected {
		if result.Request.Path[i].String() != s {
			t.Errorf("expected %q, received %q", s, result.Request.Path[i])
		}
	}
}

// should correctly parse request headers
func TestHTTPActionRequestHeaders(t *testing.T) {
	result, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "get",
			"headers": {
				"content-type": "test"
			}
		}
	}`)

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	headerLen := len(result.Request.Headers)
	if headerLen != 1 {
		t.Errorf("expected 1 header, got %d", headerLen)
	}

	if result.Request.Headers[0].Label.String() != "content-type" {
		t.Error("did not correctly parse header label")
	}

	if result.Request.Headers[0].Value.String() != "test" {
		t.Error("did not correctly parse header value")
	}
}

// should return an error when failing to compile header regexp
func TestHTTPActionRequestHeaderError(t *testing.T) {
	result, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "get",
			"headers": {
				"[": "test"
			}
		}
	}`)

	if err == nil {
		t.Error("error should not be nil")
	}

	if result != nil {
		t.Error("expected result to be nil")
	}
}

// should correctly parse request body
func TestHTTPActionRequestBody(t *testing.T) {
	result, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "get",
			"body": {
				"data": "test body"
			}
		}
	}`)

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	body := result.Request.Body.String()
	expected := "{\"data\":\"test body\"}"
	if body != expected {
		t.Errorf("expected %q and got %q", expected, body)
	}
}

// should return an error when failing to compile body regexp
func TestHTTPActionRequestBodyError(t *testing.T) {
	result, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "get",
			"body": {
				"data": "["
			}
		}
	}`)

	if err == nil {
		t.Error("error should not be nil")
	}

	if result != nil {
		t.Error("expected result to be nil")
	}
}

// should correctly parse a response status code
func TestHTTPActionResponseStatus(t *testing.T) {
	result, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "get"
		},
		"response": {
			"status": 200
		}
	}`)

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	status := result.Response.Status
	if status != 200 {
		t.Errorf("expected status to be 200, got %d", status)
	}
}

// should correctly parse response headers
func TestHTTPActionResponseHeaders(t *testing.T) {
	result, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "get"
		},
		"response": {
			"headers": {
				"content-type": "test",
				"content-length": "100"
			}
		}
	}`)

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	headers := result.Response.Headers
	expected := map[string]string{
		"content-type":   "test",
		"content-length": "100",
	}
	for key, value := range expected {
		if headers[key] != value {
			t.Errorf("expected %q to be %q was %q", key, value, headers[value])
		}
	}
}

// should correctly parse a response body
func TestHTTPActionResponseBody(t *testing.T) {
	result, err := HTTPActionFromJSON(`{
		"request": {
			"verb": "get"
		},
		"response": {
			"body": {
				"data": "test data"
			}
		}
	}`)

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	expected := "{\"data\":\"test data\"}"
	body := result.Response.Body
	if body != expected {
		t.Errorf("expected body to be %q, received %q", expected, body)
	}
}
