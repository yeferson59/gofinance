package money

import (
	"bytes"
	"database/sql/driver"
	"encoding/json/jsontext"

	"github.com/yeferson59/gofinance/v2/decimal"
)

var MoneyZero = Money{decimal.Zero, USD}
var MoneyOne = Money{decimal.One, USD}

type Money struct {
	value    decimal.Decimal
	currency Currency
}

func New(value int64, precision uint8, currency Currency) (Money, error) {
	parsedValue, err := decimal.NewFromInt64(value, precision)
	if err != nil {
		return Money{}, err
	}

	return Money{parsedValue, currency}, nil
}

func NewMoneyFromFloat64(f float64, currency Currency) (Money, error) {
	parsedValue, err := decimal.NewFromFloat64(f)
	if err != nil {
		return Money{}, err
	}

	return Money{parsedValue, currency}, nil
}

func MustMoneyFromFloat64(f float64, currency Currency) Money {
	m, err := NewMoneyFromFloat64(f, currency)
	if err != nil {
		panic(err)
	}

	return m
}

func NewMoneyFromString(s string, currency Currency) (Money, error) {
	parsedValue, err := decimal.NewFromString(s)
	if err != nil {
		return Money{}, err
	}

	return Money{parsedValue, currency}, nil
}

func MustMoneyFromString(s string, currency Currency) Money {
	m, err := NewMoneyFromString(s, currency)
	if err != nil {
		panic(err)
	}

	return m
}

// NewFromDecimal attaches a currency to d, turning it into a Money value.
func NewFromDecimal(d decimal.Decimal, currency Currency) Money {
	return Money{d, currency}
}

func (m Money) GetDecimal() decimal.Decimal {
	return m.value
}

func (m Money) GetCurrency() Currency {
	return m.currency
}

func (m *Money) SetCurrency(c Currency) {
	if m == nil {
		return
	}

	if !c.Valid() {
		m.currency = USD

		return
	}

	m.currency = c
}

// TryAdd returns the sum of m and other, keeping the currency. It returns
// ErrCurrencyMismatch if the operands' currencies differ, or the decimal
// engine's error on overflow.
func (m Money) TryAdd(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	v, err := m.value.TryAdd(other.value)
	if err != nil {
		return Money{}, err
	}

	return Money{v, m.currency}, nil
}

// Add returns the sum of m and other. Like the decimal engine's Add, it
// panics instead of returning an error — on a currency mismatch or on
// overflow. Use TryAdd for the error-returning form.
func (m Money) Add(other Money) Money {
	v, err := m.TryAdd(other)
	if err != nil {
		panic(err)
	}

	return v
}

// TrySub returns the difference of m and other, keeping the currency. It
// returns ErrCurrencyMismatch if the operands' currencies differ, or the
// decimal engine's error on overflow.
func (m Money) TrySub(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	v, err := m.value.TrySub(other.value)
	if err != nil {
		return Money{}, err
	}

	return Money{v, m.currency}, nil
}

// Sub returns the difference of m and other. Like the decimal engine's Sub,
// it panics instead of returning an error — on a currency mismatch or on
// overflow. Use TrySub for the error-returning form.
func (m Money) Sub(other Money) Money {
	v, err := m.TrySub(other)
	if err != nil {
		panic(err)
	}

	return v
}

func (m Money) RoundBank(prec uint8) Money {
	return Money{m.value.RoundBank(prec), m.currency}
}

func (m Money) RoundBankString(prec uint8) string {
	return m.value.RoundBank(prec).StringFixed(prec)
}

func (m Money) RoundAway(prec uint8) Money {
	return Money{m.value.RoundAway(prec), m.currency}
}

func (m Money) Trunc(prec uint8) Money {
	return Money{m.value.Trunc(prec), m.currency}
}

func (m Money) Abs() Money {
	return Money{m.value.Abs(), m.currency}
}

func (m Money) Neg() Money {
	return Money{m.value.Neg(), m.currency}
}

func (m Money) IsZero() bool {
	return m.value.IsZero()
}

