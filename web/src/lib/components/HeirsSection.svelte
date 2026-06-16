<script lang="ts">
	import { t } from '$lib/i18n';
	import { HEIRS, type Counts } from '$lib/heirs';

	interface Props {
		counts: Counts;
	}
	let { counts = $bindable() }: Props = $props();

	const groups = ['spouses', 'ascendants', 'descendants', 'siblings', 'collaterals'];

	function inc(key: string, max: number) {
		const cur = counts[key] ?? 0;
		if (max === 0 || cur < max) counts = { ...counts, [key]: cur + 1 };
	}

	function dec(key: string) {
		const cur = counts[key] ?? 0;
		if (cur > 0) counts = { ...counts, [key]: cur - 1 };
	}
</script>

<section class="mb-6" aria-labelledby="heirs-heading">
	<h2 id="heirs-heading" class="text-lg font-semibold mb-3">{$t('form.heirs')}</h2>
	{#each groups as group}
		{@const heirs = HEIRS.filter((h) => h.group === group)}
		<div class="mb-4">
			<h3 class="text-xs uppercase tracking-wide text-gray-400 mb-2">{$t(`group.${group}`)}</h3>
			<div class="grid grid-cols-1 sm:grid-cols-2 gap-2">
				{#each heirs as heir}
					{@const count = counts[heir.key] ?? 0}
					<div
						class="flex items-center justify-between rounded border border-(--color-border) px-3 py-2"
						title={$t(heir.helpKey)}
					>
						<div class="flex-1 min-w-0">
							<span class="text-sm font-medium" id="heir-{heir.key}">{$t(heir.labelKey)}</span>
						</div>
						<div class="flex items-center gap-2 ms-3" role="group" aria-labelledby="heir-{heir.key}">
							<button
								type="button"
								onclick={() => dec(heir.key)}
								disabled={count === 0}
								class="w-7 h-7 rounded border border-(--color-border) text-lg leading-none disabled:opacity-30 hover:bg-(--color-border)"
								aria-label="{$t('form.decrease')} {$t(heir.labelKey)}"
							>-</button>
							<span class="w-5 text-center text-sm tabular-nums" aria-live="polite" aria-atomic="true">{count}</span>
							<button
								type="button"
								onclick={() => inc(heir.key, heir.max)}
								disabled={heir.max > 0 && count >= heir.max}
								class="w-7 h-7 rounded border border-(--color-border) text-lg leading-none disabled:opacity-30 hover:bg-(--color-border)"
								aria-label="{$t('form.increase')} {$t(heir.labelKey)}"
							>+</button>
						</div>
					</div>
				{/each}
			</div>
		</div>
	{/each}
</section>
