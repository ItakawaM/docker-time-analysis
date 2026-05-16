package compute

import "math"

type Regression interface {
	Predict(x float64) float64
}

type LinearRegression struct {
	AlphaCoefficient float64 `json:"alphaCoefficient"`
	BetaCoefficient  float64 `json:"betaCoefficient"`
	RSquared         float64 `json:"rSquared"`
}

func (lf *LinearRegression) Predict(x float64) float64 {
	return lf.AlphaCoefficient*x + lf.BetaCoefficient
}

type PiecewiseRegression struct {
	Breakpoint                  float64 `json:"breakpoint"`
	LinearAlphaCoefficient      float64 `json:"linearAlphaCoefficient"`
	LinearBetaCoefficient       float64 `json:"linearBetaCoefficient"`
	ExponentialAlphaCoefficient float64 `json:"exponentialAlphaCoefficient"`
	ExponentialBetaCoefficient  float64 `json:"exponentialBetaCoefficient"`
	RSquared                    float64 `json:"rSquared"`
}

func (pf PiecewiseRegression) Predict(x float64) float64 {
	if x <= pf.Breakpoint {
		return pf.LinearAlphaCoefficient*x + pf.LinearBetaCoefficient
	}

	return pf.ExponentialBetaCoefficient * math.Pow(pf.ExponentialAlphaCoefficient, x)
}

type ExponentialRegression struct {
	AlphaCoefficient float64 `json:"alphaCoefficient"`
	BetaCoefficient  float64 `json:"betaCoefficient"`
	RSquared         float64 `json:"rSquared"`
}

func (ef *ExponentialRegression) Predict(x float64) float64 {
	return ef.BetaCoefficient * math.Pow(ef.AlphaCoefficient, x)
}
