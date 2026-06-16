<script lang="ts">
	import { t } from '$lib/i18n';

	interface Props {
		total: number;
		funeral: number;
		debts: number;
		bequests: number;
	}

	let { total = $bindable(), funeral = $bindable(), debts = $bindable(), bequests = $bindable() }: Props = $props();

	const fields: { key: keyof Props; label: string }[] = [
		{ key: 'total', label: 'form.total' },
		{ key: 'funeral', label: 'form.funeral' },
		{ key: 'debts', label: 'form.debts' },
		{ key: 'bequests', label: 'form.bequests' }
	];
</script>

<section class="mb-6">
	<h2 class="text-lg font-semibold mb-3">{$t('form.estate')}</h2>
	<div class="grid grid-cols-2 gap-3">
		<label>
			<span class="text-sm font-medium">{$t('form.total')}</span>
			<input
				type="number"
				min="0"
				step="any"
				bind:value={total}
				class="mt-1 block w-full rounded border border-(--color-border) px-3 py-2 text-sm"
			/>
		</label>
		<label>
			<span class="text-sm font-medium">{$t('form.funeral')}</span>
			<input
				type="number"
				min="0"
				step="any"
				bind:value={funeral}
				class="mt-1 block w-full rounded border border-(--color-border) px-3 py-2 text-sm"
			/>
		</label>
		<label>
			<span class="text-sm font-medium">{$t('form.debts')}</span>
			<input
				type="number"
				min="0"
				step="any"
				bind:value={debts}
				class="mt-1 block w-full rounded border border-(--color-border) px-3 py-2 text-sm"
			/>
		</label>
		<label>
			<span class="text-sm font-medium">{$t('form.bequests')}</span>
			<input
				type="number"
				min="0"
				step="any"
				bind:value={bequests}
				class="mt-1 block w-full rounded border border-(--color-border) px-3 py-2 text-sm"
			/>
		</label>
	</div>
</section>
