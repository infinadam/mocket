package router

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Header struct {
	Label *regexp.Regexp
	Value *regexp.Regexp
}

type HTTPAction struct {
	Request struct {
		Path    []*regexp.Regexp
		Headers []Header
		Body    *regexp.Regexp
	}
	Response struct {
		Status  int
		Headers map[string]string
		Body    []byte
	}
}

type httpJSON struct {
	Request struct {
		Method  string            `json:"method"`
		Path    string            `json:"path"`
		Headers map[string]string `json:"headers"`
		Body    any               `json:"body"`
	} `json:"request"`
	Response struct {
		Status  int               `json:"status"`
		Headers map[string]string `json:"headers"`
		Body    any               `json:"body"`
	} `json:"response"`
}

func requestPath(action *HTTPAction, parsed *httpJSON) error {
	method := strings.ToLower(parsed.Request.Method)

	switch method {
	case "delete", "get", "head", "options", "patch", "post", "put":
		re := regexp.MustCompile(method)
		action.Request.Path = append(action.Request.Path, re)
	default:
		return errors.New("unrecognized method")
	}

	for _, s := range strings.Split(parsed.Request.Path, "/") {
		if s == "" {
			continue
		} else if re, err := regexp.Compile(s); err != nil {
			return nil
		} else {
			action.Request.Path = append(action.Request.Path, re)
		}
	}

	return nil
}

func requestHeaders(action *HTTPAction, parsed *httpJSON) error {
	var headers []Header
	var err error
	var entry Header

	for label, value := range parsed.Request.Headers {
		if entry.Label, err = regexp.Compile(label); err != nil {
			return err
		}
		if entry.Value, err = regexp.Compile(value); err != nil {
			return err
		}
		headers = append(headers, entry)
	}
	action.Request.Headers = headers

	return nil
}

func requestBody(action *HTTPAction, parsed *httpJSON) error {
	var err error
	body, _ := json.Marshal(parsed.Request.Body)
	unquoted, _ := strconv.Unquote(string(body))
	action.Request.Body, err = regexp.Compile(unquoted)
	return err
}

func HTTPActionFromJSON(input []byte) (*HTTPAction, error) {
	var parsed httpJSON
	action := new(HTTPAction)

	unmarshal := func(action *HTTPAction, parsed *httpJSON) error {
		// passing this should guarantee successful marshalling later...
		return json.Unmarshal([]byte(input), &parsed)
	}

	parsers := []func(action *HTTPAction, parsed *httpJSON) error{
		unmarshal, requestPath, requestHeaders, requestBody,
	}
	for _, f := range parsers {
		if err := f(action, &parsed); err != nil {
			return nil, err
		}
	}

	body, _ := json.Marshal(parsed.Response.Body)
	action.Response.Body = body
	action.Response.Status = parsed.Response.Status
	action.Response.Headers = parsed.Response.Headers

	return action, nil
}

func replace(original []byte, vars map[string]string) []byte {
	re := regexp.MustCompile(`{{([[:alnum:]]+)}}`)
	matches := re.FindAllSubmatch(original, -1)
	for _, m := range matches {
		id := string(m[1])
		r := regexp.MustCompile(string(m[0]))
		original = r.ReplaceAll(original, []byte(vars[id]))
	}

	return original
}

func merge(a map[string]string, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}
	return a
}

// TODO: where are my tests?
func (a *HTTPAction) CompareHeaders(label string, value string) (bool, map[string]string) {
	for _, h := range a.Request.Headers {
		if matched, ls := match(h.Label, label); !matched {
			continue
		} else if matched, vs := match(h.Value, value); !matched {
			continue
		} else {
			return true, merge(ls, vs)
		}
	}

	return false, nil
}

func (a *HTTPAction) CompareBody(body string) (bool, map[string]string) {
	return match(a.Request.Body, body)
}

func (a *HTTPAction) Write(w http.ResponseWriter, vars map[string]string) {
	for k, v := range a.Response.Headers {
		w.Header().Set(k, string(replace([]byte(v), vars)))
	}
	w.WriteHeader(a.Response.Status)
	w.Write(replace(a.Response.Body, vars))
}
