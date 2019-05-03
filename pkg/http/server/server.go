package server

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

type (

	// Instance is a listening HTTP server instance
	Instance struct {
		Addr    string
		Handler http.Handler
		Logger  logrus.FieldLogger
		Closer  io.Closer

		DisableGracefulShutdown bool
		GracefulShutdownDelay   time.Duration
		AddrAssigned            chan string

		listeningError chan error
	}
)

// Start begins listening and serving
func (i *Instance) Start() {
	httpServer := &http.Server{
		Handler: i.Handler,
	}

	if i.Addr == "" {
		i.Addr = ":0"
	}

	listener, err := net.Listen("tcp", i.Addr)
	if err != nil {
		if i.Logger != nil {
			i.Logger.Fatalf("Failed to start server at %s", i.Addr)
		}

		return
	}

	if !i.DisableGracefulShutdown {
		i.handleServerShutdown(httpServer)
	}

	i.listeningError = make(chan error)
	go func() {
		if serveErr := httpServer.Serve(listener); serveErr != nil {
			if serveErr == http.ErrServerClosed {
				// server is shutting down, handleServerShutdown should have handled this already
				i.listeningError <- nil

				return
			}
			i.listeningError <- serveErr

			return
		}
	}()
	runtime.Gosched()

	if i.AddrAssigned != nil {
		i.AddrAssigned <- listener.Addr().String()
	}

	err = <-i.listeningError
	if err != nil && i.Logger != nil {
		i.Logger.Fatalf("Failed to start server as %s: %s", listener.Addr().String(), err.Error())
	}
}

// handleServerShutdown will handle the shutdown signal that comes to the server
// and gracefully shutdown the server
func (i *Instance) handleServerShutdown(server *http.Server) {
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGTERM, os.Interrupt)

	go func() {
		<-exit

		gracefulDelay := i.GracefulShutdownDelay
		if gracefulDelay == 0 {
			gracefulDelay = 5 * time.Second
		}

		// Killing with no mercy after a graceful delay
		go func() {
			time.Sleep(gracefulDelay)
			if i.Logger != nil {
				i.Logger.Errorf("Failed to gracefully shutdown server in %s, exiting.", gracefulDelay)
			}
			os.Exit(0)
		}()

		if i.Logger != nil {
			i.Logger.Info("Shutting down server.")
		}

		if err := server.Shutdown(context.Background()); err != nil {
			if i.Logger != nil {
				i.Logger.Errorf("Failed to shutdown server gracefully: %s", err.Error())
			}

			err = server.Close()
			if err != nil {
				panic(err.Error())
			}
		}

		if i.Closer != nil {
			err := i.Closer.Close()
			if err != nil {
				if i.Logger != nil {
					i.Logger.Errorf("Closer failed: %s", err.Error())
				}

				panic(err.Error())
			}
		}
	}()
}
