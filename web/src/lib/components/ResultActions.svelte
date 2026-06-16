<script lang="ts">
	import { t } from '$lib/i18n';
	import { faraidClient } from '$lib/client';
	import type { components } from '$lib/client';

	interface Props {
		request: components['schemas']['SolveRequest'];
	}
	let { request }: Props = $props();

	let loadingPDF = $state(false);
	let errorPDF = $state<string | null>(null);
	let loadingLink = $state(false);
	let errorLink = $state<string | null>(null);
	let copied = $state(false);

	async function downloadPDF() {
		loadingPDF = true;
		errorPDF = null;
		try {
			const resp = await fetch('/api/v1/export?format=pdf', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json', Accept: 'application/pdf' },
				body: JSON.stringify(request)
			});
			if (!resp.ok) {
				errorPDF = $t('export.error');
				return;
			}
			const blob = await resp.blob();
			const url = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = url;
			a.download = 'faraid.pdf';
			a.click();
			URL.revokeObjectURL(url);
		} catch {
			errorPDF = $t('export.error');
		} finally {
			loadingPDF = false;
		}
	}

	async function copyLink() {
		loadingLink = true;
		errorLink = null;
		try {
			const { data, error: err } = await faraidClient.POST('/cases', {
				body: { name: 'Shared case', input: request }
			});
			if (err || !data) {
				errorLink = $t('export.link_error');
				return;
			}
			const shareUrl = `${window.location.origin}${window.location.pathname}?case=${data.id}`;
			await navigator.clipboard.writeText(shareUrl);
			copied = true;
			setTimeout(() => {
				copied = false;
			}, 2000);
		} catch {
			errorLink = $t('export.link_error');
		} finally {
			loadingLink = false;
		}
	}
</script>

<div class="mt-4 flex flex-wrap gap-2">
	<button
		type="button"
		onclick={downloadPDF}
		disabled={loadingPDF}
		class="rounded border border-(--color-border) px-3 py-1.5 text-sm hover:bg-gray-50 disabled:opacity-50"
	>
		{loadingPDF ? '...' : $t('export.download_pdf')}
	</button>

	<button
		type="button"
		onclick={copyLink}
		disabled={loadingLink}
		class="rounded border border-(--color-border) px-3 py-1.5 text-sm hover:bg-gray-50 disabled:opacity-50"
	>
		{loadingLink ? '...' : copied ? $t('export.copied') : $t('export.copy_link')}
	</button>
</div>

{#if errorPDF}
	<p class="mt-1 text-xs text-red-600">{errorPDF}</p>
{/if}
{#if errorLink}
	<p class="mt-1 text-xs text-red-600">{errorLink}</p>
{/if}
