export type UploadResponse = {
	message: string;
	parsedRows: number;
};

export type ComputeRequest = {
	sampleSize: number;
};

export type ComputeResponse = {
	correlationTableData: CorrelationTableData;
	regressionData: RegressionData;
};

export type CorrelationTableData = {
	yMids: number[];
	xMids: number[];

	frequencies: number[];

	xMarginal: number[];
	yMarginal: number[];

	conditionalMeanY: number[];
};

export type RegressionData = {
	yPoints: number[];
	xPoints: number[];

	alphaCoefficient: number;
	betaCoefficient: number;

	rSquared: number;
};

export type SignificanceRequest = {
	significanceLevel: number;
};

export type SignificanceResponse = {
	fisher: StatTestResult;
	pearson: StatTestResult;
};

export type StatTestResult = {
	value: number | null;
	empirical: number;
	critical: number;
	adequate: boolean;
};
