package validator

import (
	"strconv"

	"github.com/Entrio/subenv"
)

/*
{
  "commission": {
    "commission": [
      {
        "denom": "uluna",
        "amount": "212062949.227927445774912924"
      }
    ]
  }
}
*/

type (
	DistributionCommissionResponse struct {
		Commission struct {
			Commission []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"commission"`
		} `json:"commission"`
	}

	RewardsResponse struct {
		Rewards struct {
			Rewards []struct {
				Denom  string `json:"denom"`
				Amount string `json:"amount"`
			} `json:"rewards"`
		} `json:"rewards"`
	}
)

func (d DistributionCommissionResponse) GetOutstandingCommission() float64 {
	denom := subenv.Env("DENOM_NAME", "uluna")
	for i := range d.Commission.Commission {
		if d.Commission.Commission[i].Denom != denom {
			continue
		}
		commFloat, err := strconv.ParseFloat(d.Commission.Commission[i].Amount, 64)
		if err != nil {
			return 0
		}
		return commFloat / 1000000
	}

	return 0
}

func (r RewardsResponse) GetOutstandingRewards() float64 {
	denom := subenv.Env("DENOM_NAME", "uluna")
	for i := range r.Rewards.Rewards {
		if denom != r.Rewards.Rewards[i].Denom {
			continue
		}
		rewardsFloat, err := strconv.ParseFloat(r.Rewards.Rewards[i].Amount, 64)
		if err != nil {
			return 0
		}
		return rewardsFloat / 1000000
	}

	return 0

}
