<script lang="ts">
	import type { ComputeResponse, SignificanceResponse } from '../api/models';
	import { postSignificance } from '../api/services';
	import { clamp, debounce } from './helpers';

	type Props = {
		computeResponse: ComputeResponse | null;
		onSuccess?: (data: SignificanceResponse) => void;
		onError?: () => void;
	};
	let { computeResponse, onSuccess, onError }: Props = $props();

	type Status = 'idle' | 'success' | 'error';
	let status = $state<Status>('idle');

	let errorMessage: string = $state('');
	let successMessage: string = $state('');

	const min = 0.001;
	const max = 0.99;
	const step = 0.001;

	let value: number = $state(0.05);
	let debounced: number = $state(0.05);

	const setDebounced = debounce((value: number) => {
		debounced = value;
	}, 375);

	$effect(() => {
		value = clamp(min, max, value);
	});

	$effect(() => {
		setDebounced(value);
	});

	$effect(() => {
		const controller = new AbortController();

		if (computeResponse) {
			postSignificance({ significanceLevel: debounced }, controller.signal)
				.then((data) => {
					successMessage = `Validated model adequecy at significance level ${debounced}`;
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
		}

		return () => controller.abort();
	});
</script>

<div class="slider-wrapper">
	<div class="slider-block">
		<input type="range" {min} {max} {step} bind:value />
		<input type="number" {min} {max} {step} bind:value />
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
		padding: 0.5rem;
	}

	.slider-block input[type='range'] {
		flex: 1;
	}
</style>
