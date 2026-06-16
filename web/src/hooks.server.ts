import type { Handle } from '@sveltejs/kit';

const RTL = new Set(['ar']);

function detectLocale(acceptLanguage: string | null): string {
	if (!acceptLanguage) return 'en';
	const first = acceptLanguage.split(',')[0].split('-')[0].trim().toLowerCase();
	return first === 'ar' ? 'ar' : 'en';
}

export const handle: Handle = async ({ event, resolve }) => {
	const cookie = event.cookies.get('locale');
	const locale = cookie === 'ar' || cookie === 'en'
		? cookie
		: detectLocale(event.request.headers.get('accept-language'));

	const dir = RTL.has(locale) ? 'rtl' : 'ltr';

	// Make locale available to load functions via locals.
	event.locals.locale = locale as 'en' | 'ar';

	return resolve(event, {
		transformPageChunk: ({ html }) =>
			html.replace('%lang%', locale).replace('%dir%', dir)
	});
};
