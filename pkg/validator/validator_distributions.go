package validator

import "strconv"

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
	commFloat, err := strconv.ParseFloat(d.Commission.Commission[0].Amount, 64)
	if err != nil {
		return 0
	}
	return commFloat / 1000000
}

func (r RewardsResponse) GetOutstandingRewards() float64 {
	rewardsFloat, err := strconv.ParseFloat(r.Rewards.Rewards[0].Amount, 64)
	if err != nil {
		return 0
	}
	return rewardsFloat / 1000000
}
