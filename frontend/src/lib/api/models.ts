export type UploadResponse = {
	message: string;
	parsedRows: number;
};

export type ComputeRequest = {
	sampleSize: number;
};

export type TableData = {
	yMids: number[];
	xMids: number[];

	frequencies: number[];

	xMarginal: number[];
	yMarginal: number[];

	conditionalMeanY: number[];
};

export type GraphData = {
	yPoints: number[];
	xPoints: number[];

	alphaCoefficient: number;
	betaCoefficient: number;
	rSquared: number;
};

export type ComputeResponse = {
	tableData: TableData;
	graphData: GraphData;
};
