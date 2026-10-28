package delivery

import (
	"net/http"
	"time"

	// Internal
	"github.com/revirathya/revirathya-telbot/internal/delivery/health"
	"github.com/revirathya/revirathya-telbot/internal/delivery/version"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
)

// Public
type ServerLayer struct {
	Server *http.Server
	Mux    *http.ServeMux
}

func Setup(
	env string,
	port string,
	buildVersion string,
) *ServerLayer {
	// Configure server
	mux := http.NewServeMux()

	// Set handler
	mux.HandleFunc("GET /healthz", health.Healthz())
	mux.HandleFunc("GET /version", version.Handler(buildVersion))

	// Set Server
	httpServer := &http.Server{
		Addr:              port,
		Handler:           mux,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	return &ServerLayer{
		Server: httpServer,
		Mux:    mux,
	}
}
