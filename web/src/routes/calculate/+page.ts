// Smoke test for the typed API client: verify that the POST /solve call
// type-checks correctly with the generated schema. This file intentionally
// does not call the backend at load time; the actual fetch happens from
// interactive form submissions in +page.svelte.
import type { components } from '$lib/client';

// If these type aliases compile, the generated client covers /solve correctly.
export type SolveRequest = components['schemas']['SolveRequest'];
export type SolveResult = components['schemas']['SolveResult'];

import type { PageLoad } from './$types';
export const load: PageLoad = () => ({});
