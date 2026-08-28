package money

import (
	"errors"
	"strconv"
	"strings"
)

type Currency uint8

const (
	XXX Currency = 0   // No Currency
	XTS Currency = 1   // Test Currency
	AED Currency = 2   // U.A.E. Dirham
	AFN Currency = 3   // Afghani
	ALL Currency = 4   // Lek
	AMD Currency = 5   // Armenian Dram
	ANG Currency = 6   // Netherlands Antillian Guilder
	AOA Currency = 7   // Kwanza
	ARS Currency = 8   // Argentine Peso
	AUD Currency = 9   // Australian Dollar
	AWG Currency = 10  // Aruban Guilder
	AZN Currency = 11  // Azerbaijan Manat
	BAM Currency = 12  // Convertible Mark
	BBD Currency = 13  // Barbados Dollar
	BDT Currency = 14  // Taka
	BGN Currency = 15  // Bulgarian Lev
	BHD Currency = 16  // Bahraini Dinar
	BIF Currency = 17  // Burundi Franc
	BMD Currency = 18  // Bermudian Dollar
	BND Currency = 19  // Brunei Dollar
	BOB Currency = 20  // Boliviano
	BRL Currency = 21  // Brazilian Real
	BSD Currency = 22  // Bahamian Dollar
	BTN Currency = 23  // Bhutan Ngultrum
	BWP Currency = 24  // Pula
	BYN Currency = 25  // Belarussian Ruble
	BZD Currency = 26  // Belize Dollar
	CAD Currency = 27  // Canadian Dollar
	CDF Currency = 28  // Franc Congolais
	CHF Currency = 29  // Swiss Franc
	CLP Currency = 30  // Chilean Peso
	CNY Currency = 31  // Yuan Renminbi
	COP Currency = 32  // Colombian Peso
	CRC Currency = 33  // Costa Rican Colon
	CUP Currency = 34  // Cuban Peso
	CVE Currency = 35  // Cape Verde Escudo
	CZK Currency = 36  // Czech Koruna
	DJF Currency = 37  // Djibouti Franc
	DKK Currency = 38  // Danish Krone
	DOP Currency = 39  // Dominican Peso
	DZD Currency = 40  // Algerian Dinar
	EGP Currency = 41  // Egyptian Pound
	ERN Currency = 42  // Eritean Nakfa
	ETB Currency = 43  // Ethiopian Birr
	EUR Currency = 44  // Euro
	FJD Currency = 45  // Fiji Dollar
	FKP Currency = 46  // Falkland Islands Pound
	GBP Currency = 47  // Pound Sterling
	GEL Currency = 48  // Lari
	GHS Currency = 49  // Cedi
	GIP Currency = 50  // Gibraltar Pound
	GMD Currency = 51  // Dalasi
	GNF Currency = 52  // Guinea Franc
	GTQ Currency = 53  // Quetzal
	GWP Currency = 54  // Guinea-Bissau Peso
	GYD Currency = 55  // Guyana Dollar
	HKD Currency = 56  // Hong Kong Dollar
	HNL Currency = 57  // Lempira
	HRK Currency = 58  // Croatian Kuna
	HTG Currency = 59  // Gourde
	HUF Currency = 60  // Forint
	IDR Currency = 61  // Rupiah
	ILS Currency = 62  // Israeli Shequel
	INR Currency = 63  // Indian Rupee
	IQD Currency = 64  // Iraqi Dinar
	IRR Currency = 65  // Iranian Rial
	ISK Currency = 66  // Iceland Krona
	JMD Currency = 67  // Jamaican Dollar
	JOD Currency = 68  // Jordanian Dinar
	JPY Currency = 69  // Yen
	KES Currency = 70  // Kenyan Shilling
	KGS Currency = 71  // Som
	KHR Currency = 72  // Riel
	KMF Currency = 73  // Comoro Franc
	KPW Currency = 74  // North Korean Won
	KRW Currency = 75  // Won
	KWD Currency = 76  // Kuwaiti Dinar
	KYD Currency = 77  // Cayman Islands Dollar
	KZT Currency = 78  // Tenge
	LAK Currency = 79  // Kip
	LBP Currency = 80  // Lebanese Pound
	LKR Currency = 81  // Sri Lanka Rupee
	LRD Currency = 82  // Liberian Dollar
	LSL Currency = 83  // Lesotho Loti
	LYD Currency = 84  // Libyan Dinar
	MAD Currency = 85  // Moroccan Dirham
	MDL Currency = 86  // Moldovan Leu
	MGA Currency = 87  // Malagasy Ariary
	MKD Currency = 88  // Denar
	MMK Currency = 89  // Kyat
	MNT Currency = 90  // Tugrik
	MOP Currency = 91  // Pataca
	MRU Currency = 92  // Ouguiya
	MUR Currency = 93  // Mauritius Rupee
	MVR Currency = 94  // Rufiyaa
	MWK Currency = 95  // Malawi Kwacha
	MXN Currency = 96  // Mexican Peso
	MYR Currency = 97  // Malaysian Ringgit
	MZN Currency = 98  // Mozambique Metical
	NAD Currency = 99  // Namibia Dollar
	NGN Currency = 100 // Naira
	NIO Currency = 101 // Cordoba Oro
	NOK Currency = 102 // Norwegian Krone
	NPR Currency = 103 // Nepalese Rupee
	NZD Currency = 104 // New Zealand Dollar
	OMR Currency = 105 // Rial Omani
	PAB Currency = 106 // Balboa
	PEN Currency = 107 // Sol
	PGK Currency = 108 // Kina
	PHP Currency = 109 // Philippine Peso
	PKR Currency = 110 // Pakistan Rupee
	PLN Currency = 111 // Zloty
	PYG Currency = 112 // Guarani
	QAR Currency = 113 // Qatari Rial
	RON Currency = 114 // Leu
	RSD Currency = 115 // Serbian Dinar
	RUB Currency = 116 // Russian Ruble
	RWF Currency = 117 // Rwanda Franc
	SAR Currency = 118 // Saudi Riyal
	SBD Currency = 119 // Solomon Islands Dollar
	SCR Currency = 120 // Seychelles Rupee
	SDG Currency = 121 // Sudanese Pound
	SEK Currency = 122 // Swedish Krona
	SGD Currency = 123 // Singapore Dollar
	SHP Currency = 124 // St. Helena Pound
	SLL Currency = 125 // Leone
	SOS Currency = 126 // Somali Shilling
	SRD Currency = 127 // Surinam Dollar
	SSP Currency = 128 // South Sudanese Pound
	STN Currency = 129 // Dobra
	SYP Currency = 130 // Syrian Pound
	SZL Currency = 131 // Lilangeni
	THB Currency = 132 // Baht
	TJS Currency = 133 // Somoni
	TMT Currency = 134 // Manat
	TND Currency = 135 // Tunisian Dinar
	TOP Currency = 136 // Pa'anga
	TRY Currency = 137 // Turkish Lira
	TTD Currency = 138 // Trinidad and Tobago Dollar
	TWD Currency = 139 // New Taiwan Dollar
	TZS Currency = 140 // Tanzanian Shilling
	UAH Currency = 141 // Ukrainian Hryvnia
	UGX Currency = 142 // Uganda Shilling
	USD Currency = 143 // U.S. Dollar
	UYU Currency = 144 // Peso Uruguayo
	UZS Currency = 145 // Uzbekistan Sum
	VES Currency = 146 // Sovereign Bolivar
	VND Currency = 147 // Dong
	VUV Currency = 148 // Vatu
	WST Currency = 149 // Tala
	XAF Currency = 150 // CFA Franc BEAC
	XCD Currency = 151 // East Caribbean Dollar
	XOF Currency = 152 // CFA Franc BCEAO
	XPF Currency = 153 // CFP Franc
	YER Currency = 154 // Yemeni Rial
	ZAR Currency = 155 // Rand
	ZMW Currency = 156 // Zambian Kwacha
	ZWL Currency = 157 // Zimbabwe Dollar
)

