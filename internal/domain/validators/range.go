package validators

import (
	"errors"

	"github.com/indigowar/map-of-events/internal/domain/models"
)

var (
	ErrLowBiggerThenHighValue      = errors.New("range's low value is bigger then high value")
	ErrRangeHasNegativeValues      = errors.New("range has negative value(s)")
	ErrRangeHasUnacceptablePercent = errors.New("range has a value of unaccepted percent")
)

func ValidateRange(r models.FoundingRange) error {
	high := r.High()
	low := r.Low()

	if high > low {
		return ErrLowBiggerThenHighValue
	}

	if low < 0 || high < 0 {
		return ErrRangeHasNegativeValues
	}

	return nil
}

func ValidatePercentRange(r models.FoundingRange) error {
	high := r.High()
	low := r.Low()

	if err := ValidateRange(r); err != nil {
		return err
	}

	if high > 100 || low > 100 {
		return ErrRangeHasUnacceptablePercent
	}

	return nil
}
