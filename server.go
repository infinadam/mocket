package main

import (
	"github.com/infinadam/mocket/router"
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

	var err error
	var json []byte

	if json, err = os.ReadFile(dir + "/" + e.Name()); err != nil {
		return nil, err
	}

	var action *router.HTTPAction
	if action, err = router.HTTPActionFromJSON(json); err != nil {
		return nil, err
	}

	return action, nil
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

func (s *Server) HandleRequest(w http.ResponseWriter, req *http.Request) {
	url := strings.Split(req.URL.Path, "/")
	if req.Method == "" {
		url[0] = "get"
	} else {
		url[0] = strings.ToLower(req.Method)
	}

	node := s.path.Find(url)

	if node == nil || node.Action == nil {
		w.WriteHeader(404)
		return
	}

	node.Action.Write(w)
}
