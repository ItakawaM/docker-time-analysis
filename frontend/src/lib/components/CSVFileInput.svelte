<script lang="ts">
	import type { UploadResponse } from '../api/models';
	import { postUpload } from '../api/services';

	type Props = {
		onSuccess?: (data: UploadResponse) => void;
		onError?: () => void;
	};
	let { onSuccess, onError }: Props = $props();

	let files: FileList | null = $state(null);

	type Status = 'idle' | 'uploading' | 'success' | 'error';
	let status = $state<Status>('idle');

	let errorMessage: string = $state('');
	let successMessage: string = $state('');

	async function uploadCSV() {
		let file = files?.[0];
		if (!file) return;

		status = 'uploading';
		errorMessage = '';
		successMessage = '';

		postUpload(file)
			.then((data) => {
				successMessage = `Successfully parsed ${data.parsedRows} rows`;
				status = 'success';

				onSuccess?.(data);
			})
			.catch((err) => {
				const error = err instanceof Error ? err : new Error('Unknown error');
				errorMessage = error.message;
				status = 'error';

				onError?.();
			});
	}
</script>

<div class="input-wrapper">
	<input
		type="file"
		accept=".csv"
		disabled={status === 'uploading'}
		bind:files
		onchange={uploadCSV}
	/>
	<p class="status {status}">
		{#if status === 'idle'}
			Please upload your .csv file
		{:else if status === 'uploading'}
			Uploading your file...
		{:else if status === 'error'}
			{errorMessage}
		{:else if status === 'success'}
			{successMessage}
		{/if}
	</p>
</div>

<style>
	.input-wrapper {
		display: flex;
		flex-direction: column;
		gap: 0.4rem;
	}

	input[type='file'] {
		border-radius: 0.5rem;
		border: 1px solid light-dark(#d0d0d0, #3a3a3a);
		background-color: light-dark(#ffffff, #1e1e1e);
		color: inherit;
		padding: 0.5rem 0.5rem;
		cursor: pointer;

		&:disabled {
			opacity: 0.6;
			cursor: not-allowed;
		}
	}
</style>
