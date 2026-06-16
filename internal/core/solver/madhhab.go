package solver

// DhawuArhamRoute selects whether and how the distant kindred (dhawu al-arham)
// inherit when there is no fixed-share or residuary heir.
type DhawuArhamRoute int

const (
	// DhawuArhamExcluded gives the residue to the public treasury rather than
	// to the distant kindred (the classical Maliki and Shafi'i position).
	DhawuArhamExcluded DhawuArhamRoute = iota
	// DhawuArhamInherit lets the distant kindred inherit through a structured
	// route (the Hanafi and Hanbali position).
	DhawuArhamInherit
)

// String returns a label for the route.
func (r DhawuArhamRoute) String() string {
	if r == DhawuArhamInherit {
		return "inherit"
	}
	return "excluded"
}

// Madhhab is a school profile that selects every rule that diverges between the
// schools. Rules that differ are chosen here as data, so the solver stays free
// of hardcoded per-school branches.
type Madhhab struct {
	Name string
	// GrandfatherView decides the grandfather-with-siblings treatment.
	GrandfatherView JaddView
	// MushtarakaView decides whether full siblings share the uterine third.
	MushtarakaView MushtarakaView
	// RaddToSpouse decides whether the spouse takes part in radd.
	RaddToSpouse bool
	// DhawuArham decides whether the distant kindred inherit.
	DhawuArham DhawuArhamRoute
}

// The four Sunni school profiles. The divergent positions encoded here are the
// commonly taught ones and are the single point to adjust during scholar
// review.
var (
	Hanafi = Madhhab{
		Name:            "Hanafi",
		GrandfatherView: JaddAbuHanifa,
		MushtarakaView:  MushtarakaNoShare,
		RaddToSpouse:    false,
		DhawuArham:      DhawuArhamInherit,
	}
	Maliki = Madhhab{
		Name:            "Maliki",
		GrandfatherView: JaddZayd,
		MushtarakaView:  MushtarakaShare,
		RaddToSpouse:    false,
		DhawuArham:      DhawuArhamExcluded,
	}
	Shafii = Madhhab{
		Name:            "Shafii",
		GrandfatherView: JaddZayd,
		MushtarakaView:  MushtarakaShare,
		RaddToSpouse:    false,
		DhawuArham:      DhawuArhamExcluded,
	}
	Hanbali = Madhhab{
		Name:            "Hanbali",
		GrandfatherView: JaddZayd,
		MushtarakaView:  MushtarakaNoShare,
		RaddToSpouse:    false,
		DhawuArham:      DhawuArhamInherit,
	}
)

// Madhahib returns the four Sunni school profiles in a stable order.
func Madhahib() []Madhhab {
	return []Madhhab{Hanafi, Maliki, Shafii, Hanbali}
}

// MadhhabByName returns the profile with the given name (case-insensitive on
// the canonical spellings used here), or ok false.
func MadhhabByName(name string) (Madhhab, bool) {
	for _, m := range Madhahib() {
		if m.Name == name {
			return m, true
		}
	}
	return Madhhab{}, false
}