// IsPositive reports whether m is strictly greater than zero.
func (m Money) IsPositive() bool {
	return m.value.IsPos()
}

// IsNegative reports whether m is strictly less than zero.
func (m Money) IsNegative() bool {
	return m.value.IsNeg()
}

// MulInt64 multiplies m by a plain integer factor, such as a quantity
// (e.g. unit price * quantity).
func (m Money) MulInt64(n int64) Money {
	factor, err := decimal.NewFromInt64(n, 0)
	if err != nil {
		panic(err)
	}

	return Money{m.value.Mul(factor), m.currency}
}

// MulDecimal multiplies m by a currency-less factor, such as a rate or a
// growth factor, keeping m's currency. Prefer this over Mul when the factor
// is not itself an amount of money.
func (m Money) MulDecimal(d decimal.Decimal) Money {
	return Money{m.value.Mul(d), m.currency}
}

// DivDecimal divides m by a currency-less divisor, such as a rate or a
// number of periods, keeping m's currency.
func (m Money) DivDecimal(d decimal.Decimal) (Money, error) {
	v, err := m.value.Div(d)
	if err != nil {
		return Money{}, err
	}

	return Money{v, m.currency}, nil
}

// MustDivDecimal is like DivDecimal but panics on error.
func (m Money) MustDivDecimal(d decimal.Decimal) Money {
	v, err := m.DivDecimal(d)
	if err != nil {
		panic(err)
	}

	return v
}

// DivInt64 divides m by a plain integer divisor.
func (m Money) DivInt64(n int64) (Money, error) {
	divisor, err := decimal.NewFromInt64(n, 0)
	if err != nil {
		return Money{}, err
	}

	v, err := m.value.Div(divisor)
	if err != nil {
		return Money{}, err
	}

	return Money{v, m.currency}, nil
}

// MustDivInt64 is like DivInt64 but panics on error.
func (m Money) MustDivInt64(n int64) Money {
	v, err := m.DivInt64(n)
	if err != nil {
		panic(err)
	}

	return v
}

// Min returns the smaller of m and other.
// It returns ErrCurrencyMismatch if the operands have different currencies.
func (m Money) Min(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	if m.value.Cmp(other.value) <= 0 {
		return m, nil
	}

	return other, nil
}

// Max returns the larger of m and other.
// It returns ErrCurrencyMismatch if the operands have different currencies.
func (m Money) Max(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}

	if m.value.Cmp(other.value) >= 0 {
		return m, nil
	}

	return other, nil
}

func (m Money) InexactFloat64() float64 {
	return m.value.InexactFloat64()
}

func (m Money) Cmp(other Money) int {
	return m.value.Cmp(other.value)
}

func (m Money) Floor() Money {
	return Money{m.value.Floor(), m.currency}
}

func (m Money) Ceil() Money {
	return Money{m.value.Ceil(), m.currency}
}

func (m Money) String() string {
	return m.value.String()
}

func (m Money) StringFixed(prec uint8) string {
	return m.value.StringFixed(prec)
}

func (m Money) GreaterThan(other Money) bool {
	return m.value.GreaterThan(other.value)
}

func (m Money) GreaterThanOrEqual(other Money) bool {
	return m.value.GreaterThanOrEqual(other.value)
}

func (m Money) LessThan(other Money) bool {
	return m.value.LessThan(other.value)
}

func (m Money) LessThanOrEqual(other Money) bool {
	return m.value.LessThanOrEqual(other.value)
}

func (m Money) Equal(other Money) bool {
	return m.value.Equal(other.value) && m.currency == other.currency
}

// MarshalJSON implements json.Marshaler, writing the amount as an object of
// two members: the value as a JSON string, and the ISO 4217 currency code.
//
// The value is a string rather than a number so it survives readers that parse
// every JSON number as a float64 — JavaScript above all — which is exactly the
// precision this package exists to keep. An unset or unrecognised currency is
// an error: an amount whose currency cannot be named is not an amount.
//
// This is the encoding/json v1 entry point; under encoding/json/v2,
// MarshalJSONTo takes precedence and writes the same object.
func (m Money) MarshalJSON() ([]byte, error) {
	return m.appendJSON(nil)
}

