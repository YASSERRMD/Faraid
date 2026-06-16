<script lang="ts">
	import '../app.css';
	import { locale, t, dir } from '$lib/i18n';
	import type { LayoutData } from './$types';

	let { children, data }: { children: import('svelte').Snippet; data: LayoutData } = $props();

	// Sync the server-detected locale into the client-side store on first render.
	$effect(() => {
		locale.set(data.locale);
	});

	function switchLocale() {
		const next = $locale === 'en' ? 'ar' : 'en';
		document.cookie = `locale=${next}; path=/; max-age=31536000; SameSite=Lax`;
		locale.set(next);
		document.documentElement.lang = next;
		document.documentElement.dir = next === 'ar' ? 'rtl' : 'ltr';
	}
</script>

<a
	href="#main-content"
	class="sr-only focus:not-sr-only focus:fixed focus:start-4 focus:top-4 focus:z-50 focus:rounded focus:bg-(--color-primary) focus:px-4 focus:py-2 focus:text-white focus:no-underline"
>
	{$t('a11y.skip_to_content')}
</a>

<header class="border-b border-(--color-border) bg-(--color-surface)">
	<nav class="mx-auto flex max-w-4xl items-center justify-between px-4 py-3" aria-label={$t('a11y.nav_label')}>
		<a href="/" class="text-xl font-bold text-(--color-primary)">{$t('app.title')}</a>
		<div class="flex items-center gap-6 text-sm">
			<a href="/" class="hover:text-(--color-primary)">{$t('nav.calculate')}</a>
			<a href="/compare" class="hover:text-(--color-primary)">{$t('nav.compare')}</a>
			<button
				onclick={switchLocale}
				class="rounded border border-(--color-border) px-3 py-1 text-xs hover:bg-(--color-border) transition-colors"
				aria-label={$t('a11y.switch_language')}
			>
				{$t('lang.switch')}
			</button>
		</div>
	</nav>
</header>

<main id="main-content" class="mx-auto max-w-4xl px-4 py-8" dir={$dir as 'ltr' | 'rtl'}>
	{@render children()}
</main>
