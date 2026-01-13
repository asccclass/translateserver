// router.go
package main

import (
	"net/http"
	"translateserver/handler"

	SherryServer "github.com/asccclass/sherryserver"
)

// Create your Router function
func NewRouter(srv *SherryServer.Server, documentRoot string) *http.ServeMux {
	router := http.NewServeMux()

	// Static File server
	staticfileserver := SherryServer.StaticFileServer{
		StaticPath: documentRoot,
		IndexPath:  "index.html",
	}
	staticfileserver.AddRouter(router)

	// Routes
	router.HandleFunc("/upload", handler.UploadHandler)
	router.HandleFunc("/ws", handler.WebSocketHandler)

	return router
}
