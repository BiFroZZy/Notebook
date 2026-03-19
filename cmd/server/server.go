package server

import (
	"html/template"
	"io"
	"net/http"

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
		users.GET("/about", h.AboutPage)
		users.GET("/notes", db.ShowNotes)
		users.POST("/notes/delete", db.DeleteNotes)
		users.POST("/notes/post", db.WriteNotes)
		users.GET("/info", h.UserInfoPage)
	}
}

func New(cfg *config.Config) *Server{
	e := echo.New()
	e.HideBanner = true
	return &Server{
		echo: e,
		cfg: cfg,
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
	return s.echo.Start(s.cfg.ServerPort)
}