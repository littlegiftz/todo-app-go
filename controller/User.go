package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"

	"github.com/littlegiftz/todo-app-go/db"
	"github.com/littlegiftz/todo-app-go/model"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func Login(c echo.Context) error {
	db := db.DBManager()
	user := model.User{}

	find := db.Where("email = ?", c.FormValue("email")).First(&user)
	if find.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid email or password"})
	}

	if !VerifyPassword(user.Password, c.FormValue("password")) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid email or password"})
	}

	// Create jwt token
	claims := &model.CustomClaims{
		ID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"token": tokenString})
}

func CreateUser(c echo.Context) error {
	db := db.DBManager()

	hashedPassword, err := HashPassword(c.FormValue("password"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	user := model.User{
		Email:    c.FormValue("email"),
		Password: string(hashedPassword),
	}

	create := db.Create(&user)
	if create.Error != nil {
		return c.JSON(http.StatusInternalServerError, create.Error)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": user.ID})
}

func SavePassword(c echo.Context) error {
	db := db.DBManager()
	user := model.User{}

	u := c.Get("user").(*jwt.Token)
	claims := u.Claims.(*model.CustomClaims)
	uid := claims.ID

	find := db.First(&user, uid)
	if find.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Invalid request"})
	}

	if !VerifyPassword(user.Password, c.FormValue("old_password")) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"message": "Current password is incorrect"})
	}

	newHashedPassword, err := HashPassword(c.FormValue("new_password"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	user.Password = string(newHashedPassword)

	save := db.Save(&user)
	if save.Error != nil {
		return c.JSON(http.StatusInternalServerError, save.Error)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"id": user.ID})
}
