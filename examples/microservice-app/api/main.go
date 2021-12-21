package main

import (
	pb "api/protocol"
	"fmt"
	"github.com/caarlos0/env"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net/http"
)

type NameServiceConfig struct {
	Host string `env:"NAME_HOST"`
	Port int    `env:"NAME_PORT,required"`
}

func main() {
	var nameCfg NameServiceConfig

	err := env.Parse(&nameCfg)
	if err != nil {
		log.Fatalf("Couldn't parse env: %v", err)
	}

	e := echo.New()

	nameServConn, err := grpc.Dial(fmt.Sprintf("%v:%v", nameCfg.Host, nameCfg.Port), grpc.WithInsecure())
	if err != nil {
		log.Errorf("Couldn't connect to name service server: %v", err)
	}

	nameServiceClient := pb.NewNameServiceClient(nameServConn)

	e.GET("/", func(c echo.Context) error {
		name := c.QueryParam("name")

		resp, err := nameServiceClient.NameToUpperCase(c.Request().Context(), &pb.NameReq{Name: name})
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusOK, resp.Name)
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":8080")))
}
