package router

import (
	"encoding/json"
	"errors"
	"regexp"
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
		Verb    string            `json:"verb"`
		Path    string            `json:"url"`
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
	verb := strings.ToLower(parsed.Request.Verb)

	switch verb {
	case "delete", "get", "head", "options", "patch", "post", "put":
		re, _ := regexp.Compile(verb)
		action.Request.Path = append(action.Request.Path, re)
	default:
		return errors.New("unrecognized verb")
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
	action.Request.Body, err = regexp.Compile(string(body))
	return err
}

func HTTPActionFromJSON(input string) (*HTTPAction, error) {
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
