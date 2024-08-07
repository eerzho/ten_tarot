package http

import (
	"context"
	"log"
	"net/http"

	"github.com/eerzho/ten_tarot/config"
	"github.com/eerzho/ten_tarot/internal/handler/http/v1"
	"github.com/eerzho/ten_tarot/internal/repo/mongo_repo"
	"github.com/eerzho/ten_tarot/internal/service"
	"github.com/eerzho/ten_tarot/pkg/mongo"
	"github.com/gin-gonic/gin"
)

type Server struct {
	server *http.Server
}

func New(cfg *config.Config, mg *mongo.Mongo) *Server {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// repo
	tgUserRepo := mongo_repo.NewTGUser(mg)
	tgMessageRepo := mongo_repo.NewTGMessage(mg)

	// service
	tgUserService := service.NewTGUser(tgUserRepo)
	deckService := service.NewDeck()
	tarotService := service.NewTarot(cfg.Model, cfg.GPT.Token, cfg.GPT.Prompt)
	tgMessageService := service.NewTGMessage(tgMessageRepo, deckService, tarotService)

	// handler
	v1.NewHandler(router, tgUserService, tgMessageService)

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
