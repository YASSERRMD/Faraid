package estate

import (
	"errors"

	"github.com/YASSERRMD/Faraid/internal/core/heir"
)

// Case is the complete input to an inheritance computation: the deceased's
// sex, the estate, and the surviving heirs.
type Case struct {
	DeceasedSex heir.Sex
	Estate      Estate
	Heirs       *heir.Heirs
}

// Validate checks the case for internal consistency: a known deceased sex,
// valid estate amounts, valid heirs, and spouse heirs consistent with the
// deceased's sex. A deceased man leaves wives, a deceased woman leaves a
// husband.
func (c Case) Validate() error {
	switch c.DeceasedSex {
	case heir.Male, heir.Female:
		// ok
	default:
		return errors.New("estate: deceased sex must be male or female")
	}
	if err := c.Estate.validate(); err != nil {
		return err
	}
	if c.Heirs == nil {
		return errors.New("estate: heirs must not be nil")
	}
	if err := c.Heirs.Validate(); err != nil {
		return err
	}
	if c.DeceasedSex == heir.Male && c.Heirs.Present(heir.Husband) {
		return errors.New("estate: a deceased male cannot leave a husband")
	}
	if c.DeceasedSex == heir.Female && c.Heirs.Present(heir.Wife) {
		return errors.New("estate: a deceased female cannot leave a wife")
	}
	return nil
}
