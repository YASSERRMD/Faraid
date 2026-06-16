<script lang="ts">
	import { t } from '$lib/i18n';
	import { faraidClient } from '$lib/client';
	import type { components } from '$lib/client';
	import DeceasedSection from '$lib/components/DeceasedSection.svelte';
	import EstateSection from '$lib/components/EstateSection.svelte';
	import HeirsSection from '$lib/components/HeirsSection.svelte';
	import DerivationPanel from '$lib/components/DerivationPanel.svelte';
	import ExplainButton from '$lib/components/ExplainButton.svelte';
	import NLEntry from '$lib/components/NLEntry.svelte';
	import { validate, type Counts } from '$lib/heirs';

	type SolveResult = components['schemas']['SolveResult'];

	type Madhhab = 'Hanafi' | 'Maliki' | 'Shafii' | 'Hanbali';
	let sex = $state<'male' | 'female'>('male');
	let madhhab = $state<Madhhab>('Hanafi');
	let total = $state(0);
	let funeral = $state(0);
	let debts = $state(0);
	let bequests = $state(0);
	let counts = $state<Counts>({});

	let result = $state<SolveResult | null>(null);
	let submitError = $state<string | null>(null);
	let loading = $state(false);
	let lastRequest = $state<components['schemas']['SolveRequest'] | null>(null);

	function handleParsed(parsedSex: 'male' | 'female', parsedCounts: Counts) {
		sex = parsedSex;
		counts = parsedCounts;
	}

	// Live validation: recomputed on every state change.
	const errors = $derived(validate(counts, sex));
	const canSubmit = $derived(errors.length === 0);

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!canSubmit) return;
		loading = true;
		submitError = null;
		result = null;

		// Build heirs object: only include non-zero counts.
		const heirs: Record<string, number> = {};
		for (const [k, v] of Object.entries(counts)) {
			if (v > 0) heirs[k] = v;
		}

		const { data, error } = await faraidClient.POST('/solve', {
			body: {
				deceasedSex: sex,
				heirs,
				madhhab,
				estate: total > 0 ? {
					total,
					funeral,
					debts,
					bequests,
					heirsConsentToExcessBequest: false
				} : undefined
			}
		});
		loading = false;
		if (error) {
			submitError = typeof error === 'object' && 'error' in error
				? String((error as { error: unknown }).error)
				: JSON.stringify(error);
		} else {
			result = data ?? null;
		if (result) lastRequest = { deceasedSex: sex, heirs, madhhab };
		}
	}
</script>

<svelte:head>
	<title>{$t('nav.calculate')} - {$t('app.title')}</title>
</svelte:head>

<div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
	<div>
		<h1 class="text-2xl font-bold mb-6">{$t('nav.calculate')}</h1>

		<NLEntry onParsed={handleParsed} />

		<form onsubmit={handleSubmit}>
			<DeceasedSection bind:sex bind:madhhab />
			<EstateSection bind:total bind:funeral bind:debts bind:bequests />
			<HeirsSection bind:counts />

			{#if errors.length > 0}
				<div class="mb-4 rounded border border-red-300 bg-red-50 px-3 py-2">
					{#each errors as err}
						<p class="text-sm text-red-700">{$t(err)}</p>
					{/each}
				</div>
			{/if}

			<button
				type="submit"
				disabled={!canSubmit || loading}
				class="w-full rounded bg-(--color-primary) px-4 py-3 font-semibold text-white disabled:opacity-50 hover:bg-(--color-primary-hover) transition-colors"
			>
				{loading ? '...' : $t('form.calculate')}
			</button>
		</form>
	</div>

	<div>
		{#if submitError}
			<div class="rounded border border-red-300 bg-red-50 px-3 py-2 text-sm text-red-700">
				{submitError}
			</div>
		{/if}

		{#if result && lastRequest}
			<ExplainButton request={lastRequest} />
		{/if}

		{#if result}
			<h2 class="text-xl font-bold mb-4">{$t('result.shares')}</h2>

			{#if result.specialCase}
				<p class="mb-3 text-sm italic text-gray-500">
					{$t('result.special_case')}: {result.specialCase}
				</p>
			{/if}
			{#if result.awl}
				<p class="mb-2 text-xs text-amber-600">{$t('result.awl')}</p>
			{/if}
			{#if result.radd}
				<p class="mb-2 text-xs text-amber-600">{$t('result.radd')}</p>
			{/if}

			<table class="w-full text-sm border-collapse">
				<thead>
					<tr class="border-b border-(--color-border)">
						<th class="text-start py-2 font-medium">{$t('result.heir')}</th>
						<th class="text-end py-2 font-medium">{$t('result.count')}</th>
						<th class="text-end py-2 font-medium">{$t('result.fraction')}</th>
						{#if (result.shares?.[0]?.amount ?? '0') !== '0'}
							<th class="text-end py-2 font-medium">{$t('result.amount')}</th>
						{/if}
					</tr>
				</thead>
				<tbody>
					{#each result.shares ?? [] as share}
						<tr class="border-b border-(--color-border) last:border-0">
							<td class="py-2">{$t(`heir.${share.relation}`)}</td>
							<td class="text-end py-2 tabular-nums">{share.count}</td>
							<td class="text-end py-2 tabular-nums font-mono">{share.fraction}</td>
							{#if (result.shares?.[0]?.amount ?? '0') !== '0'}
								<td class="text-end py-2 tabular-nums">{share.amount}</td>
							{/if}
						</tr>
					{/each}
				</tbody>
			</table>

			{#if result.needsReview}
				<p class="mt-4 text-sm text-amber-600">{$t('result.needs_review')}</p>
			{/if}

			{#if result.derivation && result.derivation.length > 0}
				<DerivationPanel steps={result.derivation} />
			{/if}
		{/if}
	</div>
</div>
