package msetups

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/yankeguo/minit/internal/mlog"
	"golang.org/x/net/webdav"
)

func init() {
	Register(50, setupWebDAV)
}

func setupWebDAV(logger mlog.ProcLogger) (err error) {
	envRoot := strings.TrimSpace(os.Getenv("MINIT_WEBDAV_ROOT"))
	if envRoot == "" {
		return
	}
	if err = os.MkdirAll(envRoot, 0755); err != nil {
		err = fmt.Errorf("failed initializing WebDAV root %s: %s", envRoot, err.Error())
		return
	}
	envPort := strings.TrimSpace(os.Getenv("MINIT_WEBDAV_PORT"))
	if envPort == "" {
		envPort = "7486"
	}
	logger.Printf("WebDAV started: root=%s, port=%s", envRoot, envPort)
	h := &webdav.Handler{
		FileSystem: webdav.Dir(envRoot),
		LockSystem: webdav.NewMemLS(),
		Logger: func(req *http.Request, err error) {
			if err != nil {
				logger.Printf("WebDAV: %s %s: %s", req.Method, req.URL.Path, err.Error())
			} else {
				logger.Printf("WebDAV: %s %s", req.Method, req.URL.Path)
			}
		},
	}
	envUsername := strings.TrimSpace(os.Getenv("MINIT_WEBDAV_USERNAME"))
	envPassword := strings.TrimSpace(os.Getenv("MINIT_WEBDAV_PASSWORD"))
	s := &http.Server{
		Addr: ":" + envPort,
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if envUsername != "" && envPassword != "" {
				if username, password, ok := req.BasicAuth(); !ok || username != envUsername || password != envPassword {
					rw.Header().Add("WWW-Authenticate", `Basic realm=Minit WebDAV`)
					rw.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
			h.ServeHTTP(rw, req)
		}),
	}
	// Start WebDAV server in background goroutine with error handling
	go func() {
		// Only retry on unexpected errors, not on graceful shutdown
		for {
			err := s.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				logger.Printf("WebDAV server error: %s", err.Error())
				// Wait before retrying to avoid tight loop on persistent errors
				time.Sleep(time.Second * 10)
			} else {
				// Server closed gracefully or shutting down
				if err == http.ErrServerClosed {
					logger.Printf("WebDAV server shut down")
				}
				return
			}
		}
	}()
	return
}
