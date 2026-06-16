<script lang="ts">
	import { t } from '$lib/i18n';
	import type { components } from '$lib/client';

	type DerivationStep = components['schemas']['DerivationStep'];

	interface Props {
		steps: DerivationStep[];
	}
	let { steps }: Props = $props();

	const STAGE_ORDER = ['estate', 'special_case', 'blocking', 'fixed_share', 'residuary', 'awl', 'radd', 'asl', 'result'];

	const grouped = $derived(
		STAGE_ORDER.map((stage) => ({
			stage,
			steps: steps.filter((s) => s.stage === stage)
		})).filter((g) => g.steps.length > 0)
	);

	let open = $state<Record<string, boolean>>({});
	function isOpen(stage: string) { return open[stage] !== false; }
	function toggle(stage: string) { open = { ...open, [stage]: !isOpen(stage) }; }
</script>

<section class="mt-8" aria-labelledby="derivation-heading">
	<h2 id="derivation-heading" class="text-xl font-bold mb-4">{$t('result.derivation')}</h2>
	<div class="space-y-2">
		{#each grouped as group}
			{@const panelId = `derivation-panel-${group.stage}`}
			{@const btnId = `derivation-btn-${group.stage}`}
			<div class="rounded border border-(--color-border)">
				<button
					id={btnId}
					type="button"
					onclick={() => toggle(group.stage)}
					class="flex w-full items-center justify-between px-4 py-2 text-start text-sm font-semibold hover:bg-gray-50"
					aria-expanded={isOpen(group.stage)}
					aria-controls={panelId}
				>
					<span>{$t(`stage.${group.stage}`)}</span>
					<span class="text-gray-400 text-xs" aria-hidden="true">{isOpen(group.stage) ? '▲' : '▼'}</span>
				</button>
				{#if isOpen(group.stage)}
					<ul
						id={panelId}
						role="list"
						aria-labelledby={btnId}
						class="divide-y divide-(--color-border) border-t border-(--color-border)"
					>
						{#each group.steps as step}
							<li class="px-4 py-2 text-sm">
								<div class="flex flex-wrap items-baseline gap-x-2">
									{#if step.relation}
										<span class="font-medium">{$t(`heir.${step.relation}`)}</span>
									{/if}
									<span class="text-gray-600">{step.detail}</span>
									{#if step.fraction}
										<span class="ms-auto font-mono tabular-nums text-(--color-primary)">{step.fraction}</span>
									{/if}
								</div>
								{#if step.reference}
									<p class="mt-0.5 text-xs text-gray-400">{step.reference}</p>
								{/if}
							</li>
						{/each}
					</ul>
				{/if}
			</div>
		{/each}
	</div>
</section>
