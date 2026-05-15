<script lang="ts">
	import type { ComputeResponse, SignificanceResponse, UploadResponse } from './lib/api/models';
	import { postSignificance } from './lib/api/services';
	import CorrelationTable from './lib/components/CorrelationTable.svelte';
	import CSVFileInput from './lib/components/CSVFileInput.svelte';
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

<section>
	<p class="instructions">
		The .csv file must contain <code>n_containers, startup_ms</code> columns and have at least 50 entries
	</p>

	<CSVFileInput onSuccess={handleUploadSuccess} onError={handleUploadError} />
</section>

{#if uploadResponse}
	<hr />
	<p class="instructions">Select the amount of entries you would like to sample</p>
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
	<hr />
	<p class="instructions">
		Response time (ms) vs. Docker container count — joint frequency distribution
	</p>
	<section>
		<CorrelationTable data={computeResponse} />
		<RegressionPlot data={computeResponse} />
	</section>
{/if}

{#if computeResponse}
	<hr />
	<p class="instructions">Select the significance level for regression model adequecy validation</p>
	<section>
		<SignificanceLevelSlider
			{computeResponse}
			onSuccess={handleSignificanceSuccess}
			onError={handleSignificanceError}
		/>
	</section>
{/if}

{#if significanceResponse}
	<hr />
	<p class="instructions">Regression model adequecy validation</p>
	<section>
		<SignificanceTable data={significanceResponse} />
	</section>
{/if}

<style>
	#title {
		text-align: center;
		font-size: 2.2rem;
	}

	hr {
		margin-top: 1rem;
		margin-bottom: 1rem;
		border-color: light-dark(#d0d0d0, #3a3a3a);
	}
</style>
