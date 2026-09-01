package decimal

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/json/jsontext"
	json "encoding/json/v2"
	"errors"
	"fmt"
	"strings"
)

var Zero = Decimal{value: decZero}
var One = Decimal{value: decOne}

// Decimal is a fixed-point decimal number with up to 19 digits of
// precision after the decimal point, backed by a 128-bit coefficient.
// It performs arithmetic without heap allocations or external
// dependencies.
type Decimal struct {
	value decimal128
}

func NewFromFloat64(f float64) (Decimal, error) {
	decimal, err := decFromFloat64(f)

	return Decimal{decimal}, err
}

func NewFromInt64(coef int64, prec uint8) (Decimal, error) {
	decimal, err := decFromInt64(coef, prec)

	return Decimal{decimal}, err
}

func NewFromUint64(coef uint64, prec uint8) (Decimal, error) {
	decimal, err := decFromUint64(coef, prec)

	return Decimal{decimal}, err
}

func NewFromHiLo(neg bool, hi, lo uint64, prec uint8) (Decimal, error) {
	decimal, err := decFromHiLo(neg, hi, lo, prec)

	return Decimal{decimal}, err
}

func NewFromString(s string) (Decimal, error) {
	decimal, err := parseDecimal(s)

	return Decimal{decimal}, err
}

// newFraction divides num by den, wrapping the result in a Decimal. It
// returns ErrDivideByZero if den is zero.
func newFraction(num, den decimal128) (Decimal, error) {
	result, err := num.Div(den)

	return Decimal{result}, err
}

// NewFractionFloat64 returns the Decimal value of numerator / denominator.
// It returns ErrDivideByZero if denominator is zero.
func NewFractionFloat64(numerator float64, denominator float64) (Decimal, error) {
	num, err := decFromFloat64(numerator)
	if err != nil {
		return Decimal{}, err
	}

	den, err := decFromFloat64(denominator)
	if err != nil {
		return Decimal{}, err
	}

	return newFraction(num, den)
}

// NewFractionInt64 returns the Decimal value of numerator / denominator,
// with each operand scaled by its own precision (e.g. numerator=150,
// numeratorPrec=2 represents 1.50). It returns ErrDivideByZero if
// denominator is zero.
func NewFractionInt64(numerator int64, numeratorPrec uint8, denominator int64, denominatorPrec uint8) (Decimal, error) {
	num, err := decFromInt64(numerator, numeratorPrec)
	if err != nil {
		return Decimal{}, err
	}

	den, err := decFromInt64(denominator, denominatorPrec)
	if err != nil {
		return Decimal{}, err
	}

	return newFraction(num, den)
}

// NewFractionString parses a "numerator/denominator" string (e.g. "1/3")
// and returns the resulting Decimal value. Leading/trailing whitespace
// around the whole string and around each operand is ignored. It returns
// ErrSymbolFraction if frac doesn't contain exactly one '/', and
// ErrDivideByZero if the denominator is zero.
func NewFractionString(frac string) (Decimal, error) {
	frac = strings.TrimSpace(frac)

	sep := "/"

	if strings.Count(frac, sep) != 1 {
		return Decimal{}, ErrSymbolFraction
	}

	split := strings.SplitN(frac, sep, 2)

	num, err := parseDecimal(strings.TrimSpace(split[0]))
	if err != nil {
		return Decimal{}, err
	}

	den, err := parseDecimal(strings.TrimSpace(split[1]))
	if err != nil {
		return Decimal{}, err
	}

	return newFraction(num, den)
}

