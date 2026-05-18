<script lang="ts">
	import { format } from './helpers';
	import { RegressionType, type RegressionModel } from '../api/regressions';
	import type {
		LinearRegression,
		ExponentialRegression,
		PiecewiseRegression
	} from '../api/regressions';

	type Props = {
		linearRegression: LinearRegression;
		piecewiseRegression: PiecewiseRegression;
		exponentialRegression: ExponentialRegression;
	};
	let { linearRegression, piecewiseRegression, exponentialRegression }: Props = $props();

	const bestModel: RegressionModel = $derived(
		[linearRegression, piecewiseRegression, exponentialRegression].reduce((a, b) =>
			a.qo < b.qo ? a : b
		)
	);
</script>

<table>
	<thead>
		<tr>
			<th>Model</th>
			<th>Formula</th>
			<th>R²</th>
			<th>Q<sub>o</sub></th>
			<th>Q<sub>p</sub></th>
		</tr>
	</thead>
	<tbody>
		<tr class:best={bestModel.type === RegressionType.Linear}>
			<th>Linear</th>
			<td>
				<code
					>y = {format(linearRegression.alphaCoefficient, 4)}x + {format(
						linearRegression.betaCoefficient,
						4
					)}</code
				>
			</td>
			<td>{format(linearRegression.rSquared, 6)}</td>
			<td>{format(linearRegression.qo, 2)}</td>
			<td>{format(linearRegression.qp, 2)}</td>
		</tr>
		<tr class:best={bestModel.type === RegressionType.Piecewise}>
			<th>Piecewise</th>
			<td>
				<code
					>x ≤ {format(piecewiseRegression.breakpoint, 4)}: y = {format(
						piecewiseRegression.leftAlphaCoefficient,
						4
					)}x + {format(piecewiseRegression.leftBetaCoefficient, 4)}</code
				>
				<br />
				<code
					>x &gt; {format(piecewiseRegression.breakpoint, 4)}: y = {format(
						piecewiseRegression.rightAlphaCoefficient,
						4
					)}x + {format(piecewiseRegression.rightBetaCoefficient, 4)}</code
				>
			</td>
			<td>{format(piecewiseRegression.rSquared, 6)}</td>
			<td>{format(piecewiseRegression.qo, 2)}</td>
			<td>{format(piecewiseRegression.qp, 2)}</td>
		</tr>
		<tr class:best={bestModel.type === RegressionType.Exponential}>
			<th>Exponential</th>
			<td>
				<code
					>y = {format(exponentialRegression.betaCoefficient, 4)} · {format(
						exponentialRegression.alphaCoefficient,
						4
					)}<sup>x</sup></code
				>
			</td>
			<td>{format(exponentialRegression.rSquared, 6)}</td>
			<td>{format(exponentialRegression.qo, 2)}</td>
			<td>{format(exponentialRegression.qp, 2)}</td>
		</tr>
	</tbody>
</table>

<style>
	table {
		width: 100%;
		border-collapse: collapse;
		font-size: 1rem;
		table-layout: auto;
		margin-top: 1rem;
	}

	th,
	td {
		padding: 0.3rem 0.5rem;
		text-align: right;
		border: 1px solid light-dark(#d0d0d0, #3a3a3a);
		vertical-align: middle;
		white-space: nowrap;
	}

	thead th {
		background-color: light-dark(#f5f5f5, #1e1e1e);
		font-weight: 600;
	}

	/* Formula column: allow wrapping, take remaining space */
	thead th:nth-child(2),
	tbody td:nth-child(1) {
		white-space: normal;
		word-break: break-word;
		width: 100%;
	}

	tbody th {
		font-weight: 600;
		background-color: light-dark(#f5f5f5, #1e1e1e);
	}

	tbody td:first-child {
		text-align: left;
	}

	tbody tr:hover td,
	tbody tr:hover th {
		background-color: light-dark(#efefef, #252525);
	}

	thead th:last-child,
	tbody td:last-child {
		border-left: 2px solid light-dark(#aaa, #555);
	}

	thead th:nth-child(4),
	tbody td:nth-child(3) {
		border-left: 2px solid light-dark(#aaa, #555);
	}

	.best th,
	.best td {
		background-color: light-dark(#edfaef, #1a2e1d);
	}
</style>
