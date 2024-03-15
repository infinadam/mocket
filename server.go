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

func MakeServer(dir string) (*Server, error) {
	var entries []os.DirEntry
	var err error
	server := new(Server)

	if entries, err = os.ReadDir(dir); err != nil {
		return nil, err
	}

	var json []byte
	for _, e := range entries {
		if !e.IsDir() {
			if json, err = os.ReadFile(dir + "/" + e.Name()); err != nil {
				return nil, err
			}
			action, err := router.HTTPActionFromJSON(string(json))
			if err != nil {
				return nil, err
			}
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
