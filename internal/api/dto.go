package api

// These data transfer objects mirror the schemas in openapi/faraid.yaml. All
// fractions and amounts are exact rational strings of the form "n/d".

type estateDTO struct {
	Total                       int64 `json:"total"`
	Funeral                     int64 `json:"funeral"`
	Debts                       int64 `json:"debts"`
	Bequests                    int64 `json:"bequests"`
	HeirsConsentToExcessBequest bool  `json:"heirsConsentToExcessBequest"`
}

type solveRequest struct {
	DeceasedSex string         `json:"deceasedSex"`
	Estate      estateDTO      `json:"estate"`
	Heirs       map[string]int `json:"heirs"`
	Madhhab     string         `json:"madhhab"`
}

type heirShareDTO struct {
	Relation string `json:"relation"`
	Count    int    `json:"count"`
	Fraction string `json:"fraction"`
	Parts    int64  `json:"parts"`
	Amount   string `json:"amount"`
}

type derivationStepDTO struct {
	Stage     string `json:"stage"`
	Relation  string `json:"relation,omitempty"`
	Detail    string `json:"detail,omitempty"`
	Reference string `json:"reference,omitempty"`
	Fraction  string `json:"fraction,omitempty"`
}

type solveResultDTO struct {
	Madhhab       string              `json:"madhhab"`
	Distributable string              `json:"distributable"`
	Base          int64               `json:"base"`
	Shares        []heirShareDTO      `json:"shares"`
	Excluded      []string            `json:"excluded,omitempty"`
	SpecialCase   string              `json:"specialCase,omitempty"`
	Awl           bool                `json:"awl"`
	Radd          bool                `json:"radd"`
	Residue       string              `json:"residue,omitempty"`
	NeedsReview   bool                `json:"needsReview"`
	ReviewNotes   []string            `json:"reviewNotes,omitempty"`
	Derivation    []derivationStepDTO `json:"derivation,omitempty"`
}

type madhhabDTO struct {
	Name string `json:"name"`
}

type errorDTO struct {
	Error string `json:"error"`
}
