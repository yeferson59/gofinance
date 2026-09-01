package decimal

import "errors"

var (
	ErrOverflow        = errors.New("numeric overflow")
	ErrDivideByZero    = errors.New("division by zero")
	ErrEmptyString     = errors.New("can't parse empty string")
	ErrInvalidFormat   = errors.New("invalid decimal format")
	ErrPrecOutOfRange  = errors.New("precision out of range, maximum is 19 digits after the decimal point")
	ErrIntPartOverflow = errors.New("integer part is too large to fit in int64")
	ErrLogNonPositive  = errors.New("logarithm of a non-positive number is undefined")
	ErrPowNegBase      = errors.New("negative base with a non-integer exponent is undefined")
	ErrSqrtNegative    = errors.New("square root of a negative number is undefined")
	ErrSymbolFraction  = errors.New("fraction must contain exactly one '/'")

	// ErrInvalidBinary is returned by UnmarshalBinary when the bytes are not a
	// decimal this package wrote: a wrong length, or reserved bits in use.
	ErrInvalidBinary = errors.New("invalid binary encoding")

	// ErrUnknownBinaryVersion is returned by UnmarshalBinary when the leading
	// version byte names a layout this build does not know, which is what
	// stops it from reading a newer encoding as if it were this one.
	ErrUnknownBinaryVersion = errors.New("unknown binary encoding version")
)
