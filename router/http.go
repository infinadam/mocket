package router

import (
	"encoding/json"
	"regexp"
)

type Header struct {
	Label *regexp.Regexp
	Value *regexp.Regexp
}

type HTTPAction struct {
	Request struct {
		Headers []Header
		Body    *regexp.Regexp
	}
	Response struct {
		Status  uint8
		Headers map[string]string
		Body    string
	}
}

type httpJSON struct {
	Request struct {
		Headers map[string]string `json:"headers"`
		Body    any               `json:"body"`
	} `json:"request"`
	Response struct {
		Status  uint8             `json:"status"`
		Headers map[string]string `json:"headers"`
		Body    any               `json:"body"`
	} `json:"response"`
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
		unmarshal, requestHeaders, requestBody,
	}
	for _, f := range parsers {
		if err := f(action, &parsed); err != nil {
			return nil, err
		}
	}

	body, _ := json.Marshal(parsed.Response.Body)
	action.Response.Body = string(body)
	action.Response.Status = parsed.Response.Status
	action.Response.Headers = parsed.Response.Headers

	return action, nil
}
