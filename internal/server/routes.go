package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	todoapplication "starter/internal/application/todo"
	userapplication "starter/internal/application/user"
	"starter/internal/eventbus"
	"starter/internal/handler"
	todopersistence "starter/internal/infrastructure/persistence/todo"
	userpersistence "starter/internal/infrastructure/persistence/user"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	gormDB := s.db.GetDB()
	api := r.Group("/api/v1")

	// Shared event bus (synchronous) with default logging handlers.
	bus := eventbus.New()
	bus.RegisterDefaultHandlers()

	// ---------------------------------------------------------------
	// DDD: User routes
	// ---------------------------------------------------------------
	{
		userRepo := userpersistence.NewUserRepository(gormDB)
		userSvc := userapplication.NewService(userRepo, bus)
		userHandler := handler.NewUserHandler(userSvc)

		users := api.Group("/users")
		users.GET("", userHandler.GetUsers)
		users.POST("", userHandler.CreateUser)
		users.GET("/:id", userHandler.GetUser)
	}

	// ---------------------------------------------------------------
	// DDD: TodoList routes
	// ---------------------------------------------------------------
	{
		todoRepo := todopersistence.NewTodoListRepository(gormDB)
		todoSvc := todoapplication.NewService(todoRepo, bus)
		todoHandler := handler.NewTodoHandler(todoSvc)

		todoRoutes := api.Group("/todo-lists")
		todoRoutes.GET("", todoHandler.GetAll)
		todoRoutes.POST("", todoHandler.Create)
		todoRoutes.GET("/:id", todoHandler.GetByID)
		todoRoutes.POST("/:id/lines", todoHandler.AddLine)
		todoRoutes.PATCH("/:id/lines/:lineId/complete", todoHandler.CompleteLine)
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"
	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
