<script lang="ts">
	import { t } from '$lib/i18n';
	import { llmEnabled } from '$lib/flags';
	import { faraidClient } from '$lib/client';
	import type { components } from '$lib/client';

	interface Props {
		request: components['schemas']['SolveRequest'];
	}
	let { request }: Props = $props();

	let confirmed = $state(false);
	let explanation = $state<string | null>(null);
	let inconsistent = $state(false);
	let loading = $state(false);
	let error = $state<string | null>(null);

	async function fetchExplanation() {
		loading = true;
		error = null;
		const { data, error: err } = await faraidClient.POST('/explain', { body: request });
		loading = false;
		if (err) {
			error = JSON.stringify(err);
		} else if (data) {
			explanation = data.text ?? null;
			inconsistent = !(data.consistent ?? true);
		}
	}
</script>

{#if llmEnabled}
	{#if !confirmed}
		<div class="mt-6 rounded border border-amber-300 bg-amber-50 p-4">
			<p class="text-sm font-semibold text-amber-800 mb-2">{$t('llm.explain_warning')}</p>
			<button
				type="button"
				onclick={() => { confirmed = true; fetchExplanation(); }}
				class="rounded bg-amber-600 px-3 py-1.5 text-sm text-white hover:bg-amber-700"
			>
				{$t('llm.explain_button')}
			</button>
		</div>
	{:else if loading}
		<p class="mt-4 text-sm text-gray-500">...</p>
	{:else if error}
		<p class="mt-4 text-sm text-red-600">{error}</p>
	{:else if explanation}
		<div class="mt-6 rounded border border-(--color-border) p-4">
			<div class="flex items-center gap-2 mb-2">
				<span class="rounded bg-amber-100 px-2 py-0.5 text-xs font-medium text-amber-700">
					{$t('llm.badge')}
				</span>
			</div>
			{#if inconsistent}
				<p class="text-xs text-amber-600 mb-2">{$t('llm.inconsistent')}</p>
			{/if}
			<p class="text-sm whitespace-pre-wrap">{explanation}</p>
		</div>
	{/if}
{/if}
