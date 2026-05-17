<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import * as echarts from 'echarts';
	import { format } from './helpers';
	import type { ComputeResponse } from '../api/models';
	import { predict, RegressionType, type RegressionModel } from '../api/regressions';

	type Props = {
		data: ComputeResponse;
	};
	let { data }: Props = $props();

	let chartNode: HTMLDivElement;
	let chart: echarts.ECharts;

	const linePoints: number[][] = $derived(
		data.correlationTableData.xMids.map((x, i) => [
			x,
			data.correlationTableData.conditionalMeanY[i]
		])
	);
	const scatterPoints: number[][] = $derived(
		data.regressionData.xPoints.map((x, i) => [x, data.regressionData.yPoints[i]])
	);

	const xMin: number = $derived(Math.min(...data.regressionData.xPoints));
	const xMax: number = $derived(Math.max(...data.regressionData.xPoints));

	const linearRegressionPoints: number[][] = $derived([
		[xMin, predict(data.regressionData.linearRegression, xMin)],
		[xMax, predict(data.regressionData.linearRegression, xMax)]
	]);

	const xPoints: number[] = $derived(
		Array.from({ length: 200 }, (_, i) => xMin + (i / 199) * (xMax - xMin))
	);
	const piecewiseRegressionPoints: number[][] = $derived(
		xPoints.map((x, _) => [x, predict(data.regressionData.piecewiseRegression, x)])
	);
	const exponentialRegressionPoints: number[][] = $derived(
		xPoints.map((x, _) => [x, predict(data.regressionData.exponentialRegression, x)])
	);

	const bestModel: RegressionModel = $derived(
		[
			data.regressionData.linearRegression,
			data.regressionData.piecewiseRegression,
			data.regressionData.exponentialRegression
		].reduce((a, b) => (a.qo < b.qo ? a : b))
	);

	let isDark: boolean = $state(window.matchMedia('(prefers-color-scheme: dark)').matches);
	onMount(() => {
		chart = echarts.init(chartNode, isDark ? 'dark' : 'light');
	});

	$effect(() => {
		if (!chart) return;
		chart.setOption({
			backgroundColor: isDark ? '#141414' : '#F5F5F5',
			legend: {
				selectedMode: 'multiple',
				selected: {
					'Empirical Linear Regression': true,
					'Regression Field': true,
					'Linear Regression': true,
					'Piecewise Regression': true,
					'Exponential Regression': false
				},
				bottom: 0,
				left: 'center'
			},
			title: {
				text: `Regression Visualization`
			},
			tooltip: {
				trigger: 'axis',
				axisPointer: {
					type: 'shadow'
				},
				formatter: (params: any) => {
					const lines = params.map((p: any) => {
						const [x, y] = p.data;
						return `${p.marker} ${p.seriesName}: (${format(x, 2)}, ${format(y, 2)})`;
					});

					return lines.join('<br/>');
				}
			},
			dataZoom: [
				{
					id: 'dataZoomX',
					type: 'slider',
					xAxisIndex: [0],
					filterMode: 'filter',
					bottom: 30
				},
				{
					id: 'dataZoomY',
					type: 'slider',
					yAxisIndex: [0],
					filterMode: 'empty'
				}
			],

			xAxis: {
				type: 'value',
				name: 'Amount of Docker Containers',
				nameLocation: 'middle',
				nameGap: 25
			},
			yAxis: {
				type: 'value',
				name: 'Startup Time (ms)',
				nameLocation: 'middle',
				nameGap: 40
			},
			grid: {
				left: '10%',
				right: '10%',
				top: 80,
				bottom: 100
			},

			series: [
				{
					type: 'line',
					data: linePoints,
					name: 'Empirical Linear Regression',
					itemStyle: {
						color: '#7C9EE8'
					},
					z: 5
				},
				{
					type: 'scatter',
					data: scatterPoints,
					name: 'Regression Field',
					symbolSize: 8,
					emphasis: {
						scale: 1.5
					},
					itemStyle: {
						color: '#A8A8A8'
					}
				},
				{
					type: 'line',
					data: linearRegressionPoints,
					name: 'Linear Regression',
					lineStyle: {
						type: 'dashed',
						width: 3
					},
					itemStyle: {
						color: '#E8845C'
					},
					symbol: 'none',
					z: 10
				},
				{
					type: 'line',
					data: piecewiseRegressionPoints,
					name: 'Piecewise Regression',
					lineStyle: {
						type: 'dotted',
						width: 3
					},
					itemStyle: {
						color: '#5BBD8A'
					},
					symbol: 'none',
					z: 10
				},
				{
					type: 'line',
					data: exponentialRegressionPoints,
					name: 'Exponential Regression',
					lineStyle: {
						type: 'dotted',
						width: 3
					},
					itemStyle: {
						color: '#5BBD8A'
					},
					symbol: 'none',
					z: 10
				}
			]
		});
	});

	onDestroy(() => chart?.dispose());
