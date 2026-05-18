package compute

import "math"

// Regression is an interface for different regression models that can predict values based on input x.
type Regression interface {
	Predict(x float64) float64
}

// LinearRegression represents a linear regression model of the form y = alpha*x + beta.
// It contains coefficients, R-squared goodness-of-fit, and residual sum of squares (Qo), and regression sum of squares (Qp).
type LinearRegression struct {
	Type             string  `json:"type"`
	AlphaCoefficient float64 `json:"alphaCoefficient"`
	BetaCoefficient  float64 `json:"betaCoefficient"`
	RSquared         float64 `json:"rSquared"`
	Qo               float64 `json:"qo"`
	Qp               float64 `json:"qp"`
}

// Predict calculates the predicted value for input x using the linear regression equation.
func (lf *LinearRegression) Predict(x float64) float64 {
	return lf.AlphaCoefficient*x + lf.BetaCoefficient
}

// PiecewiseRegression represents a piecewise linear regression model with two segments.
// It defines left and right linear equations separated by a breakpoint, with R-squared, residual sum of squares (Qo),
// and regression sum of squares (Qp).
type PiecewiseRegression struct {
	Type                  string  `json:"type"`
	Breakpoint            float64 `json:"breakpoint"`
	LeftAlphaCoefficient  float64 `json:"leftAlphaCoefficient"`
	LeftBetaCoefficient   float64 `json:"leftBetaCoefficient"`
	RightAlphaCoefficient float64 `json:"rightAlphaCoefficient"`
	RightBetaCoefficient  float64 `json:"rightBetaCoefficient"`
	RSquared              float64 `json:"rSquared"`
	Qo                    float64 `json:"qo"`
	Qp                    float64 `json:"qp"`
}

// Predict calculates the predicted value for input x using the appropriate piecewise regression equation based on the breakpoint.
func (pf PiecewiseRegression) Predict(x float64) float64 {
	if x <= pf.Breakpoint {
		return pf.LeftAlphaCoefficient*x + pf.LeftBetaCoefficient
	}

	return pf.RightAlphaCoefficient*x + pf.RightBetaCoefficient
}

// ExponentialRegression represents an exponential regression model of the form y = beta * (alpha ^ x).
// It contains coefficients, R-squared goodness-of-fit, residual sum of squares (Qo), and regression sum of squares (Qp).
type ExponentialRegression struct {
	Type             string  `json:"type"`
	AlphaCoefficient float64 `json:"alphaCoefficient"`
	BetaCoefficient  float64 `json:"betaCoefficient"`
	RSquared         float64 `json:"rSquared"`
	Qo               float64 `json:"qo"`
	Qp               float64 `json:"qp"`
}

// Predict calculates the predicted value for input x using the exponential regression equation.
func (ef *ExponentialRegression) Predict(x float64) float64 {
	return ef.BetaCoefficient * math.Pow(ef.AlphaCoefficient, x)
}
