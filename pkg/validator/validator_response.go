package validator

import (
	"strconv"
	"strings"
	"time"
)

/*

{
  "validator": {
    "operator_address": "terravaloper1q8w4u2wyhx574m70gwe8km5za2ptanny9mnqy3",
    "consensus_pubkey": {
      "@type": "/cosmos.crypto.ed25519.PubKey",
      "key": "KeOVNq3b/t5ozaSegdqWcC369isNQ0udRX1oNXHZFCw="
    },
    "jailed": false,
    "status": "BOND_STATUS_BONDED",
    "tokens": "5711687412644",
    "delegator_shares": "5711687412644.000000000000000000",
    "description": {
      "moniker": "Ceres Ventures",
      "identity": "D11A3578D356F50A",
      "website": "https://ceres.ventures",
      "security_contact": "technical@ceres.ventures",
      "details": "Democratising Real-World Real Estate. Ceres Ventures \u0026 Terrafirma NFTs aims to change the way people think about the entire property lifecycle, bridging fractionalised real-world real-estate with the blockchain in a way that has never been done. https://twitter.com/ceres_ventures"
    },
    "unbonding_height": "0",
    "unbonding_time": "1970-01-01T00:00:00Z",
    "commission": {
      "commission_rates": {
        "rate": "0.050000000000000000",
        "max_rate": "0.200000000000000000",
        "max_change_rate": "0.010000000000000000"
      },
      "update_time": "2022-05-28T06:00:00Z"
    },
    "min_self_delegation": "1"
  }
}

*/

type (
	ValidatorResponse struct {
		Validator Validator `json:"validator"`
	}
	Validator struct {
		OperatorAddress   string               `json:"operator_address"`
		Jailed            bool                 `json:"jailed"`
		Status            string               `json:"status"`
		DelegatedAmount   string               `json:"tokens"`
		DelegatorShares   string               `json:"delegator_shares"`
		Description       ValidatorDescription `json:"description"`
		MinSelfDelegation string               `json:"min_self_delegation"`
		UnbondingHeight   string               `json:"unbonding_height"`
		Commission        ValidatorCommission  `json:"commission"`
	}

	ValidatorCommission struct {
		CommissionRates CommissionRates `json:"commission_rates"`
		UpdateTime      time.Time       `json:"update_time"`
	}
	CommissionRates struct {
		Rate          string `json:"rate"`
		MaxRate       string `json:"max_rate"`
		MaxChangeRate string `json:"max_change_rate"`
	}
	ValidatorDescription struct {
		Moniker         string `json:"moniker"`
		Identity        string `json:"identity"`
		Website         string `json:"website"`
		SecurityContact string `json:"security_contact"`
		Details         string `json:"details"`
	}
)

func (v Validator) GetShares() string {
	return strings.Split(v.DelegatorShares, ".")[0]
}

func (v Validator) GetUTokens() uint64 {
	n, err := strconv.ParseUint(v.DelegatedAmount, 10, 64)
	if err == nil {
		return n
	}
	return 0
}
func (v Validator) GetTokens() float64 {
	delegatedFloat, err := strconv.ParseFloat(v.DelegatedAmount, 64)
	if err != nil {
		return 0
	}
	return delegatedFloat / 1000000
}
