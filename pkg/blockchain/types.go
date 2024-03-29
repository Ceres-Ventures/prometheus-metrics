package blockchain

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/ceres-ventures/prometheus-metrics/pkg/external"
)

type UpdateType int64
type UpdateField int64

const (
	ValidatorTokensAllocated UpdateField = iota
	ValidatorOutstandingCommission
	ValidatorOutstandingRewards
	LatestBlockHeight
	LunaSupply
	AverageTransactionsPerBlock
	LunaMarketData
	BlockSign
	Status
	WalletBalances
)

const (
	Update = iota
	Reset
)

type (
	MetricsData struct {
		ChainID                     string
		LastBlockHeight             int
		LatestBlockHeight           int
		LunaSupply                  float64
		AverageTransactionsPerBlock float64
		MarketData                  external.CGPriceResponse
		LastSignedCount             int64
		LastBlocksSinceSigned       int64
		LastSignedBlockHeight       int
		WalletData                  []WalletData
		ValidatorData               []ValidatorData
	}

	WalletData struct {
		Address string  `json:"address"`
		ID      string  `json:"ID"`
		Balance float64 `json:"balance"`
	}

	ValidatorData struct {
		AccountAddress        string  `json:"account_address"`
		HexAddress            string  `json:"hex_address"`
		OperatorAddress       string  `json:"operator_address"`
		PublicKey             string  `json:"public_key"`
		DelegationAmount      float64 `json:"delegation_amount"`
		CommissionOutstanding float64 `json:"commission_outstanding"`
		RewardsOutstanding    float64 `json:"rewards_outstanding"`
		VotingPower           float64 `json:"voting_power"`
	}

	Request struct {
		Method UpdateType
		Field  UpdateField
		Value  interface{}
	}

	MetricStore struct {
		currentCursor int
		cursor        [2][10]int // Dynamic cursors for keeping track of things
		Data          MetricsData
		UpdateChan    chan *Request
		Quit          chan bool
	}
)

