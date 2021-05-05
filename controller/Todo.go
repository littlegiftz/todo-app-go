package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/littlegiftz/todo-app-go/db"
	"github.com/littlegiftz/todo-app-go/model"
)

type jwtCustomClaims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

func parseDate(str string) (time.Time, error) {
	layout := "2006-01-02"
	return time.Parse(layout, str)
}

func AddTask(c echo.Context) error {
	db := db.DBManager()

	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*model.CustomClaims)
	uid := claims.ID

	duedate, err := parseDate(c.FormValue("duedate"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	todo := model.Todo{
		UserID:    uid,
		Task:      c.FormValue("task"),
		DueDate:   duedate,
		Completed: false,
	}

	create := db.Create(&todo)
	if create.Error != nil {
		return c.JSON(http.StatusInternalServerError, create.Error)
	}

	return c.JSON(http.StatusOK, todo)
}

func SaveTask(c echo.Context) error {
	db := db.DBManager()
	todo := model.Todo{}

	id := c.Param("id")
	find := db.First(&todo, id)
	if find.Error != nil {
		return c.JSON(http.StatusInternalServerError, find.Error)
	}

	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*model.CustomClaims)
	uid := claims.ID

	if uid != todo.UserID {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Unauthorized with currrent token"})
	}

	duedate := c.FormValue("duedate")
	task := c.FormValue("task")
	completed := c.FormValue("completed")

	if duedate != "" {
		duedate, err := parseDate(duedate)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		todo.DueDate = duedate
	}

	if task != "" {
		todo.Task = task
	}

	if completed != "" {
		completed, _ := strconv.ParseBool(completed)
		todo.Completed = completed
	}

	save := db.Save(&todo)
	if save.Error != nil {
		return c.JSON(http.StatusInternalServerError, save.Error)
	}

	return c.JSON(http.StatusOK, todo)
}

func DeleteTask(c echo.Context) error {
	db := db.DBManager()
	todo := model.Todo{}

	id := c.Param("id")
	find := db.First(&todo, id)
	if find.Error != nil {
		return c.JSON(http.StatusInternalServerError, find.Error)
	}

	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*model.CustomClaims)
	uid := claims.ID

	if uid != todo.UserID {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Unauthorized with current token"})
	}

	delete := db.Delete(&todo)
	if delete.Error != nil {
		return c.JSON(http.StatusInternalServerError, delete.Error)
	}

	return c.NoContent(http.StatusOK)
}

func GetTasks(c echo.Context) error {
	db := db.DBManager()
	todos := []model.Todo{}

	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*model.CustomClaims)
	uid := claims.ID

	find := db.Where("user_id = ?", uid).Order("id desc").Find(&todos)
	if find.Error != nil {
		return c.JSON(http.StatusInternalServerError, find.Error)
	}

	return c.JSON(http.StatusOK, todos)
}
