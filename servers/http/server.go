package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/weni/whatsapp-router/config"
	"github.com/weni/whatsapp-router/logger"
	"github.com/weni/whatsapp-router/repositories"
	"github.com/weni/whatsapp-router/servers/http/handlers"
	"github.com/weni/whatsapp-router/services"
	"go.mongodb.org/mongo-driver/mongo"
)

type Server struct {
	config     config.Config
	db         *mongo.Database
	httpServer *http.Server
}

func NewServer(db *mongo.Database) *Server {
	conf := config.GetConfig()
	return &Server{
		db:     db,
		config: *conf,
	}
}

func (s *Server) Start() error {
	sRouter := NewRouter(s)
	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", s.config.Server.HttpPort),
		Handler:      sRouter,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("Starting http server :%v", s.config.Server.HttpPort))
		err := s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error(err.Error())
			log.Fatal()
		}
	}()

	return nil
}

func NewRouter(s *Server) *chi.Mux {
	router := chi.NewRouter()

	contactRepoDb := repositories.NewContactRepositoryDb(s.db)
	channelRepoDb := repositories.NewChannelRepositoryDb(s.db)
	whatsappHandler := handlers.WhatsappHandler{
		ContactService:  services.NewContactService(contactRepoDb),
		ChannelService:  services.NewChannelService(channelRepoDb),
		CourierService:  services.NewCourierService(),
		WhatsappService: services.NewWhatsappService(),
	}
	courierHandler := handlers.CourierHandler{
		WhatsappService: services.NewWhatsappService(),
	}

	router.Route("/wr/", func(r chi.Router) {
		r.Use(ContentTypeJson)
		r.Route("/receive", func(r chi.Router) {
			r.Post("/", whatsappHandler.HandleIncomingRequests)
		})
	})

	router.Route("/v1", func(r chi.Router) {
		r.Post("/messages", courierHandler.HandleSendMessage)
		r.Post("/users/login", handlers.RefreshToken)
		r.Get("/settings/application", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return router
}

func ContentTypeJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json;charset=utf8")
		next.ServeHTTP(w, r)
	})
}
