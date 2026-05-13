<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import * as echarts from 'echarts';
	import type { ComputeResponse } from '../api/models';
	import { format } from './helpers';

	function plotLinearRegression(x: number, alpha: number, beta: number): number {
		return x * alpha + beta;
	}

	const COLORS = {
		light: {
			empirical: '#0D1B1E',
			scatter: '#6EA4BF',
			regression: '#AA6373',
			background: '#ffffff',
			text: '#333333'
		},
		dark: {
			empirical: '#A8DADC',
			scatter: '#6EA4BF',
			regression: '#E07A8A',
			background: '#1a1a2e',
			text: '#cccccc'
		}
	};

	type Props = {
		data: ComputeResponse;
	};
	let { data }: Props = $props();

	let chartNode: HTMLDivElement;
	let chart: echarts.ECharts;

	const linePoints: number[][] = $derived(
		data.tableData.xMids.map((x, i) => [x, data.tableData.conditionalMeanY[i]])
	);
	const scatterPoints: number[][] = $derived(
		data.graphData.xPoints.map((x, i) => [x, data.graphData.yPoints[i]])
	);

	const xMin: number = $derived(Math.min(...data.graphData.xPoints));
	const xMax: number = $derived(Math.max(...data.graphData.xPoints));

	const linearRegressionPoints: number[][] = $derived([
		[
			xMin,
			plotLinearRegression(xMin, data.graphData.alphaCoefficient, data.graphData.betaCoefficient)
		],
		[
			xMax,
			plotLinearRegression(xMax, data.graphData.alphaCoefficient, data.graphData.betaCoefficient)
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
				subtext: `y=${format(data.graphData.alphaCoefficient, 4)}x+${format(data.graphData.betaCoefficient, 4)}`
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
<p class="instructions">{`R2=${format(data.graphData.rSquared, 4)}`}</p>

<style>
	div {
		width: 100%;
		height: 600px;
	}
</style>