// MarshalJSONTo implements json.MarshalerTo, the streaming encoder
// encoding/json/v2 prefers. It writes into the encoder's own buffer and builds
// the object directly, where MarshalJSON goes through a struct and the
// reflection-based encoder.
func (m Money) MarshalJSONTo(enc *jsontext.Encoder) error {
	buf, err := m.appendJSON(enc.AvailableBuffer())
	if err != nil {
		return err
	}

	return enc.WriteValue(buf)
}

// appendJSON appends the encoded amount to dst. Both encoding paths go through
// it so they cannot drift apart.
//
// Neither member needs JSON escaping: the ISO code is three uppercase letters,
// and the value is digits with an optional minus sign and decimal point.
func (m Money) appendJSON(dst []byte) ([]byte, error) {
	isoCode, err := m.currency.GetCurrencyISOCode()
	if err != nil {
		return dst, err
	}

	dst = append(dst, `{"value":"`...)

	dst, err = m.value.AppendText(dst)
	if err != nil {
		return dst, err
	}

	dst = append(dst, `","currency":"`...)
	dst = append(dst, isoCode[:]...)

	return append(dst, `"}`...), nil
}

// UnmarshalJSON implements json.Unmarshaler.
//
// It accepts the object MarshalJSON writes, and a bare amount — a JSON number
// or a string holding one — which is read as USD, the same default Scan
// applies. Within the object the value may be either form too, the currency
// may be left out or empty for USD, and unknown members are ignored.
//
// This is the encoding/json v1 entry point; under encoding/json/v2,
// UnmarshalJSONFrom takes precedence and accepts the same forms.
func (m *Money) UnmarshalJSON(data []byte) error {
	dec := jsontext.NewDecoder(bytes.NewReader(data))

	parsed, err := parseMoneyJSON(dec)
	if err != nil {
		return err
	}

	// UnmarshalJSON is handed raw bytes rather than a positioned decoder, so
	// unlike UnmarshalJSONFrom it has to check that nothing follows the value
	// it read. Reading the offset costs nothing beyond a scan of what is left,
	// where decoding again would cost a second pass.
	if hasTrailingJSON(data, dec.InputOffset()) {
		return ErrTrailingJSONContent
	}

	*m = parsed

	return nil
}

// UnmarshalJSONFrom implements json.UnmarshalerFrom, the streaming decoder
// encoding/json/v2 prefers. It accepts the same forms as UnmarshalJSON.
//
// It consumes exactly one JSON value, including one it goes on to reject, as
// the interface requires.
func (m *Money) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	parsed, err := parseMoneyJSON(dec)
	if err != nil {
		return err
	}

	*m = parsed

	return nil
}

// parseMoneyJSON reads one amount from dec. Both decoding paths go through it
// so they cannot drift apart.
func parseMoneyJSON(dec *jsontext.Decoder) (Money, error) {
	if dec.PeekKind() == jsontext.KindBeginObject {
		return parseMoneyJSONObject(dec)
	}

	// A bare amount carries no currency, so it takes the default. Anything
	// that is not a number or a string holding one is rejected by the decimal
	// decoder, which consumes the value either way.
	var value decimal.Decimal
	if err := value.UnmarshalJSONFrom(dec); err != nil {
		return Money{}, err
	}

	return Money{value: value, currency: USD}, nil
}

// parseMoneyJSONObject reads the object form member by member rather than into
// a struct, so the amount and the currency are decoded by the same code that
// decodes them on their own.
func parseMoneyJSONObject(dec *jsontext.Decoder) (Money, error) {
	if _, err := dec.ReadToken(); err != nil { // '{'
		return Money{}, err
	}

	var (
		value     decimal.Decimal
		currency  = USD
		seenValue bool
	)

	for dec.PeekKind() != jsontext.KindEndObject {
		name, err := dec.ReadToken()
		if err != nil {
			return Money{}, err
		}

		switch name.String() {
		case "value":
			if err := value.UnmarshalJSONFrom(dec); err != nil {
				return Money{}, err
			}

			seenValue = true
		case "currency":
			currency, err = parseMoneyJSONCurrency(dec)
			if err != nil {
				return Money{}, err
			}
		default:
			// Unknown members are ignored, as the struct-based decoder did.
			if err := dec.SkipValue(); err != nil {
				return Money{}, err
			}
		}
	}

	if _, err := dec.ReadToken(); err != nil { // '}'
		return Money{}, err
	}

	if !seenValue {
		return Money{}, ErrMissingJSONValue
	}

	return Money{value: value, currency: currency}, nil
}

