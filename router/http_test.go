package router

import "testing"

// should correctly parse the request method
func TestHTTPActionRequestMethod(t *testing.T) {
	result, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "get"
		}
	}`))

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	l := len(result.Request.Path)
	if l != 1 {
		t.Fatalf("expected a path length of 1, got %d", l)
	}

	method := result.Request.Path[0].String()
	if method != "get" {
		t.Errorf("expected \"get\", received %q", method)
	}
}

// should parse HTTP methods to a standard case
func TestHTTPActionRequestMethodCase(t *testing.T) {
	result, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "GET"
		}
	}`))

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	l := len(result.Request.Path)
	if l != 1 {
		t.Fatalf("expected a path length of 1, got %d", l)
	}

	method := result.Request.Path[0].String()
	if method != "get" {
		t.Errorf("expected \"get\", received %q", method)
	}
}

// should return an error when receiving an unrecognized method
func TestHTTPActionRequestUnknownMethod(t *testing.T) {
	_, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "unknown"
		}
	}`))

	if err == nil {
		t.Errorf("expected error, but received none")
	}
}

// should correctly parse the path
func TestHTTPActionRequestPath(t *testing.T) {
	result, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "get",
			"path": "/my/test/path"
		}
	}`))

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	if result.Request.Path == nil {
		t.Fatal("path is nil")
	}

	expected := []string{"get", "my", "test", "path"}
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
	result, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "get",
			"headers": {
				"content-type": "test"
			}
		}
	}`))

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
	result, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "get",
			"headers": {
				"[": "test"
			}
		}
	}`))

	if err == nil {
		t.Error("error should not be nil")
	}

	if result != nil {
		t.Error("expected result to be nil")
	}
}

// should correctly parse request body
func TestHTTPActionRequestBody(t *testing.T) {
	result, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "get",
			"body": {
				"data": "test body"
			}
		}
	}`))

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
	result, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "get",
			"body": {
				"data": "["
			}
		}
	}`))

	if err == nil {
		t.Error("error should not be nil")
	}

	if result != nil {
		t.Error("expected result to be nil")
	}
}

// should correctly parse a response status code
func TestHTTPActionResponseStatus(t *testing.T) {
	result, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "get"
		},
		"response": {
			"status": 200
		}
	}`))

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
	result, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "get"
		},
		"response": {
			"headers": {
				"content-type": "test",
				"content-length": "100"
			}
		}
	}`))

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
	result, err := HTTPActionFromJSON([]byte(`{
		"request": {
			"method": "get"
		},
		"response": {
			"body": {
				"data": "test data"
			}
		}
	}`))

	if err != nil {
		t.Fatalf("received error (%v)", err)
	}

	expected := "{\"data\":\"test data\"}"
	body := result.Response.Body
	if string(body) != expected {
		t.Errorf("expected body to be %q, received %q", expected, body)
	}
}

// should replace variables with values
func TestHTTPActionReplace(t *testing.T) {
	result := replace([]byte("test {{replace}}"), map[string]string{
		"replace": "aaa",
	})

	if string(result) != "test aaa" {
		t.Errorf("expected \"test aaa\", got %q", result)
	}
}

// should replace all variables, even with empty strings
func TestHTTPActionReplaceNil(t *testing.T) {
	result := replace([]byte("test {{nothing}}"), nil)

	if string(result) != "test " {
		t.Errorf("expected \"test \", got %q", result)
	}
}