func (ms MetricStore) ToPrometheusString() string {
	stats, market, validators, wallets := strings.Builder{}, strings.Builder{}, strings.Builder{}, strings.Builder{}

	market.WriteString(
		fmt.Sprintf(
			"# TYPE luna_circulation_supply counter\nluna_circulation_supply %f\n"+
				"# TYPE market_price gauge\nmarket_price{currency=\"usd\"} %f\n"+
				"# TYPE market_price_ath gauge\nmarket_price_ath{currency=\"usd\"} %f\n"+
				"# TYPE market_price_atl gauge\nmarket_price_atl{currency=\"usd\"} %f\n",
			ms.Data.LunaSupply,
			ms.Data.MarketData.MarketData.CurrentPrice.USD,
			ms.Data.MarketData.MarketData.ATH.USD,
			ms.Data.MarketData.MarketData.ATL.USD,
		),
	)

	stats.WriteString(
		fmt.Sprintf(
			"# HELP validator_tokens_allocated The current validator delegations.\n# TYPE validator_tokens_allocated gauge\nvalidator_tokens_allocated %f\n"+
				"# HELP validator_outstanding_commission Current outstanding commision\n# TYPE validator_outstanding_commission gauge\nvalidator_outstanding_commission %f\n"+
				"# HELP validator_outstanding_rewards Current outstanding commision\n# TYPE validator_outstanding_rewards gauge\nvalidator_outstanding_rewards %f\n"+
				"# HELP latest_block_height Most recent block\n# TYPE latest_block_height gauge\nlatest_block_height{chain_id=\"%s\"} %d\n"+
				"# HELP avg_transaction_per_block Average number of transaction for last 10 blocks\n# TYPE avg_transaction_per_block gauge\navg_transaction_per_block %f\n"+
				"# HELP transactions_last_block Transactions for last block\n# TYPE transactions_last_block gauge\ntransactions_last_block %d\n"+
				"market_price{currency=\"usd\"} %f\n"+
				"# HELP market_price_ath All time high price\n# TYPE market_price_ath gauge\nmarket_price_ath{currency=\"aud\"} %f\n"+
				"market_price_ath{currency=\"usd\"} %f\n"+
				"# HELP market_price_atl All time low price\n# TYPE market_price_atl gauge\nmarket_price_atl{currency=\"aud\"} %f\n"+
				"market_price_atl{currency=\"usd\"} %f\n"+
				"# HELP market_price_high24 All time low price\n# TYPE market_price_high24 gauge\nmarket_price_high24{currency=\"aud\"} %f\n"+
				"market_price_high24{currency=\"usd\"} %f\n"+
				"# HELP market_price_low24 All time low price\n# TYPE market_price_low24 gauge\nmarket_price_low24{currency=\"aud\"} %f\n"+
				"market_price_low24{currency=\"usd\"} %f\n"+
				"# HELP block_since_sign How many blocks ago we signed a block\n# TYPE block_since_sign gauge\nblock_since_sign  %d\n"+
				"# HELP last_signed_height What was the last block we signed\n# TYPE last_signed_height gauge\nlast_signed_height  %d\n"+
				"# HELP current_luna_worth How much commission is worth\n# TYPE current_luna_worth gauge\ncurrent_luna_worth{currency=\"usd\"}  %f\n"+
				"# HELP current_luna_worth How much commission is worth\n# TYPE current_luna_worth gauge\ncurrent_luna_worth{currency=\"aud\"}  %f\n"+
				"# HELP voting_power Voting power the validator has\n# TYPE voting_power gauge\nvoting_power %f\n"+
				"# HELP wallet_balance Current wallet balance\n# TYPE wallet_balance gauge\nwallet_balance{wallet=\"%s\"} %f\n"+
				"",

			ms.Data.ValidatorTokensAllocated,
			ms.Data.ValidatorOutstandingCommission,
			ms.Data.ValidatorOutstandingRewards,
			ms.Data.ChainID,
			ms.cursor[0][ms.currentCursor],
			ms.Data.AverageTransactionsPerBlock,
			ms.cursor[1][ms.currentCursor],
			ms.Data.MarketData.MarketData.CurrentPrice.AUD,
			ms.Data.MarketData.MarketData.CurrentPrice.USD,

			ms.Data.MarketData.MarketData.ATH.AUD,
			ms.Data.MarketData.MarketData.ATH.USD,

			ms.Data.MarketData.MarketData.ATL.AUD,
			ms.Data.MarketData.MarketData.ATL.USD,

			ms.Data.MarketData.MarketData.High24.AUD,
			ms.Data.MarketData.MarketData.High24.USD,

			ms.Data.MarketData.MarketData.Low24.AUD,
			ms.Data.MarketData.MarketData.Low24.USD,

			ms.Data.LastSignedCount,
			ms.Data.LastSignedBlockHeight,
			ms.Data.ValidatorOutstandingCommission*ms.Data.MarketData.MarketData.CurrentPrice.USD,
			ms.Data.ValidatorOutstandingCommission*ms.Data.MarketData.MarketData.CurrentPrice.AUD,

			ms.Data.ValidatorVotingPower,
			ms.Data.WalletData.ID,
			ms.Data.WalletData.Balance,
		),
	)

	return stats.String()
}

func NewMetricStore() *MetricStore {
	s := &MetricStore{
		currentCursor: -1,
		cursor:        [2][10]int{{-1, -1, -1, -1, -1, -1, -1, -1, -1, -1}, {-1, -1, -1, -1, -1, -1, -1, -1, -1, -1}}, // Keep track of last 10 block heights and the num of tx per block
		Data:          MetricsData{},
		UpdateChan:    make(chan *Request),
		Quit:          make(chan bool),
	}
	s.Start()
	return s
}

func (ms *MetricStore) Start() {
	go func() {
		// Start chan thread for updates etc
		for {
			select {
			case u := <-ms.UpdateChan:
				switch u.Method {
				case Update:
					// we want to update the value
					ms.processUpdate(u.Field, u.Value)
				case Reset:
					// we want to reset the value
					ms.processReset(u.Field)
				}
			case <-ms.Quit:
				return
			}
		}
	}()
}

func (ms *MetricStore) AddUpdate(field UpdateField, value interface{}) {
	go func() {
		ms.UpdateChan <- &Request{
			Method: Update,
			Field:  field,
			Value:  value,
		}
	}()
}
func (ms *MetricStore) AddReset(field UpdateField) {
	go func() {
		ms.UpdateChan <- &Request{
			Method: Reset,
			Field:  field,
			Value:  nil,
		}
	}()
}

