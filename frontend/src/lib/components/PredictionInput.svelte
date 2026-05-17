<script lang="ts">
	import {
		predict,
		type ExponentialRegression,
		type LinearRegression,
		type PiecewiseRegression,
		type RegressionModel
	} from '../api/regressions';
	import { debounce, format } from './helpers';

	type Props = {
		linearRegression: LinearRegression;
		piecewiseRegression: PiecewiseRegression;
		exponentialRegression: ExponentialRegression;
	};
	let { linearRegression, piecewiseRegression, exponentialRegression }: Props = $props();

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

	const bestModel: RegressionModel = $derived(
		[linearRegression, piecewiseRegression, exponentialRegression].reduce((a, b) =>
			a.rSquared > b.rSquared ? a : b
		)
	);

	let predicted: number = $derived(predict(bestModel, debounced));
</script>

<div class="wrapper">
	<input type="number" {min} bind:value />
	<p class="status">
		{format(predicted, 4)}(ms) &nbsp;·&nbsp; {bestModel.type} Regression model
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