var currencyCode = map[Currency][3]byte{
	XXX: {'X', 'X', 'X'},
	XTS: {'X', 'T', 'S'},
	AED: {'A', 'E', 'D'},
	AFN: {'A', 'F', 'N'},
	ALL: {'A', 'L', 'L'},
	AMD: {'A', 'M', 'D'},
	ANG: {'A', 'N', 'G'},
	AOA: {'A', 'O', 'A'},
	ARS: {'A', 'R', 'S'},
	AUD: {'A', 'U', 'D'},
	AWG: {'A', 'W', 'G'},
	AZN: {'A', 'Z', 'N'},
	BAM: {'B', 'A', 'M'},
	BBD: {'B', 'B', 'D'},
	BDT: {'B', 'D', 'T'},
	BGN: {'B', 'G', 'N'},
	BHD: {'B', 'H', 'D'},
	BIF: {'B', 'I', 'F'},
	BMD: {'B', 'M', 'D'},
	BND: {'B', 'N', 'D'},
	BOB: {'B', 'O', 'B'},
	BRL: {'B', 'R', 'L'},
	COP: {'C', 'O', 'P'},
	USD: {'U', 'S', 'D'},
	VND: {'V', 'N', 'D'},
	VUV: {'V', 'U', 'V'},
	WST: {'W', 'S', 'T'},
	XAF: {'X', 'A', 'F'},
	XCD: {'X', 'C', 'D'},
	XOF: {'X', 'O', 'F'},
	XPF: {'X', 'P', 'F'},
	YER: {'Y', 'E', 'R'},
	ZAR: {'Z', 'A', 'R'},
	ZMW: {'Z', 'M', 'W'},
	ZWL: {'Z', 'W', 'L'},
	BSD: {'B', 'S', 'D'},
	BTN: {'B', 'T', 'N'},
	BWP: {'B', 'W', 'P'},
	BYN: {'B', 'Y', 'N'},
	BZD: {'B', 'Z', 'D'},
	CAD: {'C', 'A', 'D'},
	CDF: {'C', 'D', 'F'},
	CHF: {'C', 'H', 'F'},
	CLP: {'C', 'L', 'P'},
	CNY: {'C', 'N', 'Y'},
	UYU: {'U', 'Y', 'U'},
	UZS: {'U', 'Z', 'S'},
	VES: {'V', 'E', 'S'},
	CRC: {'C', 'R', 'C'},
	CUP: {'C', 'U', 'P'},
	CVE: {'C', 'V', 'E'},
	CZK: {'C', 'Z', 'K'},
	DJF: {'D', 'J', 'F'},
	DKK: {'D', 'K', 'K'},
	DOP: {'D', 'O', 'P'},
	DZD: {'D', 'Z', 'D'},
	EGP: {'E', 'G', 'P'},
	ERN: {'E', 'R', 'N'},
	ETB: {'E', 'T', 'B'},
	EUR: {'E', 'U', 'R'},
	FJD: {'F', 'J', 'D'},
	FKP: {'F', 'K', 'P'},
	GBP: {'G', 'B', 'P'},
	GEL: {'G', 'E', 'L'},
	GHS: {'G', 'H', 'S'},
	GIP: {'G', 'I', 'P'},
	GMD: {'G', 'M', 'D'},
	GNF: {'G', 'N', 'F'},
	GTQ: {'G', 'T', 'Q'},
	GWP: {'G', 'W', 'P'},
	GYD: {'G', 'Y', 'D'},
	HKD: {'H', 'K', 'D'},
	HNL: {'H', 'N', 'L'},
	HRK: {'H', 'R', 'K'},
	HTG: {'H', 'T', 'G'},
	HUF: {'H', 'U', 'F'},
	IDR: {'I', 'D', 'R'},
	ILS: {'I', 'L', 'S'},
	INR: {'I', 'N', 'R'},
	IQD: {'I', 'Q', 'D'},
	IRR: {'I', 'R', 'R'},
	ISK: {'I', 'S', 'K'},
	JMD: {'J', 'M', 'D'},
	JOD: {'J', 'O', 'D'},
	JPY: {'J', 'P', 'Y'},
	KES: {'K', 'E', 'S'},
	KGS: {'K', 'G', 'S'},
	KHR: {'K', 'H', 'R'},
	KMF: {'K', 'M', 'F'},
	KPW: {'K', 'P', 'W'},
	KRW: {'K', 'R', 'W'},
	KWD: {'K', 'W', 'D'},
	KYD: {'K', 'Y', 'D'},
	KZT: {'K', 'Z', 'T'},
	LAK: {'L', 'A', 'K'},
	LBP: {'L', 'B', 'P'},
	LKR: {'L', 'K', 'R'},
	LRD: {'L', 'R', 'D'},
	LSL: {'L', 'S', 'L'},
	LYD: {'L', 'Y', 'D'},
	MAD: {'M', 'A', 'D'},
	MDL: {'M', 'D', 'L'},
	MGA: {'M', 'G', 'A'},
	MKD: {'M', 'K', 'D'},
	MMK: {'M', 'M', 'K'},
	MNT: {'M', 'N', 'T'},
	MOP: {'M', 'O', 'P'},
	MRU: {'M', 'R', 'U'},
	MUR: {'M', 'U', 'R'},
	MVR: {'M', 'V', 'R'},
	MWK: {'M', 'W', 'K'},
	MXN: {'M', 'X', 'N'},
	MYR: {'M', 'Y', 'R'},
	MZN: {'M', 'Z', 'N'},
	NAD: {'N', 'A', 'D'},
	NGN: {'N', 'G', 'N'},
	NIO: {'N', 'I', 'O'},
	NOK: {'N', 'O', 'K'},
	NPR: {'N', 'P', 'R'},
	NZD: {'N', 'Z', 'D'},
	OMR: {'O', 'M', 'R'},
	PAB: {'P', 'A', 'B'},
	PEN: {'P', 'E', 'N'},
	PGK: {'P', 'G', 'K'},
	PHP: {'P', 'H', 'P'},
	PKR: {'P', 'K', 'R'},
	PLN: {'P', 'L', 'N'},
	PYG: {'P', 'Y', 'G'},
	QAR: {'Q', 'A', 'R'},
	RON: {'R', 'O', 'N'},
	RSD: {'R', 'S', 'D'},
	RUB: {'R', 'U', 'B'},
	RWF: {'R', 'W', 'F'},
	SAR: {'S', 'A', 'R'},
	SBD: {'S', 'B', 'D'},
	SCR: {'S', 'C', 'R'},
	SDG: {'S', 'D', 'G'},
	SEK: {'S', 'E', 'K'},
	SGD: {'S', 'G', 'D'},
	SHP: {'S', 'H', 'P'},
	SLL: {'S', 'L', 'L'},
	SOS: {'S', 'O', 'S'},
	SRD: {'S', 'R', 'D'},
	SSP: {'S', 'S', 'P'},
	STN: {'S', 'T', 'N'},
	SYP: {'S', 'Y', 'P'},
	SZL: {'S', 'Z', 'L'},
	THB: {'T', 'H', 'B'},
	TJS: {'T', 'J', 'S'},
	TMT: {'T', 'M', 'T'},
	TND: {'T', 'N', 'D'},
	TOP: {'T', 'O', 'P'},
	TRY: {'T', 'R', 'Y'},
	TTD: {'T', 'T', 'D'},
	TWD: {'T', 'W', 'D'},
	TZS: {'T', 'Z', 'S'},
	UAH: {'U', 'A', 'H'},
	UGX: {'U', 'G', 'X'},
}

