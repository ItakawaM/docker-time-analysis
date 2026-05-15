<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import * as echarts from 'echarts';
	import type { ComputeResponse } from '../api/models';
	import { format } from './helpers';

	function plotLinearRegression(x: number, alpha: number, beta: number): number {
		return x * alpha + beta;
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
		[
			xMin,
			plotLinearRegression(
				xMin,
				data.regressionData.alphaCoefficient,
				data.regressionData.betaCoefficient
			)
		],
		[
			xMax,
			plotLinearRegression(
				xMax,
				data.regressionData.alphaCoefficient,
				data.regressionData.betaCoefficient
			)
		]
	]);

	let isDark: boolean = $state(window.matchMedia('(prefers-color-scheme: dark)').matches);
	onMount(() => {
		chart = echarts.init(chartNode, isDark ? 'dark' : 'light');
	});

	$effect(() => {
		if (!chart) return;
		chart.setOption({
			backgroundColor: isDark ? '#141414' : 'ffffff',
			legend: {
				selectedMode: 'multiple',
				bottom: 0,
				left: 'center'
			},
			title: {
				text: `Regression Visualization`,
				subtext: `y=${format(data.regressionData.alphaCoefficient, 4)}x+${format(data.regressionData.betaCoefficient, 4)} R2=${format(data.regressionData.rSquared, 4)}`
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
					name: 'Empirical Linear Regression'
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
						color: '#6EA4BF'
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
						color: '#AA6373'
					},
					symbol: 'none',
					z: 10
				}
			]
		});
	});

	onDestroy(() => chart?.dispose());
</script>

<p></p>
<div bind:this={chartNode}></div>

<style>
	div {
		width: min(85ch, 100% - 4rem);
		height: 600px;
	}
</style>
