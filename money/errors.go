package money

import "errors"

var (
	// ErrInvalidISOCode is returned when a currency ISO code cannot be
	// recognized.
	ErrInvalidISOCode = errors.New("invalid iso code")

	// ErrNoAllocationRatios is returned by Allocate when called without any
	// ratios to split by.
	ErrNoAllocationRatios = errors.New("money: no allocation ratios given")

	// ErrZeroAllocationRatios is returned by Allocate when every ratio is zero,
	// making the split undefined.
	ErrZeroAllocationRatios = errors.New("money: allocation ratios sum to zero")

	// ErrInvalidAllocationCount is returned by AllocateEvenly when asked to
	// split into zero or fewer parts.
	ErrInvalidAllocationCount = errors.New("money: allocation count must be positive")

	// ErrInvalidExchangeRate is returned by Convert when given a rate that is
	// zero or negative.
	ErrInvalidExchangeRate = errors.New("money: exchange rate must be positive")

	// ErrCurrencyMismatch is returned by operations that require both operands
	// to share the same currency, such as SafeAdd and SafeSub.
	ErrCurrencyMismatch = errors.New("money: currency mismatch")

	// ErrTrailingJSONContent is returned by UnmarshalJSON when something
	// follows the JSON value it decoded. json.Unmarshal never passes such
	// input, but these methods are also called directly.
	ErrTrailingJSONContent = errors.New("money: unexpected content after the JSON value")

	// ErrMissingJSONValue is returned when a JSON object describing an amount
	// has no "value" member. The currency may be left out — it defaults to USD
	// — but the amount cannot.
	ErrMissingJSONValue = errors.New("money: JSON object has no \"value\" member")

	// ErrInvalidBinary is returned by UnmarshalBinary when the bytes are not
	// something this package wrote: too short, or the wrong length.
	ErrInvalidBinary = errors.New("money: invalid binary encoding")

	// ErrUnknownBinaryVersion is returned by UnmarshalBinary when the leading
	// version byte names a layout this build does not know, which is what
	// stops it from reading a newer encoding as if it were this one.
	ErrUnknownBinaryVersion = errors.New("money: unknown binary encoding version")
)
