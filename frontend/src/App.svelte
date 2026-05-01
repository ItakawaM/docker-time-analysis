<script lang="ts">
	import type { UploadResponse } from './lib/api/models';
	import CSVFileInput from './lib/components/CSVFileInput.svelte';

	let uploadResponse: UploadResponse | null = $state(null);
	function handleSuccess(data: UploadResponse): void {
		uploadResponse = data;
	}

	function handleError(): void {
		uploadResponse = null;
	}
</script>

<section>
	<h1 id="title">Docker Time Analyzer</h1>

	<p class="instructions">
		The .csv file must contain <code>n_containers, startup_ms</code> columns and have at least 50 entries
	</p>

	<CSVFileInput onSuccess={handleSuccess} onError={handleError} />

	{#if uploadResponse}
		<p>Here comes the slider...</p>
	{/if}
</section>

<style>
	#title {
		text-align: center;
		font-size: 2.2rem;
	}

	.instructions {
		text-align: left;
		color: light-dark(#555, #aaa);
		margin-inline: auto;
	}
</style>
