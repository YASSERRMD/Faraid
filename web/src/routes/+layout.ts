import type { LayoutLoad } from './$types';

export const load: LayoutLoad = ({ data }) => {
	return { locale: data?.locale ?? 'en' };
};
