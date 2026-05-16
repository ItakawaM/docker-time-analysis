<script lang="ts">
	import type { ComputeResponse } from '../api/models';
	import { postCompute } from '../api/services';
	import { clamp, debounce } from './helpers';

	type Props = {
		totalRows: number;
		onSuccess?: (data: ComputeResponse) => void;
		onError?: () => void;
	};
	let { totalRows, onSuccess, onError }: Props = $props();

	type Status = 'idle' | 'success' | 'error';
	let status = $state<Status>('idle');

	let errorMessage: string = $state('');
	let successMessage: string = $state('');

	const min = 50;
	let value: number = $state(min);

	let debounced: number = $state(min);
	const setDebounced = debounce((value: number) => {
		debounced = value;
	}, 375);

	$effect(() => {
		value = clamp(min, totalRows, value);
	});

	$effect(() => {
		setDebounced(value);
	});

	$effect(() => {
		const controller = new AbortController();

		postCompute({ sampleSize: debounced }, controller.signal)
			.then((data) => {
				successMessage = `${debounced} rows sampled from ${totalRows} uploaded`;
				status = 'success';

				onSuccess?.(data);
			})
			.catch((err) => {
				if (err.name === 'AbortError') return;

				const error = err instanceof Error ? err : new Error('Unknown error');
				errorMessage = error.message;
				status = 'error';

				onError?.();
			});

		return () => controller.abort();
	});
</script>

<div class="slider-wrapper">
	<div class="slider-block">
		<input type="range" {min} step="1" bind:value max={totalRows} />
		<input type="number" {min} step="1" bind:value max={totalRows} />
	</div>
	<p class="status {status}">
		{#if status === 'error'}
			{errorMessage}
		{:else if status == 'success'}
			{successMessage}
		{/if}
	</p>
</div>

<style>
	.slider-wrapper {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	.slider-block {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		border-radius: 0.5rem;
		border: 1px solid light-dark(#d0d0d0, #3a3a3a);
		background-color: light-dark(#f5f5f5, #1e1e1e);
		padding: 0.5rem 0.5rem;
	}

	.slider-block input[type='range'] {
		flex: 1;
	}
</style>