</script>

<div id="chart" bind:this={chartNode}></div>
<table>
	<thead>
		<tr>
			<th>Model</th>
			<th>Formula</th>
			<th>R²</th>
			<th>Q<sub>o</sub></th>
		</tr>
	</thead>
	<tbody>
		<tr class:best={bestModel.type === RegressionType.Linear}>
			<th>Linear</th>
			<td>
				<code
					>y = {format(data.regressionData.linearRegression.alphaCoefficient, 4)}x + {format(
						data.regressionData.linearRegression.betaCoefficient,
						4
					)}</code
				>
			</td>
			<td>{format(data.regressionData.linearRegression.rSquared, 10)}</td>
			<td>{format(data.regressionData.linearRegression.qo, 10)}</td>
		</tr>
		<tr class:best={bestModel.type === RegressionType.Piecewise}>
			<th>Piecewise</th>
			<td>
				<code
					>x ≤ {format(data.regressionData.piecewiseRegression.breakpoint, 4)}: y = {format(
						data.regressionData.piecewiseRegression.leftAlphaCoefficient,
						4
					)}x + {format(data.regressionData.piecewiseRegression.leftBetaCoefficient, 4)}</code
				>
				<br />
				<code
					>x &gt; {format(data.regressionData.piecewiseRegression.breakpoint, 4)}: y = {format(
						data.regressionData.piecewiseRegression.rightAlphaCoefficient,
						4
					)}x + {format(data.regressionData.piecewiseRegression.rightBetaCoefficient, 4)}</code
				>
			</td>
			<td>{format(data.regressionData.piecewiseRegression.rSquared, 10)}</td>
			<td>{format(data.regressionData.piecewiseRegression.qo, 10)}</td>
		</tr>
		<tr class:best={bestModel.type === RegressionType.Exponential}>
			<th>Exponential</th>
			<td>
				<code
					>y = {format(data.regressionData.exponentialRegression.betaCoefficient, 4)} · {format(
						data.regressionData.exponentialRegression.alphaCoefficient,
						4
					)}<sup>x</sup></code
				>
			</td>
			<td>{format(data.regressionData.exponentialRegression.rSquared, 10)}</td>
			<td>{format(data.regressionData.exponentialRegression.qo, 10)}</td>
		</tr>
	</tbody>
</table>

<style>
	#chart {
		width: 100%;
		height: 600px;
	}

	table {
		width: 100%;
		border-collapse: collapse;
		font-size: 1rem;
		table-layout: fixed;
		margin-top: 1rem;
	}

	th,
	td {
		padding: 0.3rem 0.6rem;
		text-align: right;
		border: 1px solid light-dark(#d0d0d0, #3a3a3a);
		vertical-align: middle;
	}

	thead th {
		background-color: light-dark(#f5f5f5, #1e1e1e);
		font-weight: 600;
	}

	thead th:nth-child(1),
	tbody th {
		width: 15%;
	} /* Model */
	thead th:nth-child(2),
	tbody td:nth-child(1) {
		width: 45%;
	} /* Formula */
	thead th:nth-child(3),
	tbody td:nth-child(2) {
		width: 20%;
	} /* R² */
	thead th:nth-child(4),
	tbody td:nth-child(3) {
		width: 20%;
	} /* Qo */

	tbody th {
		font-weight: 600;
		background-color: light-dark(#f5f5f5, #1e1e1e);
	}

	tbody td:first-child {
		text-align: left;
		white-space: normal;
		word-break: break-word;
	}

	tbody tr:hover td,
	tbody tr:hover th {
		background-color: light-dark(#efefef, #252525);
	}

	thead th:last-child,
	tbody td:last-child {
		border-left: 2px solid light-dark(#aaa, #555);
	}

	.best th,
	.best td {
		background-color: light-dark(#edfaef, #1a2e1d);
	}
</style>
