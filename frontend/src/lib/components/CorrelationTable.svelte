<script lang="ts">
	import type { ComputeResponse } from '../api/models';
	import { format } from './helpers';

	type Props = {
		data: ComputeResponse;
	};
	let { data }: Props = $props();

	const rows: number = $derived(data.tableData.yMarginal.length);

	function freq(row: number, column: number): number {
		return data.tableData.frequencies[row * rows + column];
	}
</script>

<div class="tables">
	<table>
		<thead>
			<tr>
				<th>Y \ X</th>
				{#each data.tableData.xMids as x}
					<th>{format(x)}</th>
				{/each}
				<th>Σ Y</th>
			</tr>
		</thead>
		<tbody>
			{#each data.tableData.yMids as y, row}
				<tr>
					<th>{format(y)}</th>
					{#each data.tableData.xMids as _, col}
						<td>{format(freq(row, col))}</td>
					{/each}
					<td>{format(data.tableData.yMarginal[row])}</td>
				</tr>
			{/each}
		</tbody>
		<tfoot>
			<tr>
				<th>Σ X</th>
				{#each data.tableData.xMarginal as xm}
					<td>{format(xm)}</td>
				{/each}
				<td
					>{format(
						data.tableData.xMarginal.reduce((accumulator, currentValue) => {
							return accumulator + currentValue;
						}, 0)
					)}</td
				>
			</tr>
		</tfoot>
	</table>

	<table>
		<tbody>
			<tr>
				<th>X</th>
				{#each data.tableData.xMids as x}
					<td>{format(x)}</td>
				{/each}
			</tr>
			<tr>
				<th>Ŷ</th>
				{#each data.tableData.conditionalMeanY as y}
					<td>{format(y)}</td>
				{/each}
			</tr>
			<tr>
				<th>N</th>
				{#each data.tableData.xMarginal as n}
					<td>{format(n)}</td>
				{/each}
			</tr>
		</tbody>
	</table>
</div>

<style>
	table {
		width: 100%;
		border-collapse: collapse;
		font-size: 1rem;
	}

	th,
	td {
		padding: 0.3rem 0.6rem;
		text-align: right;
		border: 1px solid light-dark(#d0d0d0, #3a3a3a);
	}

	thead th {
		background-color: light-dark(#f5f5f5, #1e1e1e);
		font-weight: 600;
	}

	tbody th {
		font-weight: 600;
		background-color: light-dark(#f5f5f5, #1e1e1e);
	}

	tfoot th,
	tfoot td {
		background-color: light-dark(#f5f5f5, #1e1e1e);
		font-weight: 600;
	}

	tbody tr:hover td,
	tbody tr:hover th {
		background-color: light-dark(#efefef, #252525);
	}

	thead th:last-child,
	tbody td:last-child,
	tfoot td:last-child {
		border-left: 2px solid light-dark(#aaa, #555);
	}

	.tables {
		display: flex;
		flex-direction: column;
		gap: 1.5rem;
	}
</style>
