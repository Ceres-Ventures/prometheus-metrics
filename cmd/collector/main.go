package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

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

//TODO: Load rpc data for validators from: https://rpc-terra.wildsage.io/validators?height=6581264 <-- Pub key + HEX
// https://lcd-terra.wildsage.io/cosmos/staking/v1beta1/validators  <-- Pub key + address
// Create a global validator storage that gets updated all the time

func main() {
	lvl, err := log.ParseLevel(subenv.Env("LOG_LEVEL", "info"))
	if err != nil {
		log.Warn("failed to parse log level, reverting to ino")
		lvl = log.InfoLevel
	}
	log.SetLevel(lvl)
	e := echo.New()
	e.HideBanner = true

	dis := job.CreateNewDispatcher()
	metricStore = blockchain.NewMetricStore()
	metricStore.Start()

	latestBlockData, err := blockchain.GetLatestBlockData()
	if err != nil {
		panic(err.Error())
	}
	metricStore.AddUpdate(blockchain.LatestBlockHeight, latestBlockData)

	/*
		s, err := blockchain.GetStatus()
		if err != nil {
			panic(err.Error())
		}
		metricStore.AddUpdate(blockchain.Status, s)
	*/

	e.GET("/", func(c echo.Context) error {
		r, e := blockchain.GetLatestBlockData()
		if e != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("failed to query\n%s", e.Error()))
		}
		return c.JSON(http.StatusOK, r)
	})

	e.GET("/metrics", getMetrics)
	latestBlocksJob := func() error {
		latestBlockData, e := blockchain.GetLatestBlockData()
		if e != nil {
			return e
		}

		metricStore.AddUpdate(blockchain.LatestBlockHeight, latestBlockData)

		balances, e := blockchain.GetBalances()
		if e != nil {
			return e
		}
		metricStore.AddUpdate(blockchain.WalletBalances, balances)

		/*
			s, err := blockchain.GetStatus()
			if err != nil {
				panic(err.Error())
			}
			metricStore.AddUpdate(blockchain.Status, s)
		*/

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
		time.Sleep(time.Second * 30)
		return nil
	}
	dis.AddJob(marketJob, true, -1, 0)

	delegations := func() error {
		validators := validator.GetValidatorData()

		valDistr, err := validator.QueryValidatorCommissions()
		if err != nil {
			return err
		}
		valRewards, err := validator.QueryValidatorRewards()
		if err != nil {
			return err
		}

		metricStore.AddUpdate(blockchain.ValidatorTokensAllocated, validators.Validator.GetTokens())
		metricStore.AddUpdate(blockchain.ValidatorOutstandingCommission, valDistr.GetOutstandingCommission())
		metricStore.AddUpdate(blockchain.ValidatorOutstandingRewards, valRewards.GetOutstandingRewards())
		time.Sleep(time.Second * 2)
		return nil
	}
	dis.AddJob(delegations, true, -1, 0)

	dis.Start(5)

	log.WithFields(
		log.Fields{
			"port": subenv.EnvI("BIND_PORT", 9292),
			"ip":   subenv.Env("BIND_IP", "0.0.0.0"),
			"LCD":  subenv.Env("LCD_URL", "http://188.40.140.51:1317"),
			"HASH": subenv.Env("VALIDATOR_KEY_ADDRESS", "C24A7D204E0A07736EAF8A7E76820CD868565B0E"),
		},
	).Info("Starting metrics server")

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
