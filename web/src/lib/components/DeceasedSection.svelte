<script lang="ts">
	import { t } from '$lib/i18n';

	type Madhhab = 'Hanafi' | 'Maliki' | 'Shafii' | 'Hanbali';
	const madhahib: Madhhab[] = ['Hanafi', 'Maliki', 'Shafii', 'Hanbali'];

	interface Props {
		sex: 'male' | 'female';
		madhhab: Madhhab;
	}

	let { sex = $bindable(), madhhab = $bindable() }: Props = $props();
</script>

<section class="mb-6">
	<h2 class="text-lg font-semibold mb-3">{$t('form.deceased_sex')}</h2>
	<div class="flex gap-4 mb-4">
		<label class="flex items-center gap-2 cursor-pointer">
			<input type="radio" name="sex" value="male" bind:group={sex} class="accent-(--color-primary)" />
			<span>{$t('form.deceased_male')}</span>
		</label>
		<label class="flex items-center gap-2 cursor-pointer">
			<input type="radio" name="sex" value="female" bind:group={sex} class="accent-(--color-primary)" />
			<span>{$t('form.deceased_female')}</span>
		</label>
	</div>

	<label class="block">
		<span class="text-sm font-medium">{$t('form.madhhab')}</span>
		<select
			bind:value={madhhab}
			class="mt-1 block w-full rounded border border-(--color-border) bg-(--color-surface) px-3 py-2 text-sm"
		>
			{#each madhahib as m}
				<option value={m}>{$t(`madhhab.${m}`)}</option>
			{/each}
		</select>
	</label>
</section>
