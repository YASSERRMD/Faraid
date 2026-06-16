// Heir metadata: API key, translation key, max count (0 = unbounded), group.
export interface HeirDef {
	key: string;
	labelKey: string;
	helpKey: string;
	max: number;
	group: string;
}

export const HEIR_GROUPS = ['spouses', 'ascendants', 'descendants', 'siblings', 'collaterals'];

export const HEIRS: HeirDef[] = [
	// Spouses
	{ key: 'husband', labelKey: 'heir.husband', helpKey: 'help.husband', max: 1, group: 'spouses' },
	{ key: 'wife', labelKey: 'heir.wife', helpKey: 'help.wife', max: 4, group: 'spouses' },
	// Ascendants
	{ key: 'father', labelKey: 'heir.father', helpKey: 'help.father', max: 1, group: 'ascendants' },
	{ key: 'mother', labelKey: 'heir.mother', helpKey: 'help.mother', max: 1, group: 'ascendants' },
	{ key: 'paternal_grandfather', labelKey: 'heir.paternal_grandfather', helpKey: 'help.paternal_grandfather', max: 1, group: 'ascendants' },
	{ key: 'paternal_grandmother', labelKey: 'heir.paternal_grandmother', helpKey: 'help.paternal_grandmother', max: 1, group: 'ascendants' },
	{ key: 'maternal_grandmother', labelKey: 'heir.maternal_grandmother', helpKey: 'help.maternal_grandmother', max: 1, group: 'ascendants' },
	// Descendants
	{ key: 'son', labelKey: 'heir.son', helpKey: 'help.son', max: 0, group: 'descendants' },
	{ key: 'daughter', labelKey: 'heir.daughter', helpKey: 'help.daughter', max: 0, group: 'descendants' },
	{ key: 'sons_son', labelKey: 'heir.sons_son', helpKey: 'help.sons_son', max: 0, group: 'descendants' },
	{ key: 'sons_daughter', labelKey: 'heir.sons_daughter', helpKey: 'help.sons_daughter', max: 0, group: 'descendants' },
	// Siblings
	{ key: 'full_brother', labelKey: 'heir.full_brother', helpKey: 'help.full_brother', max: 0, group: 'siblings' },
	{ key: 'full_sister', labelKey: 'heir.full_sister', helpKey: 'help.full_sister', max: 0, group: 'siblings' },
	{ key: 'consanguine_brother', labelKey: 'heir.consanguine_brother', helpKey: 'help.consanguine_brother', max: 0, group: 'siblings' },
	{ key: 'consanguine_sister', labelKey: 'heir.consanguine_sister', helpKey: 'help.consanguine_sister', max: 0, group: 'siblings' },
	{ key: 'uterine_brother', labelKey: 'heir.uterine_brother', helpKey: 'help.uterine_brother', max: 0, group: 'siblings' },
	{ key: 'uterine_sister', labelKey: 'heir.uterine_sister', helpKey: 'help.uterine_sister', max: 0, group: 'siblings' },
	// Collaterals
	{ key: 'full_brothers_son', labelKey: 'heir.full_brothers_son', helpKey: 'help.full_brothers_son', max: 0, group: 'collaterals' },
	{ key: 'consanguine_brothers_son', labelKey: 'heir.consanguine_brothers_son', helpKey: 'help.consanguine_brothers_son', max: 0, group: 'collaterals' },
	{ key: 'full_paternal_uncle', labelKey: 'heir.full_paternal_uncle', helpKey: 'help.full_paternal_uncle', max: 0, group: 'collaterals' },
	{ key: 'consanguine_paternal_uncle', labelKey: 'heir.consanguine_paternal_uncle', helpKey: 'help.consanguine_paternal_uncle', max: 0, group: 'collaterals' },
	{ key: 'full_paternal_uncles_son', labelKey: 'heir.full_paternal_uncles_son', helpKey: 'help.full_paternal_uncles_son', max: 0, group: 'collaterals' },
	{ key: 'consanguine_paternal_uncles_son', labelKey: 'heir.consanguine_paternal_uncles_son', helpKey: 'help.consanguine_paternal_uncles_son', max: 0, group: 'collaterals' }
];

export type Counts = Record<string, number>;

// validate returns a list of error message keys for the current heir
// configuration, mirroring the backend validation rules.
export function validate(counts: Counts, sex: 'male' | 'female'): string[] {
	const errors: string[] = [];
	const get = (k: string) => counts[k] ?? 0;

	if (get('husband') > 0 && get('wife') > 0) errors.push('validation.husband_wife');
	if (sex === 'male' && get('husband') > 0) errors.push('validation.male_no_husband');
	if (sex === 'female' && get('wife') > 0) errors.push('validation.female_no_wife');

	for (const heir of HEIRS) {
		const n = get(heir.key);
		if (n < 0) { errors.push('validation.non_negative'); break; }
		if (heir.max > 0 && n > heir.max) {
			errors.push(heir.max === 1 ? 'validation.max_one' : 'validation.max_four');
			break;
		}
	}
	return errors;
}