var currencyByISOCode = func() map[[3]byte]Currency {
	mapped := make(map[[3]byte]Currency, len(currencyCode))

	for currency, code := range currencyCode {
		mapped[code] = currency
	}

	return mapped
}()

// defaultCurrencyPrec is the ISO 4217 minor unit (decimal places) used for
// currencies without an explicit entry in currencyPrec.
const defaultCurrencyPrec uint8 = 2

// currencyPrec overrides the default precision (2 decimal places) for
// currencies whose ISO 4217 minor unit differs.
var currencyPrec = map[Currency]uint8{
	// Zero decimal places
	BIF: 0,
	CLP: 0,
	DJF: 0,
	GNF: 0,
	ISK: 0,
	JPY: 0,
	KMF: 0,
	KRW: 0,
	PYG: 0,
	RWF: 0,
	UGX: 0,
	VND: 0,
	VUV: 0,
	XAF: 0,
	XOF: 0,
	XPF: 0,
	// Three decimal places
	BHD: 3,
	IQD: 3,
	JOD: 3,
	KWD: 3,
	LYD: 3,
	OMR: 3,
	TND: 3,
}

// currencySymbol maps a currency to its common display symbol (e.g. "$"
// for USD). Currencies without a distinct, unambiguous symbol are omitted;
// Symbol falls back to the ISO 4217 code for those.
var currencySymbol = map[Currency]string{
	AED: "د.إ",
	ARS: "$",
	AUD: "A$",
	BDT: "৳",
	BRL: "R$",
	CAD: "C$",
	CHF: "CHF",
	CLP: "$",
	CNY: "¥",
	COP: "$",
	EGP: "E£",
	EUR: "€",
	GBP: "£",
	HKD: "HK$",
	IDR: "Rp",
	ILS: "₪",
	INR: "₹",
	JPY: "¥",
	KRW: "₩",
	MXN: "$",
	MYR: "RM",
	NGN: "₦",
	NOK: "kr",
	NZD: "NZ$",
	PEN: "S/",
	PHP: "₱",
	PKR: "₨",
	PLN: "zł",
	RUB: "₽",
	SAR: "﷼",
	SEK: "kr",
	SGD: "S$",
	THB: "฿",
	TRY: "₺",
	UAH: "₴",
	USD: "$",
	VND: "₫",
	ZAR: "R",
}

