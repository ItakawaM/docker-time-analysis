export enum RegressionType {
	Linear = 'Linear',
	Piecewise = 'Piecewise',
	Exponential = 'Exponential'
}

export type RegressionModel = LinearRegression | ExponentialRegression | PiecewiseRegression;

export type LinearRegression = {
	type: RegressionType.Linear;
	alphaCoefficient: number;
	betaCoefficient: number;
	rSquared: number;
	qo: number;
};

export type ExponentialRegression = {
	type: RegressionType.Exponential;
	alphaCoefficient: number;
	betaCoefficient: number;
	rSquared: number;
	qo: number;
};

export type PiecewiseRegression = {
	type: RegressionType.Piecewise;
	breakpoint: number;
	leftAlphaCoefficient: number;
	leftBetaCoefficient: number;
	rightAlphaCoefficient: number;
	rightBetaCoefficient: number;
	rSquared: number;
	qo: number;
};

export function predict(model: RegressionModel, x: number): number {
	switch (model.type) {
		case RegressionType.Linear:
			return model.alphaCoefficient * x + model.betaCoefficient;

		case RegressionType.Exponential:
			return model.betaCoefficient * Math.pow(model.alphaCoefficient, x);

		case RegressionType.Piecewise:
			if (x <= model.breakpoint) {
				return model.leftAlphaCoefficient * x + model.leftBetaCoefficient;
			}

			return model.rightAlphaCoefficient * x + model.rightBetaCoefficient;
	}
}
