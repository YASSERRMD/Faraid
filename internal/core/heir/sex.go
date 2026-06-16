package heir

// Sex is the legal sex of an heir. It determines the residuary 2 to 1
// distribution and several share rules. SexUnknown is reserved for the
// intersex (khuntha) case, which a later phase resolves with a cautious
// scenario based procedure.
type Sex uint8

const (
	SexUnknown Sex = iota
	Male
	Female
)

// String returns a lowercase label for the sex.
func (s Sex) String() string {
	switch s {
	case Male:
		return "male"
	case Female:
		return "female"
	default:
		return "unknown"
	}
}
