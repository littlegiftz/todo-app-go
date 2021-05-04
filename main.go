package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/littlegiftz/todo-app-go/controller"
	"github.com/littlegiftz/todo-app-go/db"
	"github.com/littlegiftz/todo-app-go/model"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	// Initilize database
	db.Init()

	jwtMiddleware := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &model.CustomClaims{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	})

	// Create routes
	e := echo.New()

	// Middleware
	e.HTTPErrorHandler = customHTTPErrorHandler
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.POST("/login", controller.Login)
	e.POST("/register", controller.CreateUser)

	u := e.Group("/user")
	u.Use(jwtMiddleware)
	u.POST("/change-password", controller.SavePassword)

	t := e.Group("/todo")
	t.Use(jwtMiddleware)
	t.GET("/", controller.GetTasks)
	t.POST("/", controller.AddTask)
	t.POST("/:id", controller.SaveTask)
	t.DELETE("/:id", controller.DeleteTask)

	e.Logger.Fatal(e.Start(":1323"))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	if err == middleware.ErrJWTMissing {
		c.Error(echo.NewHTTPError(http.StatusUnauthorized, "Login required"))
		return
	}
	c.Echo().DefaultHTTPErrorHandler(err, c)
}
