package main

import (
	"github.com/infinadam/mocket/router"
	"io"
	"net/http"
	"os"
	"strings"
)

type Server struct {
	path router.Path
}

func actionFromEntry(dir string, e os.DirEntry) (*router.HTTPAction, error) {
	if e.IsDir() {
		return nil, nil
	}

	if json, err := os.ReadFile(dir + "/" + e.Name()); err != nil {
		return nil, err
	} else if action, err := router.HTTPActionFromJSON(json); err != nil {
		return nil, err
	} else {
		return action, nil
	}
}

func MakeServer(dir string) (*Server, error) {
	var entries []os.DirEntry
	var err error
	server := new(Server)

	if entries, err = os.ReadDir(dir); err != nil {
		return nil, err
	}

	for _, e := range entries {
		if action, err := actionFromEntry(dir, e); err != nil {
			return nil, err
		} else {
			child := server.path.Add(action.Request.Path)
			child.Action = action
		}
	}

	return server, nil
}

func merge(a map[string]string, b map[string]string) map[string]string {
	for k, v := range b {
		a[k] = v
	}

	return a
}

// TODO: clean me up!
func (s *Server) HandleRequest(w http.ResponseWriter, req *http.Request) {
	url := strings.Split(req.URL.Path, "/")
	if req.Method == "" {
		url[0] = "get"
	} else {
		url[0] = strings.ToLower(req.Method)
	}

	node, groups := s.path.Find(url, nil)

	if node == nil || node.Action == nil {
		w.WriteHeader(404)
		return
	}

	for l, v := range req.Header {
		_, vars := node.Action.CompareHeaders(l, strings.Join(v, ","))
		groups = merge(groups, vars)
	}

	if body, err := io.ReadAll(req.Body); err != nil {
		w.WriteHeader(404)
		return
	} else if matched, vars := node.Action.CompareBody(string(body)); !matched {
		w.WriteHeader(404)
		return
	} else {
		groups = merge(groups, vars)
	}

	node.Action.Write(w, groups)
}
