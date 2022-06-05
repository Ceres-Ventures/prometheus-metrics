package main

import (
	"fmt"
	"net/http"

	"github.com/ceres-ventures/prometheus-metrics/pkg"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/metrics", getMetrics)

	e.Logger.Fatal(e.Start(":1323"))
}

func getMetrics(c echo.Context) error {
	validator, err := pkg.MakeRequest()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to query remote stats\n%s", err.Error()))
	}
}
