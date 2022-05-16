package app

import (
	"client/internal/config"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type App struct {
	router *mux.Router
	logger *logrus.Logger
	config *config.Config
}

func NewApp(config *config.Config) *App {
	app := &App{router: mux.NewRouter(), logger: logrus.New(), config: config}
	app.addRouters()
	app.logger.Infoln("App configured")
	return app
}

func (a *App) addRouters() {
	a.router.HandleFunc("/", indexHandler)
	a.router.HandleFunc("/predict", a.predictHandler).Methods("POST")
	a.router.HandleFunc("/help", helpHandler)
	a.router.HandleFunc("/download", downloadHandler)
}

func (a *App) Run() {
	s := &http.Server{
		Handler:      a.router,
		Addr:         fmt.Sprintf("%s:%s", a.config.Host, a.config.Port),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}
	a.logger.Infoln("Server started")
	s.ListenAndServe()
}