// parseMoneyJSONCurrency reads the "currency" member, where an empty code
// means "unset" and is read as USD rather than rejected: a document that
// leaves the currency blank describes dollars the same way one that omits the
// member does. Currency's own decoder has no such amount to fall back on, so
// it rejects the empty string.
func parseMoneyJSONCurrency(dec *jsontext.Decoder) (Currency, error) {
	if dec.PeekKind() != jsontext.KindString {
		return parseCurrencyJSON(dec)
	}

	token, err := dec.ReadToken()
	if err != nil {
		return 0, err
	}

	code := token.String()
	if code == "" {
		return USD, nil
	}

	return GetCurrencyFromISOCode(code)
}

// The binary layout, version 1: a version byte, the three-letter ISO 4217
// code, then the amount in the decimal package's own binary form. Twenty-two
// bytes with that package's current layout.
//
// The currency is written before the amount so both sit at fixed offsets, and
// the amount is delegated rather than reproduced: its layout is the decimal
// package's to define and to version, and it carries its own version byte for
// exactly that reason.
const (
	moneyBinaryVersion = 1

	// moneyBinaryPrefixLen covers the version byte and the ISO code. What
	// follows is the decimal's own encoding, whose length is its business.
	moneyBinaryPrefixLen = 4
)

// MarshalBinary implements encoding.BinaryMarshaler.
//
// It exists because a Money's fields are unexported, which leaves
// encoding/gob — and every codec that follows the same rule — unable to encode
// one at all: gob consults GobEncoder and BinaryMarshaler, and never falls
// back to MarshalText. It is also what lets an amount be handed straight to a
// cache client such as go-redis.
//
// An amount whose currency cannot be named is an error, as it is for the JSON
// encoders.
func (m Money) MarshalBinary() ([]byte, error) {
	return m.AppendBinary(nil)
}

// AppendBinary implements encoding.BinaryAppender, appending the same bytes
// MarshalBinary returns to b. On error b is returned unchanged, so a failure
// never truncates what the caller had already built.
func (m Money) AppendBinary(b []byte) ([]byte, error) {
	isoCode, err := m.currency.GetCurrencyISOCode()
	if err != nil {
		return b, err
	}

	b = append(b, moneyBinaryVersion)
	b = append(b, isoCode[:]...)

	return m.value.AppendBinary(b)
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
//
// Unlike the JSON decoder, it has no bare-amount form to fall back on: these
// bytes were written by this package, and one that carries no currency is not
// an amount this package wrote. data is only read, never retained.
func (m *Money) UnmarshalBinary(data []byte) error {
	if len(data) < moneyBinaryPrefixLen {
		return ErrInvalidBinary
	}

	if data[0] != moneyBinaryVersion {
		return ErrUnknownBinaryVersion
	}

	currency, err := getCurrencyByISOCode([3]byte(data[1:moneyBinaryPrefixLen]))
	if err != nil {
		return err
	}

	// The decimal validates its own length, so a truncated or padded amount is
	// caught there rather than by a length check that would have to know its
	// layout.
	var value decimal.Decimal
	if err := value.UnmarshalBinary(data[moneyBinaryPrefixLen:]); err != nil {
		return err
	}

	m.value = value
	m.currency = currency

	return nil
}

func (m Money) Value() (driver.Value, error) {
	return m.value.String(), nil
}

func (m *Money) Scan(src any) error {
	if err := m.value.Scan(src); err != nil {
		return err
	}

	m.SetCurrency(USD)

	return nil
}
