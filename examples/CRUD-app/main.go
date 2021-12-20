package main

import (
	"database/sql"
	"fmt"
	"github.com/DarkReduX/Kubernetes-WorkStart-GUIDE/tree/master/examples/crud-app/config"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	postgresConfig := config.NewPostgresConfig()
	// set connection with database

	postgresDb, err := sql.Open("postgres", fmt.Sprintf("port=%s host=%s user=%s password=%s dbname=%s sslmode=%s",
		postgresConfig.PORT,
		postgresConfig.Host,
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.DbName,
		"disable"))

	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err = postgresDb.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// init echo
	e := echo.New()

	e.GET("users", func(c echo.Context) error {
		var users []string
		row, err := postgresDb.QueryContext(c.Request().Context(), `SELECT * FROM users`)
		if err != nil {
			log.Errorf("Query error: %v", err)
			return c.String(http.StatusInternalServerError, err.Error())
		}

		for row.Next() {
			var user string
			err := row.Scan(&user)
			if err != nil {
				log.Errorf("Scan err: %v", err)
			}

			users = append(users, user)
		}

		return nil
	})

	e.POST("users", func(c echo.Context) error {
		user := c.QueryParam("user")
		if user == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		_, err := postgresDb.ExecContext(c.Request().Context(), `INSERT INTO USERS VALUES($1)`, user)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.String(http.StatusOK, fmt.Sprintf("Successful added user: %v", user))
	})

	e.DELETE("users/:user", func(c echo.Context) error {
		user := c.Param("user")
		if user == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		_, err := postgresDb.ExecContext(c.Request().Context(), `DELETE FROM USERS WHERE username = $1`, user)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.String(http.StatusOK, fmt.Sprintf("Successful deleted user: %v", user))
	})

	e.PUT("users/:user", func(c echo.Context) error {
		currUser := c.Param("newUser")
		if currUser == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		newUser := c.QueryParam("newUser")
		if newUser == "" {
			return c.NoContent(http.StatusBadRequest)
		}

		_, err := postgresDb.ExecContext(c.Request().Context(), `UPDATE USERS SET username = $1 WHERE username = $2`, newUser, currUser)
		if err != nil {
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.String(http.StatusOK, fmt.Sprintf("Successful renamed user %v to : %v", currUser, newUser))
	})

	e.Logger.Fatal(e.Start(":8080"))
}
