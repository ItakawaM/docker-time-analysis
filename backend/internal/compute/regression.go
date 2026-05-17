package compute

import "math"

type Regression interface {
	Predict(x float64) float64
}

type LinearRegression struct {
	Type             string  `json:"type"`
	AlphaCoefficient float64 `json:"alphaCoefficient"`
	BetaCoefficient  float64 `json:"betaCoefficient"`
	RSquared         float64 `json:"rSquared"`
	Qo               float64 `json:"qo"`
}

func (lf *LinearRegression) Predict(x float64) float64 {
	return lf.AlphaCoefficient*x + lf.BetaCoefficient
}

type PiecewiseRegression struct {
	Type                  string  `json:"type"`
	Breakpoint            float64 `json:"breakpoint"`
	LeftAlphaCoefficient  float64 `json:"leftAlphaCoefficient"`
	LeftBetaCoefficient   float64 `json:"leftBetaCoefficient"`
	RightAlphaCoefficient float64 `json:"rightAlphaCoefficient"`
	RightBetaCoefficient  float64 `json:"rightBetaCoefficient"`
	RSquared              float64 `json:"rSquared"`
	Qo                    float64 `json:"qo"`
}

func (pf PiecewiseRegression) Predict(x float64) float64 {
	if x <= pf.Breakpoint {
		return pf.LeftAlphaCoefficient*x + pf.LeftBetaCoefficient
	}

	return pf.RightAlphaCoefficient*x + pf.RightBetaCoefficient
}

type ExponentialRegression struct {
	Type             string  `json:"type"`
	AlphaCoefficient float64 `json:"alphaCoefficient"`
	BetaCoefficient  float64 `json:"betaCoefficient"`
	RSquared         float64 `json:"rSquared"`
	Qo               float64 `json:"qo"`
}

func (ef *ExponentialRegression) Predict(x float64) float64 {
	return ef.BetaCoefficient * math.Pow(ef.AlphaCoefficient, x)
}
