package money

import (
	"errors"
	"strconv"
	"strings"
)

type Currency uint

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

var currencyCode = map[Currency]string{
	XXX: "XXX",
	XTS: "XTS",
	AED: "AED",
	AFN: "AFN",
	ALL: "ALL",
	AMD: "AMD",
	ANG: "ANG",
	AOA: "AOA",
	ARS: "ARS",
	AUD: "AUD",
	AWG: "AWG",
	AZN: "AZN",
	BAM: "BAM",
	BBD: "BBD",
	BDT: "BDT",
	BGN: "BGN",
	BHD: "BHD",
	BIF: "BIF",
	BMD: "BMD",
	BND: "BND",
	BOB: "BOB",
	BRL: "BRL",
	BSD: "BSD",
	BTN: "BTN",
	BWP: "BWP",
	BYN: "BYN",
	BZD: "BZD",
	CAD: "CAD",
	CDF: "CDF",
	CHF: "CHF",
	CLP: "CLP",
	CNY: "CNY",
	COP: "COP",
	CRC: "CRC",
	CUP: "CUP",
	CVE: "CVE",
	CZK: "CZK",
	DJF: "DJF",
	DKK: "DKK",
	DOP: "DOP",
	DZD: "DZD",
	EGP: "EGP",
	ERN: "ERN",
	ETB: "ETB",
	EUR: "EUR",
	FJD: "FJD",
	FKP: "FKP",
	GBP: "GBP",
	GEL: "GEL",
	GHS: "GHS",
	GIP: "GIP",
	GMD: "GMD",
	GNF: "GNF",
	GTQ: "GTQ",
	GWP: "GWP",
	GYD: "GYD",
	HKD: "HKD",
	HNL: "HNL",
	HRK: "HRK",
	HTG: "HTG",
	HUF: "HUF",
	IDR: "IDR",
	ILS: "ILS",
	INR: "INR",
	IQD: "IQD",
	IRR: "IRR",
	ISK: "ISK",
	JMD: "JMD",
	JOD: "JOD",
	JPY: "JPY",
	KES: "KES",
	KGS: "KGS",
	KHR: "KHR",
	KMF: "KMF",
	KPW: "KPW",
	KRW: "KRW",
	KWD: "KWD",
	KYD: "KYD",
	KZT: "KZT",
	LAK: "LAK",
	LBP: "LBP",
	LKR: "LKR",
	LRD: "LRD",
	LSL: "LSL",
	LYD: "LYD",
	MAD: "MAD",
	MDL: "MDL",
	MGA: "MGA",
	MKD: "MKD",
	MMK: "MMK",
	MNT: "MNT",
	MOP: "MOP",
	MRU: "MRU",
	MUR: "MUR",
	MVR: "MVR",
	MWK: "MWK",
	MXN: "MXN",
	MYR: "MYR",
	MZN: "MZN",
	NAD: "NAD",
	NGN: "NGN",
	NIO: "NIO",
	NOK: "NOK",
	NPR: "NPR",
	NZD: "NZD",
	OMR: "OMR",
	PAB: "PAB",
	PEN: "PEN",
	PGK: "PGK",
	PHP: "PHP",
	PKR: "PKR",
	PLN: "PLN",
	PYG: "PYG",
	QAR: "QAR",
	RON: "RON",
	RSD: "RSD",
	RUB: "RUB",
	RWF: "RWF",
	SAR: "SAR",
	SBD: "SBD",
	SCR: "SCR",
	SDG: "SDG",
	SEK: "SEK",
	SGD: "SGD",
	SHP: "SHP",
	SLL: "SLL",
	SOS: "SOS",
	SRD: "SRD",
	SSP: "SSP",
	STN: "STN",
	SYP: "SYP",
	SZL: "SZL",
	THB: "THB",
	TJS: "TJS",
	TMT: "TMT",
	TND: "TND",
	TOP: "TOP",
	TRY: "TRY",
	TTD: "TTD",
	TWD: "TWD",
	TZS: "TZS",
	UAH: "UAH",
	UGX: "UGX",
	USD: "USD",
	UYU: "UYU",
	UZS: "UZS",
	VES: "VES",
	VND: "VND",
	VUV: "VUV",
	WST: "WST",
	XAF: "XAF",
	XCD: "XCD",
	XOF: "XOF",
	XPF: "XPF",
	YER: "YER",
	ZAR: "ZAR",
	ZMW: "ZMW",
	ZWL: "ZWL",
}

var currencyByISOCode = func() map[string]Currency {
	mapped := make(map[string]Currency, len(currencyCode))
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

	return isoCode, nil
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
		return isoCode
	}

	return "Currency(" + strconv.FormatUint(uint64(c), 10) + ")"
}

func (c Currency) GetCurrencyISOCode() (string, error) {
	isoCode, ok := currencyCode[c]
	if !ok {
		return "", errors.New("invalid currency ISO code")
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

func CurrencyFromISOCode(code string) (Currency, error) {
	normalized := strings.ToUpper(strings.TrimSpace(code))
	if normalized == "" {
		return 0, errors.New("invalid currency ISO code")
	}

	currency, ok := currencyByISOCode[normalized]
	if !ok {
		return 0, errors.New("invalid currency ISO code")
	}

	return currency, nil
}
