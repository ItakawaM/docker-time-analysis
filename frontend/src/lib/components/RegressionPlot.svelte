<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import * as echarts from 'echarts';
	import type { ComputeResponse, LinearRegression, PiecewiseRegression } from '../api/models';
	import { format } from './helpers';

	function plotLinearRegression(x: number, linearRegression: LinearRegression): number {
		return x * linearRegression.alphaCoefficient + linearRegression.betaCoefficient;
	}

	function plotPiecewiseRegression(x: number, piecewiseRegression: PiecewiseRegression): number {
		if (x <= piecewiseRegression.breakpoint) {
			return (
				x * piecewiseRegression.linearAlphaCoefficient + piecewiseRegression.linearBetaCoefficient
			);
		}

		return (
			piecewiseRegression.exponentialBetaCoefficient *
			Math.pow(piecewiseRegression.exponentialAlphaCoefficient, x)
		);
	}

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
		[xMin, plotLinearRegression(xMin, data.regressionData.linearRegression)],
		[xMax, plotLinearRegression(xMax, data.regressionData.linearRegression)]
	]);

	const xPoints: number[] = $derived(
		Array.from({ length: 200 }, (_, i) => xMin + (i / 199) * (xMax - xMin))
	);
	const piecewiseRegressionPoints: number[][] = $derived(
		xPoints.map((x, _) => [x, plotPiecewiseRegression(x, data.regressionData.piecewiseRegression)])
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
					}
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
				}
			]
		});
	});

	onDestroy(() => chart?.dispose());
</script>

<div id="chart" bind:this={chartNode}></div>
<div class="models">
	<div class="model">
		<p class="model-label">Linear Regression Model</p>
		<p>
			<code
				>y = {format(data.regressionData.linearRegression.alphaCoefficient, 4)}x + {format(
					data.regressionData.linearRegression.betaCoefficient,
					4
				)}</code
			>
			&nbsp;·&nbsp; R² = {format(data.regressionData.linearRegression.rSquared, 10)}
			&nbsp;·&nbsp; Q<sub>o</sub> = {format(data.regressionData.linearRegression.qo, 10)}
		</p>
	</div>
	<div class="model">
		<p class="model-label">Piecewise Regression Model</p>
		<p>
			x ≤ {format(data.regressionData.piecewiseRegression.breakpoint, 4)}:
			<code
				>y = {format(data.regressionData.piecewiseRegression.linearAlphaCoefficient, 4)}x + {format(
					data.regressionData.piecewiseRegression.linearBetaCoefficient,
					4
				)}</code
			>
		</p>
		<p>
			x &gt; {format(data.regressionData.piecewiseRegression.breakpoint, 4)}:
			<code
				>y = {format(data.regressionData.piecewiseRegression.exponentialBetaCoefficient, 4)} · {format(
					data.regressionData.piecewiseRegression.exponentialAlphaCoefficient,
					4
				)}<sup>x</sup></code
			>
			&nbsp;·&nbsp; R² = {format(data.regressionData.piecewiseRegression.rSquared, 10)}
			&nbsp;·&nbsp; Q<sub>o</sub> = {format(data.regressionData.piecewiseRegression.qo, 10)}
		</p>
	</div>
	<p class="status success">
		{#if data.regressionData.linearRegression.qo > data.regressionData.piecewiseRegression.qo}
			Piecewise Regression Model
		{:else}
			Linear Regression Model
		{/if} is better based on residual variance
	</p>
</div>

<style>
	#chart {
		width: 100%;
		height: 600px;
	}

	.models {
		display: flex;
		flex-direction: column;
		gap: 1rem;
		font-size: 1rem;
		color: light-dark(#3a3a3a, #aaa);
	}

	.model {
		margin-top: 0.75rem;
	}

	.model-label {
		font-size: 0.85rem;
		text-transform: uppercase;
		letter-spacing: 0.05em;
		color: light-dark(#3a3a3a, #aaa);
		margin-bottom: 0.15rem;
	}
</style>