func MustFromFloat64(f float64) Decimal {
	d, err := NewFromFloat64(f)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFromInt64(coef int64, prec uint8) Decimal {
	d, err := NewFromInt64(coef, prec)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFromUint64(coef uint64, prec uint8) Decimal {
	d, err := NewFromUint64(coef, prec)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFromHiLo(neg bool, hi, lo uint64, prec uint8) Decimal {
	d, err := NewFromHiLo(neg, hi, lo, prec)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFromString(s string) Decimal {
	d, err := NewFromString(s)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFractionFloat64(numerator float64, denominator float64) Decimal {
	d, err := NewFractionFloat64(numerator, denominator)
	if err != nil {
		panic(err)
	}

	return d
}

func MustFractionInt64(numerator int64, numeratorPrec uint8, denominator int64, denominatorPrec uint8) Decimal {
	d, err := NewFractionInt64(numerator, numeratorPrec, denominator, denominatorPrec)
	if err != nil {
		panic(err)
	}

	return d
}

// TryAdd is like Add but returns an error instead of panicking on overflow.
func (d Decimal) TryAdd(other Decimal) (Decimal, error) {
	v, err := d.value.Add(other.value)

	return Decimal{v}, err
}

func (d Decimal) Add(other Decimal) Decimal {
	v, err := d.TryAdd(other)
	if err != nil {
		panic(err)
	}

	return v
}

// TrySub is like Sub but returns an error instead of panicking on overflow.
func (d Decimal) TrySub(other Decimal) (Decimal, error) {
	v, err := d.value.Sub(other.value)

	return Decimal{v}, err
}

func (d Decimal) Sub(other Decimal) Decimal {
	v, err := d.TrySub(other)
	if err != nil {
		panic(err)
	}

	return v
}

// TryMul is like Mul but returns an error instead of panicking on overflow.
func (d Decimal) TryMul(other Decimal) (Decimal, error) {
	v, err := d.value.Mul(other.value)

	return Decimal{v}, err
}

func (d Decimal) Mul(other Decimal) Decimal {
	v, err := d.TryMul(other)
	if err != nil {
		panic(err)
	}

	return v
}

func (d Decimal) Div(other Decimal) (Decimal, error) {
	div, err := d.value.Div(other.value)

	return Decimal{div}, err
}

func (d Decimal) MustDiv(other Decimal) Decimal {
	div, err := d.Div(other)
	if err != nil {
		panic(err)
	}

	return div
}

// Pow returns d raised to the power of exponent, computed natively on the
// decimal engine with 120-bit binary fixed-point internals, so the result
// is accurate to Decimal's full 19-digit precision (rounded half away from
// zero). Conventions: d^0 = 1 (including 0^0); 0^exponent = 0 for a
// positive exponent and ErrDivideByZero for a negative one; a negative
// base requires an integer exponent (the result's sign follows its
// parity), otherwise ErrPowNegBase is returned. Results beyond Decimal's
// range return ErrOverflow; results too small to represent round to zero.
func (d Decimal) Pow(exponent Decimal) (Decimal, error) {
	v, err := d.value.Pow(exponent.value)

	return Decimal{v}, err
}

// MustPow is like Pow but panics on error.
func (d Decimal) MustPow(exponent Decimal) Decimal {
	result, err := d.Pow(exponent)
	if err != nil {
		panic(err)
	}

	return result
}

// Sqrt returns the square root of d, computed directly with Newton's
// integer iteration on the exact 256-bit radicand, so the result is always
// the 19-fractional-digit value nearest to the true root (exact for
// perfect squares, e.g. Sqrt(2.25) = 1.5). It returns ErrSqrtNegative if
// d is negative. For other roots, use Pow (e.g. d.Pow(1/3) for cube roots).
func (d Decimal) Sqrt() (Decimal, error) {
	v, err := d.value.Sqrt()

	return Decimal{v}, err
}

// MustSqrt is like Sqrt but panics on error.
func (d Decimal) MustSqrt() Decimal {
	v, err := d.Sqrt()
	if err != nil {
		panic(err)
	}

	return v
}

// Ln returns the natural logarithm (base e) of d, accurate to Decimal's
// full 19-digit precision (rounded half away from zero). It returns
// ErrLogNonPositive if d is zero or negative, since the natural logarithm
// is undefined there.
func (d Decimal) Ln() (Decimal, error) {
	v, err := d.value.Ln()

	return Decimal{v}, err
}

// MustLn is like Ln but panics on error.
func (d Decimal) MustLn() Decimal {
	result, err := d.Ln()
	if err != nil {
		panic(err)
	}

	return result
}

// Log10 returns the base-10 logarithm of d, accurate to Decimal's full
// 19-digit precision. It returns ErrLogNonPositive if d isn't strictly
// positive.
func (d Decimal) Log10() (Decimal, error) {
	v, err := d.value.Log10()

	return Decimal{v}, err
}

// MustLog10 is like Log10 but panics on error.
func (d Decimal) MustLog10() Decimal {
	result, err := d.Log10()
	if err != nil {
		panic(err)
	}

	return result
}

// Log2 returns the base-2 logarithm of d, accurate to Decimal's full
// 19-digit precision. It returns ErrLogNonPositive if d isn't strictly
// positive.
func (d Decimal) Log2() (Decimal, error) {
	v, err := d.value.Log2()

	return Decimal{v}, err
}

// MustLog2 is like Log2 but panics on error.
func (d Decimal) MustLog2() Decimal {
	result, err := d.Log2()
	if err != nil {
		panic(err)
	}

	return result
}

// Log returns the logarithm of d in the given base, computed as
// Ln(d) / Ln(base) and accurate to Decimal's full 19-digit precision. It
// returns ErrLogNonPositive if d or base aren't strictly positive, and
// ErrDivideByZero if base equals 1 (which makes the logarithm undefined).
func (d Decimal) Log(base Decimal) (Decimal, error) {
	v, err := d.value.Log(base.value)

	return Decimal{v}, err
}

// MustLog is like Log but panics on error.
func (d Decimal) MustLog(base Decimal) Decimal {
	result, err := d.Log(base)
	if err != nil {
		panic(err)
	}

	return result
}

func (d Decimal) Mod(other Decimal) (Decimal, error) {
	mod, err := d.value.Mod(other.value)

	return Decimal{mod}, err
}

func (d Decimal) Div64(other uint64) (Decimal, error) {
	div, err := d.value.Div64(other)

	return Decimal{div}, err
}

func (d Decimal) InexactFloat64() float64 {
	return d.value.InexactFloat64()
}

func (d Decimal) RoundBank(prec uint8) Decimal {
	return Decimal{d.value.RoundBank(prec)}
}

func (d Decimal) RoundAway(prec uint8) Decimal {
	return Decimal{d.value.RoundAway(prec)}
}

func (d Decimal) RoundHAZ(prec uint8) Decimal {
	return Decimal{d.value.RoundHAZ(prec)}
}

func (d Decimal) RoundHTZ(prec uint8) Decimal {
	return Decimal{d.value.RoundHTZ(prec)}
}

func (d Decimal) Trunc(prec uint8) Decimal {
	return Decimal{d.value.Trunc(prec)}
}

func (d Decimal) Abs() Decimal {
	return Decimal{d.value.Abs()}
}

func (d Decimal) Floor() Decimal {
	return Decimal{d.value.Floor()}
}

func (d Decimal) Ceil() Decimal {
	return Decimal{d.value.Ceil()}
}

func (d Decimal) Neg() Decimal {
	return Decimal{d.value.Neg()}
}

func (d Decimal) Sign() int {
	return d.value.Sign()
}

func (d Decimal) Cmp(other Decimal) int {
	return d.value.Cmp(other.value)
}

func (d Decimal) String() string {
	return d.value.String()
}

func (d Decimal) StringFixed(prec uint8) string {
	return d.value.StringFixed(prec)
}

func (d Decimal) Float64() (float64, error) {
	return d.value.InexactFloat64(), nil
}

func (d Decimal) Int64() (int64, error) {
	return d.value.Int64()
}

func (d Decimal) IsZero() bool {
	return d.value.IsZero()
}

func (d Decimal) IsNeg() bool {
	return d.value.IsNeg()
}

func (d Decimal) IsPos() bool {
	return d.value.IsPos()
}

func (d Decimal) GreaterThan(other Decimal) bool {
	return d.value.GreaterThan(other.value)
}

func (d Decimal) GreaterThanOrEqual(other Decimal) bool {
	return d.value.GreaterThanOrEqual(other.value)
}

func (d Decimal) LessThan(other Decimal) bool {
	return d.value.LessThan(other.value)
}

func (d Decimal) LessThanOrEqual(other Decimal) bool {
	return d.value.LessThanOrEqual(other.value)
}

func (d Decimal) Equal(other Decimal) bool {
	return d.value.Equal(other.value)
}

// MarshalJSON implements json.Marshaler, writing the decimal as a bare JSON
// number carrying every digit it holds.
//
// The text String produces is already a valid JSON number — a plain literal,
// never scientific notation — so encoding is a copy rather than a conversion,
// and no precision is lost the way it would be through float64.
//
// This is the encoding/json v1 entry point. Under encoding/json/v2,
// MarshalJSONTo takes precedence and produces the same number, except where
// JSON itself requires a string.
func (d Decimal) MarshalJSON() ([]byte, error) {
	return d.AppendText(nil)
}

// MarshalJSONTo implements json.MarshalerTo, the streaming encoder
// encoding/json/v2 prefers. It writes into the encoder's own buffer instead of
// allocating a fresh slice per value.
//
// It writes a bare JSON number, except in the two positions where JSON has no
// number to offer:
//
//   - as an object name, since a JSON object member name must be a string.
//     This is what lets a Decimal be used as a map key; MarshalJSON alone
//     cannot, because a decoder rejects the number it returns.
//   - when StringifyNumbers is set, which is what the `json:",string"` struct
//     tag turns on. Producers use it precisely so a decimal survives readers
//     that parse every JSON number as a float64.
//
// In both cases the digits are identical; only the quotes differ, and
// UnmarshalJSONFrom reads either form back to the same value.
func (d Decimal) MarshalJSONTo(enc *jsontext.Encoder) error {
	quoted := jsonWantsString(enc)

	// AvailableBuffer hands back the encoder's spare capacity, so the encoded
	// number needs no slice of its own the way MarshalJSON's return value
	// does. The digits themselves still cost the string String builds.
	buf := enc.AvailableBuffer()

	if quoted {
		buf = append(buf, '"')
	}

	buf, err := d.AppendText(buf)
	if err != nil {
		return err
	}

	if quoted {
		buf = append(buf, '"')
	}

	// The text is digits, an optional minus sign and at most one point, so it
	// needs no JSON escaping in either form.
	return enc.WriteValue(buf)
}

// jsonWantsString reports whether the value about to be written must be a JSON
// string rather than a bare number.
//
// StringifyNumbers is the documented option for this and covers the
// `json:",string"` tag, but encoding/json/v2 does not set it for map keys: it
// simply rejects a non-string object name. The position on the encoder's stack
// is what distinguishes them — inside an object, an even number of tokens
// written so far means the next one is a name, since names and values are
// counted separately.
func jsonWantsString(enc *jsontext.Encoder) bool {
	if stringify, _ := json.GetOption(enc.Options(), json.StringifyNumbers); stringify {
		return true
	}

	kind, length := enc.StackIndex(enc.StackDepth())

	return kind == jsontext.KindBeginObject && length%2 == 0
}

// UnmarshalJSON implements json.Unmarshaler.
//
// It accepts a JSON number and a JSON string holding the same digits. The
// string form is not what MarshalJSON writes, but it is what a producer sends
// to protect precision from readers that treat every JSON number as a float64
// — JavaScript above all — and it is the form a Decimal takes as a map key.
// Rejecting it would make this package unable to read documents it can write.
//
// The digits are parsed by the same strict parser NewFromString uses, so a
// quoted value is held to the same standard as an unquoted one: exponent
// notation, padding and more than 19 fractional digits are all rejected.
//
// This is the encoding/json v1 entry point; under encoding/json/v2,
// UnmarshalJSONFrom takes precedence and accepts the same two forms.
func (d *Decimal) UnmarshalJSON(data []byte) error {
	dec, err := parseDecimalJSON(data)
	if err != nil {
		return err
	}

	d.value = dec

	return nil
}

// UnmarshalJSONFrom implements json.UnmarshalerFrom, the streaming decoder
// encoding/json/v2 prefers. It accepts the same forms as UnmarshalJSON.
//
// It reads a whole JSON value rather than a single token, so an object or
// array is consumed before being rejected, leaving the decoder positioned on
// the next value as the interface requires.
func (d *Decimal) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	value, err := dec.ReadValue()
	if err != nil {
		return err
	}

	parsed, err := parseDecimalJSONValue(value)
	if err != nil {
		return err
	}

	d.value = parsed

	return nil
}

// MarshalText implements encoding.TextMarshaler, writing the decimal in the
// same plain form String produces: an optional minus sign, digits, and at most
// one decimal point, with no quotes and no exponent.
//
// MarshalJSON already covers JSON, but it is not consulted by encoders that
// work in plain text — YAML, TOML, XML, flag.TextVar, log/slog — which
// otherwise cannot encode a Decimal at all, since its only field is
// unexported. The text pair serves them from the one canonical form the
// package already uses for String and Value.
//
// JSON does not go through here: MarshalJSONTo and MarshalJSON take
// precedence, and they handle the one place JSON needs a string of its own
// accord.
func (d Decimal) MarshalText() ([]byte, error) {
	return d.AppendText(nil)
}

// AppendText implements encoding.TextAppender, appending the same text
// MarshalText returns to b. Callers encoding many decimals can reuse one
// buffer and avoid the allocation MarshalText makes on every call.
func (d Decimal) AppendText(b []byte) ([]byte, error) {
	return append(b, d.String()...), nil
}

// UnmarshalText implements encoding.TextUnmarshaler, parsing exactly what
// MarshalText writes.
//
// It shares the parser behind NewFromString, so it holds the same line the
// rest of the package holds: exponent notation ("1e2"), empty input,
// surrounding space and more than 19 fractional digits are all rejected. A
// document either round-trips exactly or fails loudly; it never decodes to a
// number that quietly lost digits.
//
// text is only borrowed for the duration of the call; nothing derived from it
// is retained.
func (d *Decimal) UnmarshalText(text []byte) error {
	dec, err := parseDecimal(string(text))
	if err != nil {
		return err
	}

	d.value = dec

	return nil
}

// The binary layout, version 1: a version byte, a header byte holding the sign
// in its top bit and the scale in its low five bits, then the 128-bit
// coefficient as two big-endian uint64s. Eighteen bytes, fixed.
//
// The version byte is what makes this safe to persist. Unlike the text and
// JSON forms, a binary encoding pins the internal representation, and a value
// written today may be read back by a build whose decimal128 no longer looks
// the same. A reader that finds a version it does not know says so instead of
// reading the bytes as something they are not.
const (
	binaryVersion   = 1
	binaryLen       = 18
	binaryNegBit    = 0x80
	binaryScaleMask = 0x1f
)

// MarshalBinary implements encoding.BinaryMarshaler, writing the decimal as
// the fixed 18-byte layout described above.
//
// It exists because a Decimal's only field is unexported, which leaves
// encoding/gob — and every codec that follows the same rule — unable to
// encode one at all: gob consults GobEncoder and BinaryMarshaler, and never
// falls back to MarshalText. It is also what lets a Decimal be handed straight
// to a cache client such as go-redis.
//
// Unlike the text form, this preserves the representation and not just the
// value: 1.50 comes back with the trailing zero it was written with, where
// MarshalText normalises it to 1.5.
//
// It is not the compact choice. A 128-bit coefficient costs sixteen bytes
// whatever digits it holds, so ordinary amounts encode smaller as text — this
// is for reaching binary codecs, not for saving space.
func (d Decimal) MarshalBinary() ([]byte, error) {
	return d.AppendBinary(nil)
}

// AppendBinary implements encoding.BinaryAppender, appending the same bytes
// MarshalBinary returns to b, so a caller encoding many decimals can reuse one
// buffer.
func (d Decimal) AppendBinary(b []byte) ([]byte, error) {
	header := d.value.scale
	if d.value.neg {
		header |= binaryNegBit
	}

	b = append(b, binaryVersion, header)
	b = binary.BigEndian.AppendUint64(b, d.value.coef.hi)

	return binary.BigEndian.AppendUint64(b, d.value.coef.lo), nil
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
//
// Everything the layout does not define is rejected rather than ignored: a
// length other than 18, an unknown version, a scale past the 19 digits this
// package keeps, and any of the header's reserved bits in use. Binary input is
// machine-written, so anything unexpected in it means the reader and the
// writer disagree about the format, which is worth an error rather than a
// silently different number.
//
// data is only read, never retained.
func (d *Decimal) UnmarshalBinary(data []byte) error {
	if len(data) != binaryLen {
		return ErrInvalidBinary
	}

	if data[0] != binaryVersion {
		return ErrUnknownBinaryVersion
	}

	header := data[1]
	if header&^(binaryNegBit|binaryScaleMask) != 0 {
		return ErrInvalidBinary
	}

	scale := header & binaryScaleMask
	if scale > maxScale {
		return ErrPrecOutOfRange
	}

	coef := u128{
		hi: binary.BigEndian.Uint64(data[2:10]),
		lo: binary.BigEndian.Uint64(data[10:18]),
	}

	// Zero has exactly one encoding, because newDec collapses a zero
	// coefficient to the canonical zero whatever sign or scale came with it.
	// A blob carrying either is therefore not one this package wrote, and
	// accepting it would break the property that makes a binary format worth
	// having: equal values encode to equal bytes, so a blob can be compared,
	// deduplicated or used as a cache key.
	if coef.IsZero() && header != 0 {
		return ErrInvalidBinary
	}

	d.value = newDec(header&binaryNegBit != 0, coef, scale)

	return nil
}

func (d Decimal) Value() (driver.Value, error) {
	return d.String(), nil
}

func parseDecimalJSON(data []byte) (decimal128, error) {
	dec := jsontext.NewDecoder(bytes.NewReader(data))

	value, err := dec.ReadValue()
	if err != nil {
		return decimal128{}, err
	}

	// The value aliases the decoder's buffer and the next read overwrites it,
	// so it is parsed before anything else is read rather than cloned.
	parsed, err := parseDecimalJSONValue(value)
	if err != nil {
		return decimal128{}, err
	}

	// UnmarshalJSON is handed raw bytes rather than a positioned decoder, so
	// unlike UnmarshalJSONFrom it has to check that nothing follows the value
	// it read. Without this, "00" decodes as zero: the decoder returns the
	// first number and leaves the second unread. json.Unmarshal never passes
	// such input, but this method is also called directly.
	//
	// The check reads the offset rather than decoding again, so a well-formed
	// document costs nothing beyond a scan of what is left.
	if hasTrailingJSON(data, dec.InputOffset()) {
		return decimal128{}, ErrInvalidFormat
	}

	return parsed, nil
}

// hasTrailingJSON reports whether anything but whitespace follows the value
// that ended at offset.
//
// JSON's whitespace is only space, tab, carriage return and newline.
// bytes.TrimSpace would also swallow a vertical tab or a form feed, which JSON
// does not allow, and accepting them here would make this decoder read
// documents encoding/json/v2 rejects.
func hasTrailingJSON(data []byte, offset int64) bool {
	for _, c := range data[offset:] {
		switch c {
		case ' ', '\t', '\r', '\n':
		default:
			return true
		}
	}

	return false
}

// parseDecimalJSONValue parses the one JSON value a Decimal is written as: a
// bare number, or a string holding the same digits. Both decoding paths go
// through it so the v1 and v2 entry points cannot drift apart.
func parseDecimalJSONValue(value jsontext.Value) (decimal128, error) {
	switch value.Kind() {
	case jsontext.KindNumber:
		return parseDecimal(string(value))
	case jsontext.KindString:
		// Escapes are pointless in a number but legal in JSON, so the string
		// is unquoted properly rather than by trimming the quotes.
		text, err := jsontext.AppendUnquote(nil, value)
		if err != nil {
			return decimal128{}, err
		}

		return parseDecimal(string(text))
	default:
		return decimal128{}, ErrInvalidFormat
	}
}

func (d *Decimal) Scan(src any) error {
	var (
		dec decimal128
		err error
	)

	switch v := src.(type) {
	case []byte:
		dec, err = parseDecimal(string(v))
	case string:
		dec, err = parseDecimal(v)
	case uint64:
		dec, err = decFromUint64(v, 0)
	case int64:
		dec, err = decFromInt64(v, 0)
	case int:
		dec, err = decFromInt64(int64(v), 0)
	case int32:
		dec, err = decFromInt64(int64(v), 0)
	case float64:
		dec, err = decFromFloat64(v)
	case nil:
		err = errors.New("can't scan nil")
	default:
		err = fmt.Errorf("can't scan %T", src)
	}

	if err != nil {
		return err
	}

	d.value = dec

	return nil
}
