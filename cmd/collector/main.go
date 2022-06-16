package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Entrio/subenv"
	"github.com/ceres-ventures/prometheus-metrics/pkg/blockchain"
	"github.com/ceres-ventures/prometheus-metrics/pkg/external"
	"github.com/ceres-ventures/prometheus-metrics/pkg/job"
	"github.com/ceres-ventures/prometheus-metrics/pkg/validator"
	"github.com/labstack/echo/v4"
)

var (
	metricStore *blockchain.MetricStore
)

func main() {
	e := echo.New()

	dis := job.CreateNewDispatcher()
	metricStore = blockchain.NewMetricStore()
	metricStore.Start()

	r, err := blockchain.GetLatestBlockData()
	if err != nil {
		panic(err.Error())
	}

	metricStore.AddUpdate(blockchain.LatestBlockHeight, r)

	e.GET("/", func(c echo.Context) error {
		r, e := blockchain.GetLatestBlockData()
		if e != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to query\n%s", e.Error()))
		}
		return c.JSON(http.StatusOK, r)
	})

	e.GET("/metrics", getMetrics)
	latestBlocksJob := func() error {
		r, e := blockchain.GetLatestBlockData()
		if e != nil {
			return e
		}

		metricStore.AddUpdate(blockchain.LatestBlockHeight, r)
		//time.Sleep(time.Millisecond * 150)
		//metricStore.AddUpdate(blockchain.AverageTransactionsPerBlock, r.GetNumOfTxs())
		time.Sleep(time.Second * 2)
		return nil
	}
	dis.AddJob(latestBlocksJob, true, -1, 0)

	supplyJob := func() error {
		r, e := blockchain.GetLunaSupply()
		if e != nil {
			return e
		}

		metricStore.AddUpdate(blockchain.LunaSupply, r.GetLunaSupply())
		time.Sleep(time.Second * 2)
		return nil
	}
	dis.AddJob(supplyJob, true, -1, 0)

	marketJob := func() error {
		r, e := external.GetLunaMarketData()
		if e != nil {
			return e
		}

		metricStore.AddUpdate(blockchain.LunaMarketData, r)
		time.Sleep(time.Second * 20)
		return nil
	}
	dis.AddJob(marketJob, true, -1, 0)

	delegations := func() error {
		val, err := validator.GetValidatorData()
		if err != nil {
			return err
		}
		valDistr, err := validator.QueryValidatorCommissions()
		if err != nil {
			return err
		}
		valRewards, err := validator.QueryValidatorRewards()
		if err != nil {
			return err
		}

		metricStore.AddUpdate(blockchain.ValidatorTokensAllocated, val.Validator.GetTokens())
		metricStore.AddUpdate(blockchain.ValidatorOutstandingCommission, valDistr.GetOutstandingCommission())
		metricStore.AddUpdate(blockchain.ValidatorOutstandingRewards, valRewards.GetOutstandingRewards())
		time.Sleep(time.Second * 2)
		return nil
	}
	dis.AddJob(delegations, true, -1, 0)

	dis.Start(5)

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
	return c.String(200, metricStore.ToPrometheusString())
}
