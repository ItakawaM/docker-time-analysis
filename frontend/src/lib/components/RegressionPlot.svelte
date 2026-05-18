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

<style>
	#chart {
		width: 100%;
		height: 600px;
	}
</style>
