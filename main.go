package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

const (
	// shutdownTimeout is the time to give handlers to finish their work.
	shutdownTimeout = 10 * time.Second

	// readTimeout is the maximum duration for reading the entire request, including the body.
	readTimeout = 10 * time.Second

	// writeTimeout is the maximum duration before timing out writes of the response.
	writeTimeout = 10 * time.Second

	// maxHeaderBytes is the maximum number of bytes the server will read.
	maxHeaderBytes = 1 << 20

	port = 8080
)

func main() {
	r := httprouter.New()

	// /healthz/live.
	{
		h := func(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
			w.WriteHeader(http.StatusOK)
		}
		r.GET("/healthz/live", h)
		log.Info().Msg("initialized /healthz/live")
	}

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        r,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	go func() {
		listener, errListen := net.Listen("tcp", srv.Addr)
		if errListen != nil {
			log.Fatal().Err(errListen).Msg("failed to listen")
		}
		log.Info().Msgf("listening on %s", srv.Addr)
		if err := srv.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	sig := <-quit
	log.Info().Msgf("received %s, gracefully shutting down", sig.String())

	if err := shutdown(srv); err != nil {
		log.Error().Err(err).Msg("failed to shutdown server")
	}
}

func shutdown(srv *http.Server) error {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		if cErr := srv.Close(); cErr != nil {
			return fmt.Errorf("closing server: %v, %v", err, cErr)
		}

		return fmt.Errorf("shutting down server: %v", err)
	}

	return nil
}