// Symbol returns the common display symbol for c (e.g. "$" for USD, "€"
// for EUR). If c has no distinct symbol on record, Symbol falls back to
// its ISO 4217 code.
func (c Currency) Symbol() (string, error) {
	isoCode, err := c.GetCurrencyISOCode()
	if err != nil {
		return "", err
	}

	if symbol, ok := currencySymbol[c]; ok {
		return symbol, nil
	}

	return string(isoCode[:]), nil
}

// ResolveCurrency returns the single currency the given amounts are expressed
// in, ignoring any that were never set.
//
// It exists for the types that carry several optional amounts — a simple or
// compound interest configuration, an annuity — where a caller supplies only
// the ones the calculation needs. An unset Money is the zero value, which
// carries XXX, the ISO 4217 code for "no currency". Deriving a result's
// currency from one particular field therefore yields XXX whenever that field
// is the unset one, and combining such a result with a real amount fails with
// a currency mismatch. Resolving across every field instead keeps a partially
// configured value in one currency.
//
// It returns ErrCurrencyMismatch when two amounts that are set disagree, and
// XXX with no error when none of them is set.
func ResolveCurrency(amounts ...Money) (Currency, error) {
	resolved := XXX

	for _, amount := range amounts {
		currency := amount.GetCurrency()
		if currency == XXX {
			continue
		}

		if resolved != XXX && currency != resolved {
			return XXX, ErrCurrencyMismatch
		}

		resolved = currency
	}

	return resolved, nil
}

