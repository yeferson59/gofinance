package decimal

import "errors"

var (
	// ErrOverflow is returned when a result does not fit in the 128-bit
	// coefficient.
	ErrOverflow = errors.New("numeric overflow")
	// ErrDivideByZero is returned when dividing by zero.
	ErrDivideByZero = errors.New("division by zero")
	// ErrEmptyString is returned when parsing an empty string.
	ErrEmptyString = errors.New("can't parse empty string")
	// ErrInvalidFormat is returned for text that is not a plain decimal number,
	// and for NaN and ±Inf floats.
	ErrInvalidFormat = errors.New("invalid decimal format")
	// ErrPrecOutOfRange is returned for a precision above 19 fractional digits.
	ErrPrecOutOfRange = errors.New("precision out of range, maximum is 19 digits after the decimal point")
	// ErrIntPartOverflow is returned by Int64 when the integer part does not fit
	// in an int64.
	ErrIntPartOverflow = errors.New("integer part is too large to fit in int64")
	// ErrLogNonPositive is returned by the logarithms for a zero or negative
	// argument.
	ErrLogNonPositive = errors.New("logarithm of a non-positive number is undefined")
	// ErrPowNegBase is returned by Pow for a negative base with a non-integer
	// exponent.
	ErrPowNegBase = errors.New("negative base with a non-integer exponent is undefined")
	// ErrSqrtNegative is returned by Sqrt for a negative argument.
	ErrSqrtNegative = errors.New("square root of a negative number is undefined")
	// ErrSymbolFraction is returned by NewFractionString when the text does not
	// contain exactly one '/'.
	ErrSymbolFraction = errors.New("fraction must contain exactly one '/'")

	// ErrInvalidBinary is returned by UnmarshalBinary when the bytes are not a
	// decimal this package wrote: a wrong length, or reserved bits in use.
	ErrInvalidBinary = errors.New("invalid binary encoding")

	// ErrUnknownBinaryVersion is returned by UnmarshalBinary when the leading
	// version byte names a layout this build does not know, which is what
	// stops it from reading a newer encoding as if it were this one.
	ErrUnknownBinaryVersion = errors.New("unknown binary encoding version")

	// ErrScanNull is returned by Scan for a NULL column, which has no Decimal
	// value. Scan into a *Decimal to read a nullable column.
	ErrScanNull = errors.New("decimal: can't scan NULL into a Decimal; scan into *Decimal for a nullable column")
)
