package blockchain

import "fmt"

type UpdateType int64
type UpdateField int64

const (
	ValidatorTokensAllocated UpdateField = iota
	ValidatorOutstandingCommission
	ValidatorOutstandingRewards
	LatestBlockHeight
	LunaSupply
	AverageTransactionsPerBlock
)

const (
	Update = iota
	Reset
)

type (
	MetricsData struct {
		ValidatorTokensAllocated       float64
		ValidatorOutstandingCommission float64
		ValidatorOutstandingRewards    float64
		LatestBlockHeight              int
		LunaSupply                     float64
		AverageTransactionsPerBlock    float64
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
	return fmt.Sprintf(
		"# HELP validator_tokens_allocated The current validator delegations.\n# TYPE validator_tokens_allocated gauge\nvalidator_tokens_allocated %f\n"+
			"# HELP validator_outstanding_commission Current outstanding commision\n# TYPE validator_outstanding_commission gauge\nvalidator_outstanding_commission %f\n"+
			"# HELP validator_outstanding_rewards Current outstanding commision\n# TYPE validator_outstanding_rewards gauge\nvalidator_outstanding_rewards %f\n"+
			"# HELP latest_block_height Most recent block\n# TYPE latest_block_height gauge\nlatest_block_height %d\n"+
			"# HELP luna_circulation_supply Current luna in circulation\n# TYPE luna_circulation_supply counter\nluna_circulation_supply %f\n"+
			"# HELP avg_transaction_per_block Average number of transaction for last 10 blocks\n# TYPE avg_transaction_per_block gauge\navg_transaction_per_block %f\n"+
			"# HELP transactions_last_block Transactions for last block\n# TYPE transactions_last_block gauge\ntransactions_last_block %d\n"+
			"",
		ms.Data.ValidatorTokensAllocated,
		ms.Data.ValidatorOutstandingCommission,
		ms.Data.ValidatorOutstandingRewards,
		ms.cursor[0][ms.currentCursor],
		ms.Data.LunaSupply,
		ms.Data.AverageTransactionsPerBlock,
		ms.cursor[1][ms.currentCursor],
	)
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
		height := value.(int)
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