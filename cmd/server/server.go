package server

import (
	"context"
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/http-swagger"

	_ "notebook/docs"
	"notebook/internal/auth"
	"notebook/internal/config"
	db "notebook/internal/database"
	l "notebook/internal/logger"
	h "notebook/web/handlers"
)
var (
	logger = l.NewLogger()
)

type Server struct{
	ctx context.Context
	echo *echo.Echo
	cfg *config.Config
}

type Template struct{
	templates *template.Template
}

// Метод для рендера шаблонов, обращается к структур Template
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error{
	return t.templates.ExecuteTemplate(w, name, data) 
}

func (s *Server) Routes() {
	s.echo.GET("/swagger/*", echo.WrapHandler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
        httpSwagger.DocExpansion("none"),
	)))
	s.echo.GET("/", func(c echo.Context) error{
		return c.Render(http.StatusOK, "auth.html", nil)
	})

	public := s.echo.Group("/public")
	{	
		public.POST("/auth/post", auth.Authorization)
		public.GET("/auth", h.AuthPage)
		public.GET("/reg", h.RegPage)
		public.POST("/reg/post", auth.Registration)
	}
	users := s.echo.Group("/users/:user_id")
	users.Use(auth.AuthMiddleware)
	{
		users.GET("/notes", db.ShowNotes)
		users.POST("/notes/delete", db.DeleteNotes)
		users.POST("/notes/post", db.WriteNotes)
		users.GET("/info", h.UserInfoPage)
	}
}

func New(cfg *config.Config) *Server{
	ctx := context.Background()
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	return &Server{
		echo: e,
		cfg: cfg,
		ctx: ctx,
	}
}

func (s *Server) ShowTemplates(){
	tmpl, err := template.ParseGlob(
		"web/templates/*.html", 
	)
	if err != nil{
		logger.Err(err).Msg("Error in parsing HTML templates in handler!")
	}

	s.echo.Renderer = &Template{templates: tmpl}
}

func (s *Server) SetMiddleware(){
	s.echo.Use(l.LoggerMiddleware())
	s.echo.Use(middleware.Recover())
	s.echo.Use(middleware.CORS())
	s.echo.Use(middleware.Secure())
}

func (s *Server) Start() error{
	s.SetMiddleware()
	s.Routes()
	s.ShowTemplates()
	logger.Info().Msg("Starting the server")
	
	s.cfg.ShutdownTimeout = 30*time.Second
	serverError := make(chan error, 1)
	go func(){
		if err := s.echo.Start(s.cfg.ServerPort); !errors.Is(err, http.ErrServerClosed){
			logger.Error().Err(err).Msg("Error occured while starting server")
			serverError <- err
		}
		close(serverError)
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select{
		case err := <-serverError:
			return err
		case <- stop:
			logger.Info().Msg("Shutdown signal recieved")
		case <-s.ctx.Done():
			logger.Info().Msg("Context is done")
	}
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		s.cfg.ShutdownTimeout,
	)
	defer cancel()
	if err := s.echo.Shutdown(shutdownCtx); err != nil{
		if closeErr := s.echo.Close(); closeErr != nil{
			return errors.Join(err, closeErr)
		}
	}
	return nil
}