func (ms *MetricStore) processUpdate(field UpdateField, value interface{}) {
	switch field {
	case ValidatorTokensAllocated:
		ms.Data.ValidatorTokensAllocated = value.(float64)
	case ValidatorOutstandingCommission:
		ms.Data.ValidatorOutstandingCommission = value.(float64)
	case ValidatorOutstandingRewards:
		ms.Data.ValidatorOutstandingRewards = value.(float64)
	case LatestBlockHeight:
		bd := value.(*LatestBlockResponse)
		height := bd.GetBlockHeightInt()
		ms.Data.ChainID = bd.Block.Header.Chain

		if ms.currentCursor == -1 {
			ms.currentCursor = 0
			ms.cursor[0][ms.currentCursor] = height // set block height
		} else {
			if ms.cursor[0][ms.currentCursor] != height {
				if ms.currentCursor+1 == 10 {
					ms.currentCursor = 0
				} else {
					ms.currentCursor++
				}
				// move the cursor by one and set the value
				ms.cursor[0][ms.currentCursor] = height
			}
		}

		txs := bd.GetNumOfTxs()
		ms.cursor[1][ms.currentCursor] = txs

		// need to recalculate running average
		total := 0
		for i := 0; i < 10; i++ {
			if ms.cursor[1][i] > 0 {
				total += ms.cursor[1][i]
			}
		}
		ms.Data.AverageTransactionsPerBlock = float64(total) / 10

		if ms.Data.LastBlockHeight != height {
			logrus.Infof("New Block: %d", height)
			if bd.Block.Header.Proposer == OwnValidatorHashAddress {
				// We signed this block
				ms.Data.LastSignedBlockHeight = height
				ms.Data.LastSignedCount = ms.Data.LastBlocksSinceSigned
				ms.Data.LastBlocksSinceSigned = 0
				logrus.Infof("Block %d signed!", height)
			} else {
				// we are not the block signer
				logrus.Debugf("Proposer %s != %s. Havent signed for %d\n", bd.Block.Header.Proposer, OwnValidatorHashAddress, ms.Data.LastBlocksSinceSigned)
				ms.Data.LastBlocksSinceSigned++
			}
		}
		ms.Data.LastBlockHeight = height
	case LunaSupply:
		ms.Data.LunaSupply = value.(float64)
	case AverageTransactionsPerBlock:
		if ms.currentCursor == -1 {
			ms.currentCursor = 0
		}
		txs := value.(int)
		ms.cursor[1][ms.currentCursor] = txs

		// need to recalculate running average
		total := 0
		for i := 0; i < 10; i++ {
			if ms.cursor[1][i] > 0 {
				total += ms.cursor[1][i]
			}
		}
		ms.Data.AverageTransactionsPerBlock = float64(total) / 10
	case LunaMarketData:
		data := value.(*external.CGPriceResponse)
		if data == nil {
			logrus.Warn("Market data had no error but got a nil value. Most likely roo many requests!")
			logrus.Info(ms.Data.MarketData)
			return
		}
		ms.Data.MarketData = *data
		fmt.Println(ms.Data.MarketData)
	case BlockSign:
	case Status:
		res := value.(*StatusResponse)
		f, _ := strconv.ParseFloat(res.Result.ValidatorInfo.VotingPower, 64)
		ms.Data.ValidatorVotingPower = f / 1000000
	case WalletBalances:
		res := value.([]*BalanceResponse)

		for _, wallet := range res {
			var data WalletData
			data.Address = wallet.WalletAddress

			if len(wallet.Balances) == 0 {
				// Could do some more cleanup but most likely answer is that the balance of this denom is 0
				data.Balance = 0
				data.ID = wallet.WalletAddress
				return
			}

			bal, _ := strconv.ParseFloat(wallet.Balances[0].Amount, 64)
			data.Balance = bal / 1000000
			ms.Data.WalletData = append(ms.Data.WalletData, data)
		}

	}
}
func (ms *MetricStore) processReset(field UpdateField) {
	fmt.Printf("Processing reset for %v\n", field)
	switch field {
	case ValidatorTokensAllocated:
		ms.Data.ValidatorTokensAllocated = 0.0
	case ValidatorOutstandingCommission:
		ms.Data.ValidatorOutstandingCommission = 0.0
	case ValidatorOutstandingRewards:
		ms.Data.ValidatorOutstandingRewards = 0.0
	case LatestBlockHeight:
		ms.Data.LatestBlockHeight = 0
	case LunaSupply:
		ms.Data.LunaSupply = 0.0
	case AverageTransactionsPerBlock:
		ms.Data.AverageTransactionsPerBlock = 0.0
	}
}
