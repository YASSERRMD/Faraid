import { env } from '$env/dynamic/public';

// llmEnabled is true only when PUBLIC_LLM_ENABLED is the exact string "true".
// Anything else (unset, "false", empty) keeps the trial tier hidden.
// Components check this before rendering any LLM affordance.
export const llmEnabled = env.PUBLIC_LLM_ENABLED === 'true';
