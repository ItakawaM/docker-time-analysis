<script lang="ts">
	import type { ComputeResponse, SignificanceResponse, UploadResponse } from './lib/api/models';
	import { postSignificance } from './lib/api/services';
	import CorrelationTable from './lib/components/CorrelationTable.svelte';
	import CSVFileInput from './lib/components/CSVFileInput.svelte';
	import PredictionInput from './lib/components/PredictionInput.svelte';
	import RegressionPlot from './lib/components/RegressionPlot.svelte';
	import SampleSizeSlider from './lib/components/SampleSizeSlider.svelte';
	import SignificanceLevelSlider from './lib/components/SignificanceLevelSlider.svelte';
	import SignificanceTable from './lib/components/SignificanceTable.svelte';

	let uploadResponse: UploadResponse | null = $state(null);
	function handleUploadSuccess(data: UploadResponse): void {
		uploadResponse = data;
	}

	function handleUploadError(): void {
		uploadResponse = null;
		computeResponse = null;
	}

	let computeResponse: ComputeResponse | null = $state(null);
	function handleComputeSuccess(data: ComputeResponse): void {
		computeResponse = data;
	}

	function handleComputeError(): void {
		computeResponse = null;
	}

	let significanceResponse: SignificanceResponse | null = $state(null);
	function handleSignificanceSuccess(data: SignificanceResponse): void {
		significanceResponse = data;
	}

	function handleSignificanceError(): void {
		significanceResponse = null;
	}
</script>

<h1 id="title">Docker Time Analyzer</h1>
<p id="subtitle">Container startup regression &amp; significance testing</p>

<section>
	<p class="instructions">
		Upload a <code>.csv</code> with <code>n_containers</code> and <code>startup_ms</code> columns — minimum
		50 rows
	</p>
	<CSVFileInput onSuccess={handleUploadSuccess} onError={handleUploadError} />
</section>

{#if uploadResponse}
	<hr />
	<p class="instructions">Choose how many rows to sample</p>
	<section>
		{#key uploadResponse}
			<SampleSizeSlider
				totalRows={uploadResponse.parsedRows}
				onSuccess={handleComputeSuccess}
				onError={handleComputeError}
			/>
		{/key}
	</section>
{/if}

{#if computeResponse}
	<p class="instructions">Startup time (ms) vs. container count — frequency distribution</p>
	<section>
		<CorrelationTable data={computeResponse} />
		<RegressionPlot data={computeResponse} />
	</section>
{/if}

{#if computeResponse}
	<p class="instructions">Predict startup time</p>
	<section>
		<PredictionInput
			linearRegressionModel={computeResponse.regressionData.linearRegression}
			piecewiseRegressionModel={computeResponse.regressionData.piecewiseRegression}
		/>
	</section>
{/if}

{#if computeResponse}
	<hr />
	<p class="instructions">Choose a significance level to validate the regression model</p>
	<section>
		<SignificanceLevelSlider
			{computeResponse}
			onSuccess={handleSignificanceSuccess}
			onError={handleSignificanceError}
		/>
	</section>
{/if}

{#if significanceResponse}
	<p class="instructions">Regression model validation results</p>
	<section>
		<SignificanceTable data={significanceResponse} />
	</section>
{/if}

<style>
	#title {
		text-align: center;
		font-size: 2rem;
		font-weight: 700;
	}

	#subtitle {
		text-align: center;
		font-size: 1.2rem;
		color: light-dark(#666, #888);
		text-transform: uppercase;
	}

	section {
		margin-block: 1rem;
	}
</style>
