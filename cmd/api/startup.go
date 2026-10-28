package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/delivery"
)

const (
	shutdownTimeout = 10 * time.Second
)

func RunServers(
	ctx context.Context,
	httpSrvLayer *delivery.ServerLayer,
) *errgroup.Group {
	// Init errgroup
	g, gctx := errgroup.WithContext(ctx)

	// Serve http
	g.Go(func() error {
		if err := httpSrvLayer.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	// Graceful shutdown
	g.Go(func() error {
		<-gctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		httpErr := httpSrvLayer.Server.Shutdown(shutdownCtx)
		return errors.Join(httpErr)
	})

	return g
}
