<script lang="ts">
	import { linear } from 'svelte/easing';
	import type { LinearRegression, PiecewiseRegression } from '../api/models';
	import { debounce, format } from './helpers';

	type Props = {
		linearRegressionModel: LinearRegression;
		piecewiseRegressionModel: PiecewiseRegression;
	};
	let { linearRegressionModel, piecewiseRegressionModel }: Props = $props();

	function predict(value: number): number {
		if (linearRegressionModel.rSquared > piecewiseRegressionModel.rSquared) {
			return linearRegressionModel.alphaCoefficient * value + linearRegressionModel.betaCoefficient;
		}

		if (value < piecewiseRegressionModel.breakpoint) {
			return (
				piecewiseRegressionModel.linearAlphaCoefficient * value +
				piecewiseRegressionModel.linearBetaCoefficient
			);
		}

		return (
			piecewiseRegressionModel.exponentialBetaCoefficient *
			Math.pow(piecewiseRegressionModel.exponentialAlphaCoefficient, value)
		);
	}

	const min: number = 0;

	let value: number = $state(0);
	let debounced: number = $state(min);
	const setDebounced: Function = debounce((value: number) => {
		debounced = value;
	}, 375);

	$effect(() => {
		value = Math.max(0, value);
	});

	$effect(() => {
		setDebounced(value);
	});

	let predicted: number = $derived(predict(debounced));
</script>

<div class="wrapper">
	<input type="number" {min} bind:value />
	<p class="status">
		{format(predicted, 4)}(ms) &nbsp;·&nbsp;
		{#if linearRegressionModel.rSquared > piecewiseRegressionModel.rSquared}
			Linear
		{:else}
			Piecewise
		{/if} Regression model
	</p>
</div>

<style>
	.wrapper {
		display: flex;
		align-items: center;
		gap: 0.5rem;
		border-radius: 0.5rem;
		border: 1px solid light-dark(#d0d0d0, #3a3a3a);
		background-color: light-dark(#f5f5f5, #1e1e1e);
		padding: 0.5rem;
	}

	input[type='number'] {
		flex: 1;
	}

	p {
		flex: 1;
	}
</style>
