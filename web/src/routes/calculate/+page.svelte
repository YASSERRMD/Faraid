<script lang="ts">
	import { t } from '$lib/i18n';
	import { faraidClient } from '$lib/client';
	import type { components } from '$lib/client';

	type SolveResult = components['schemas']['SolveResult'];

	let result = $state<SolveResult | null>(null);
	let error = $state<string | null>(null);
	let loading = $state(false);

	// Minimal smoke call: will be replaced by the full heir-builder form in Phase 43.
	async function smokeTest() {
		loading = true;
		error = null;
		const { data, error: err } = await faraidClient.POST('/solve', {
			body: {
				deceasedSex: 'male',
				heirs: { wife: 1, son: 1 },
				madhhab: 'Hanafi'
			}
		});
		loading = false;
		if (err) {
			error = JSON.stringify(err);
		} else {
			result = data ?? null;
		}
	}
</script>

<svelte:head>
	<title>{$t('nav.calculate')} - {$t('app.title')}</title>
</svelte:head>

<h1 class="text-2xl font-bold mb-6">{$t('nav.calculate')}</h1>

<p class="text-sm text-gray-500 mb-4">
	Full heir builder coming in the next phase. Use the smoke test below to verify the API client.
</p>

<button
	onclick={smokeTest}
	disabled={loading}
	class="rounded bg-(--color-primary) px-4 py-2 text-white text-sm disabled:opacity-50"
>
	{loading ? 'Loading...' : 'API Smoke Test (wife + son, Hanafi)'}
</button>

{#if error}
	<p class="mt-4 text-red-600 text-sm">{error}</p>
{/if}

{#if result}
	<div class="mt-6 rounded border border-(--color-border) p-4 font-mono text-sm">
		<p class="font-semibold mb-2">{$t('result.shares')}</p>
		{#each result.shares ?? [] as share}
			<div class="flex gap-4">
				<span class="w-40">{share.relation}</span>
				<span>{share.fraction}</span>
			</div>
		{/each}
	</div>
{/if}
