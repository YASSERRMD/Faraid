import { writable, derived } from 'svelte/store';
import en from './locales/en.json';
import ar from './locales/ar.json';

export type Locale = 'en' | 'ar';

const translations: Record<Locale, Record<string, string>> = { en, ar };

// The supported locales and their display directions.
export const LOCALES: Locale[] = ['en', 'ar'];
export const RTL_LOCALES = new Set<Locale>(['ar']);

// Writable locale store. Components import and set this to switch language.
export const locale = writable<Locale>('en');

// Derived direction: 'ltr' for English, 'rtl' for Arabic.
export const dir = derived(locale, ($l) => (RTL_LOCALES.has($l) ? 'rtl' : 'ltr'));

// t(key) returns the translation for the current locale, falling back to the
// English string when the key is absent from the active translation table.
export const t = derived(locale, ($l) => (key: string): string => {
	return translations[$l][key] ?? translations['en'][key] ?? key;
});