// String returns c's ISO 4217 code, so a Currency prints as "USD" rather than
// as the number behind it. An unrecognised currency prints as
// "Currency(<n>)", which names the problem instead of hiding it behind a bare
// integer.
//
// It exists because Currency is an integer type: without String, printing one
// with %v or passing it to Println produced a meaningless number.
func (c Currency) String() string {
	if isoCode, ok := currencyCode[c]; ok {
		return string(isoCode[:])
	}

	return "Currency(" + strconv.FormatUint(uint64(c), 10) + ")"
}

func (c Currency) GetCurrencyISOCode() ([3]byte, error) {
	isoCode, ok := currencyCode[c]
	if !ok {
		return [3]byte{}, errors.New("invalid currency ISO code")
	}

	return isoCode, nil
}

func (c Currency) GetCurrencyPrecisionCode() (uint8, error) {
	if _, ok := currencyCode[c]; !ok {
		return 0, errors.New("invalid currency")
	}

	if prec, ok := currencyPrec[c]; ok {
		return prec, nil
	}

	return defaultCurrencyPrec, nil
}

func GetCurrencyFromISOCode(code string) (Currency, error) {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	if normalized == "" {
		return 0, errors.New("invalid currency ISO code")
	}

	newCode := make([]byte, 0, 3)

	for i := range 3 {
		newCode = append(newCode, normalized[i])
	}

	currency, ok := currencyByISOCode[[3]byte(newCode)]
	if !ok {
		return 0, errors.New("invalid currency ISO code")
	}

	return currency, nil
}
