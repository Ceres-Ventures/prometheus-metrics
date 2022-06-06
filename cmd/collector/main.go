package main

import (
	"fmt"
	"net/http"

	"github.com/Entrio/subenv"
	"github.com/ceres-ventures/prometheus-metrics/pkg/validator"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/metrics", getMetrics)

	e.Logger.Fatal(
		e.Start(
			fmt.Sprintf(
				"%s:%d",
				subenv.Env("BIND_IP", "0.0.0.0"),
				subenv.EnvI("BIND_PORT", 9292),
			),
		),
	)
}

func getMetrics(c echo.Context) error {
	val, err := validator.MakeRequest()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to query remote delegations stats\n%s", err.Error()))
	}

	valDistr, err := validator.QueryValidatorCommissions()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to query remote commission stats\n%s", err.Error()))
	}

	valRewards, err := validator.QueryValidatorRewards()
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to query remote commission stats\n%s", err.Error()))
	}

	return c.String(
		200,
		fmt.Sprintf(
			"# HELP validator_tokens_allocated The current validator delegations.\n# TYPE validator_tokens_allocated gauge\nvalidator_tokens_allocated %f\n"+
				"# HELP validator_outstanding_commission Current outstanding commision\n# TYPE validator_outstanding_commission gauge\nvalidator_outstanding_commission %f\n"+
				"# HELP validator_outstanding_rewards Current outstanding commision\n# TYPE validator_outstanding_rewards gauge\nvalidator_outstanding_rewards %f\n"+
				"",
			val.Validator.GetTokens(), valDistr.GetOutstandingCommission(), valRewards.GetOutstandingRewards(),
		),
	)
}
