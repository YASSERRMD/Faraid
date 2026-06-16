<script lang="ts">
	import { t } from '$lib/i18n';
</script>

<svelte:head>
	<title>Faraid</title>
</svelte:head>

<div class="text-center py-12">
	<h1 class="text-4xl font-bold text-(--color-primary) mb-3">{$t('app.title')}</h1>
	<p class="text-xl text-gray-600 mb-2">{$t('app.subtitle')}</p>
	<p class="text-sm text-gray-400 mb-10">{$t('app.tagline')}</p>
	<a
		href="/calculate"
		class="inline-block rounded bg-(--color-primary) px-6 py-3 text-white font-semibold hover:bg-(--color-primary-hover) transition-colors"
	>
		{$t('nav.calculate')}
	</a>
</div>
