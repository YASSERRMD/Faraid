<script lang="ts">
	import { t } from '$lib/i18n';
	import { llmEnabled } from '$lib/flags';
	import { faraidClient } from '$lib/client';
	import type { components } from '$lib/client';
	import type { Counts } from '$lib/heirs';

	interface Props {
		onParsed: (sex: 'male' | 'female', heirs: Counts) => void;
	}
	let { onParsed }: Props = $props();

	let text = $state('');
	let loading = $state(false);
	let error = $state<string | null>(null);
	let preview = $state<components['schemas']['ParseProposal'] | null>(null);

	async function parse() {
		loading = true;
		error = null;
		preview = null;
		const { data, error: err } = await faraidClient.POST('/parse', { body: { text } });
		loading = false;
		if (err) {
			error = JSON.stringify(err);
		} else if (data) {
			preview = data;
		}
	}

	function applyPreview() {
		if (!preview) return;
		const sex = preview.deceasedSex === 'female' ? 'female' : 'male';
		const counts: Counts = {};
		for (const [k, v] of Object.entries(preview.heirs ?? {})) {
			if (typeof v === 'number' && v > 0) counts[k] = v;
		}
		onParsed(sex, counts);
		preview = null;
		text = '';
	}
</script>

{#if llmEnabled}
	<div class="mb-6 rounded border border-amber-200 bg-amber-50 p-4">
		<p class="text-xs text-amber-700 mb-2">{$t('llm.nl_warning')}</p>
		<textarea
			bind:value={text}
			placeholder={$t('llm.nl_placeholder')}
			rows="3"
			class="w-full rounded border border-(--color-border) bg-white px-3 py-2 text-sm resize-none"
		></textarea>
		<button
			type="button"
			onclick={parse}
			disabled={loading || text.trim() === ''}
			class="mt-2 rounded bg-amber-600 px-3 py-1.5 text-sm text-white disabled:opacity-50 hover:bg-amber-700"
		>
			{loading ? '...' : $t('llm.nl_parse')}
		</button>

		{#if error}
			<p class="mt-2 text-xs text-red-600">{error}</p>
		{/if}

		{#if preview}
			<div class="mt-3 rounded border border-amber-300 bg-white p-3 text-sm">
				<p class="font-semibold mb-1">Preview: {preview.deceasedSex}</p>
				<ul class="text-xs space-y-0.5">
					{#each Object.entries(preview.heirs ?? {}) as [k, v]}
						{#if typeof v === 'number' && v > 0}
							<li>{k}: {v}</li>
						{/if}
					{/each}
				</ul>
				<button
					type="button"
					onclick={applyPreview}
					class="mt-2 rounded bg-(--color-primary) px-3 py-1 text-xs text-white hover:bg-(--color-primary-hover)"
				>
					{$t('llm.nl_confirm')}
				</button>
			</div>
		{/if}
	</div>
{/if}
