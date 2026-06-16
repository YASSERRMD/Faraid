<script lang="ts">
	import { t } from '$lib/i18n';
	import { faraidClient } from '$lib/client';
	import type { components } from '$lib/client';
	import DeceasedSection from '$lib/components/DeceasedSection.svelte';
	import EstateSection from '$lib/components/EstateSection.svelte';
	import HeirsSection from '$lib/components/HeirsSection.svelte';
	import { validate, type Counts } from '$lib/heirs';

	type CompareResult = components['schemas']['Comparison'];
	type SolveResult = components['schemas']['SolveResult'];
	type Madhhab = 'Hanafi' | 'Maliki' | 'Shafii' | 'Hanbali';

	let sex = $state<'male' | 'female'>('male');
	let madhhab = $state<Madhhab>('Hanafi');
	let total = $state(0);
	let funeral = $state(0);
	let debts = $state(0);
	let bequests = $state(0);
	let counts = $state<Counts>({});

	let result = $state<CompareResult | null>(null);
	let submitError = $state<string | null>(null);
	let loading = $state(false);

	const errors = $derived(validate(counts, sex));
	const canSubmit = $derived(errors.length === 0);

	const madhahib: Madhhab[] = ['Hanafi', 'Maliki', 'Shafii', 'Hanbali'];

	async function handleSubmit(e: SubmitEvent) {
		e.preventDefault();
		if (!canSubmit) return;
		loading = true;
		submitError = null;
		result = null;

		const heirs: Record<string, number> = {};
		for (const [k, v] of Object.entries(counts)) {
			if (v > 0) heirs[k] = v;
		}

		const { data, error } = await faraidClient.POST('/compare', {
			body: {
				deceasedSex: sex,
				heirs,
				estate: total > 0 ? { total, funeral, debts, bequests, heirsConsentToExcessBequest: false } : undefined
			}
		});
		loading = false;
		if (error) {
			submitError = JSON.stringify(error);
		} else {
			result = data ?? null;
		}
	}

	// Collect all heir relations that appear in any school's result.
	const allHeirs = $derived((): string[] => {
		if (!result) return [];
		const seen = new Set<string>();
		for (const m of madhahib) {
			const r = result.results?.[m] as SolveResult | undefined;
			for (const s of r?.shares ?? []) seen.add(s.relation);
		}
		return [...seen];
	});

	// Return the fraction for a relation in a school, or '-' when absent.
	function fractionFor(school: Madhhab, relation: string): string {
		const r = result?.results?.[school] as SolveResult | undefined;
		return r?.shares?.find((s) => s.relation === relation)?.fraction ?? '-';
	}

	// A relation diverges when at least two schools disagree on its fraction.
	function diverges(relation: string): boolean {
		const fractions = madhahib.map((m) => fractionFor(m, relation));
		return new Set(fractions).size > 1;
	}
</script>

<svelte:head>
	<title>{$t('nav.compare')} - {$t('app.title')}</title>
</svelte:head>

<h1 class="text-2xl font-bold mb-6">{$t('nav.compare')}</h1>

<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
	<div class="lg:col-span-1">
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
				{loading ? '...' : $t('nav.compare')}
			</button>
		</form>
	</div>

	<div class="lg:col-span-2">
		{#if submitError}
			<p class="text-sm text-red-600">{submitError}</p>
		{/if}

		{#if result}
			{#if result.divergences && result.divergences.length > 0}
				<div class="mb-4 rounded border border-amber-300 bg-amber-50 px-3 py-2">
					<p class="text-xs font-semibold text-amber-700 mb-1">Divergences</p>
					{#each result.divergences as d}
						<p class="text-xs text-amber-700">{d}</p>
					{/each}
				</div>
			{/if}

			<div class="overflow-x-auto">
				<table class="w-full text-sm border-collapse">
					<thead>
						<tr class="border-b border-(--color-border)">
							<th class="text-start py-2 pe-4 font-medium">{$t('result.heir')}</th>
							{#each madhahib as m}
								<th class="text-end py-2 px-2 font-medium">{$t(`madhhab.${m}`)}</th>
							{/each}
						</tr>
					</thead>
					<tbody>
						{#each allHeirs() as relation}
							<tr
								class="border-b border-(--color-border) last:border-0"
								class:bg-amber-50={diverges(relation)}
							>
								<td class="py-2 pe-4">{$t(`heir.${relation}`)}</td>
								{#each madhahib as m}
									<td class="text-end py-2 px-2 tabular-nums font-mono">
										{fractionFor(m, relation)}
									</td>
								{/each}
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
			<p class="mt-2 text-xs text-gray-400">Rows highlighted in amber differ between schools.</p>
		{/if}
	</div>
</div>
