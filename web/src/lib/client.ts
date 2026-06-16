import createClient from 'openapi-fetch';
import type { paths } from './api.gen';

// faraidClient is the single typed HTTP client for the Faraid backend. It
// targets /api/v1 relative to the current origin, so it works in both the
// dev proxy and the production deployment without configuration.
export const faraidClient = createClient<paths>({ baseUrl: '/api/v1' });

// Re-export the schema types that components use directly so they only need
// one import.
export type { components, operations } from './api.gen';
