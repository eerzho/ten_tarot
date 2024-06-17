package http

import (
	"context"
	"log"
	"net/http"

	"github.com/eerzho/event_manager/config"
	"github.com/eerzho/event_manager/internal/handler/http/v1"
	"github.com/eerzho/event_manager/internal/service"
	"github.com/eerzho/event_manager/pkg/logger"
	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
}

func New(l logger.Logger, cfg *config.Config, tgUserService *service.TGUser, tgMessageService *service.TGMessage) *Server {
	router := gin.Default()

	v1.NewHandler(l, router, tgUserService, tgMessageService)

	server := &http.Server{
		Addr:    ":" + cfg.HTTP.Port,
		Handler: router,
	}

	return &Server{
		server: server,
	}
}

func (h *Server) Run() {
	const op = "./internal/app/http/server::Run"

	log.Printf("%s: http server starting at %s", op, h.server.Addr)
	err := h.server.ListenAndServe()
	if err != nil {
		log.Printf("%s: %v", op, err)
	}
}

func (h *Server) Shutdown(ctx context.Context) {
	const op = "./internal/app/app/server::Shutdown"

	log.Printf("%s: http server shutting down", op)
	err := h.server.Shutdown(ctx)
	if err != nil {
		log.Printf("%s: %v", op, err)
	}
}
