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
	Breakpoint       float64 `json:"breakpoint"`
	LinearAlpha      float64 `json:"linearAlpha"`
	LinearBeta       float64 `json:"linearBeta"`
	ExponentialAlpha float64 `json:"exponentialAlpha"`
	ExponentialBeta  float64 `json:"exponentialBeta"`
	RSquared         float64 `json:"rSquared"`
}

func (pf PiecewiseRegression) Predict(x float64) float64 {
	if x <= pf.Breakpoint {
		return pf.LinearAlpha*x + pf.LinearBeta
	}

	return pf.ExponentialBeta * math.Pow(pf.ExponentialAlpha, x)
}

type ExponentialRegression struct {
	AlphaCoefficient float64 `json:"alphaCoefficient"`
	BetaCoefficient  float64 `json:"betaCoefficient"`
	RSquared         float64 `json:"rSquared"`
}

func (ef *ExponentialRegression) Predict(x float64) float64 {
	return ef.BetaCoefficient * math.Pow(ef.AlphaCoefficient, x)
}
