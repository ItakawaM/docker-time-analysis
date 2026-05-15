<script lang="ts">
	import type { SignificanceResponse } from '../api/models';
	import { format } from './helpers';

	type Props = {
		data: SignificanceResponse;
	};

	let { data }: Props = $props();
</script>

<div class="tables">
	{#each [{ label: 'Fisher', result: data.fisher }, { label: 'Pearson', result: data.pearson }] as { label, result }}
		<table>
			<thead>
				<tr>
					<th colspan="2">{label} Test</th>
				</tr>
			</thead>
			<tbody>
				{#if result.value !== null && result.value !== undefined}
					<tr>
						<th>Pearson Coefficient</th>
						<td>{format(result.value, 4)}</td>
					</tr>
					<tr>
						<th>Level</th>
						<td>
							{#if Math.abs(result.value) >= 0.7}
								Strong
							{:else if Math.abs(result.value) >= 0.5}
								Medium
							{:else if Math.abs(result.value) >= 0.3}
								Weak
							{:else}
								Extremely Weak
							{/if}

							{#if result.value > 0}
								Positive
							{:else if result.value < 0}
								Negative
							{/if}
						</td>
					</tr>
				{/if}
				<tr>
					<th>Empirical</th>
					<td>{format(result.empirical, 4)}</td>
				</tr>
				<tr>
					<th>Critical</th>
					<td>{format(result.critical, 4)}</td>
				</tr>
			</tbody>
			<tfoot>
				<tr>
					<th>Adequacy</th>
					<td class:adequate={result.adequate} class:inadequate={!result.adequate}>
						{result.adequate ? 'Adequate' : 'Not adequate'}
					</td>
				</tr>
			</tfoot>
		</table>
	{/each}
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

	.adequate {
		color: light-dark(#2a7a2a, #6fcf6f);
	}

	.inadequate {
		color: light-dark(#a02020, #e07070);
	}
</style>
