package types

var (
	// Currencies contains the in-memory database of supported currencies
	Currencies map[string]ConcreteCurrency
)

func init() {
	Currencies = make(map[string]ConcreteCurrency)
	currencyData := []struct {
		d uint
		m int64
		s string
	}{
		{2, 1, "AED"},
		{2, 1, "AFN"},
		{2, 1, "ALL"},
		{2, 1, "AMD"},
		{2, 1, "ANG"},
		{2, 1, "AOA"},
		{2, 1, "ARS"},
		{2, 1, "AUD"},
		{2, 1, "AWG"},
		{2, 1, "AZN"},
		{2, 1, "BAM"},
		{2, 1, "BBD"},
		{2, 1, "BDT"},
		{2, 1, "BGN"},
		{3, 1, "BHD"},
		{0, 1, "BIF"},
		{2, 1, "BMD"},
		{2, 1, "BND"},
		{2, 1, "BOB"},
		{2, 1, "BOV"},
		{2, 1, "BRL"},
		{2, 1, "BSD"},
		{2, 1, "BTN"},
		{2, 1, "BWP"},
		{2, 1, "BYN"},
		{0, 1, "BYR"},
		{2, 1, "BZD"},
		{2, 1, "CAD"},
		{2, 1, "CDF"},
		{2, 1, "CHE"},
		{2, 1, "CHF"},
		{2, 1, "CHW"},
		{4, 1, "CLF"},
		{0, 1, "CLP"},
		{2, 1, "CNY"},
		{2, 1, "COP"},
		{2, 1, "COU"},
		{2, 1, "CRC"},
		{2, 1, "CUC"},
		{2, 1, "CUP"},
		{0, 1, "CVE"},
		{2, 1, "CZK"},
		{0, 1, "DJF"},
		{2, 1, "DKK"},
		{2, 1, "DOP"},
		{2, 1, "DZD"},
		{2, 1, "EGP"},
		{2, 1, "ERN"},
		{2, 1, "ETB"},
		{2, 1, "EUR"},
		{2, 1, "FJD"},
		{2, 1, "FKP"},
		{2, 1, "GBP"},
		{2, 1, "GEL"},
		{2, 1, "GHS"},
		{2, 1, "GIP"},
		{2, 1, "GMD"},
		{0, 1, "GNF"},
		{2, 1, "GTQ"},
		{2, 1, "GYD"},
		{2, 1, "HKD"},
		{2, 1, "HNL"},
		{2, 1, "HRK"},
		{2, 1, "HTG"},
		{2, 1, "HUF"},
		{2, 1, "IDR"},
		{2, 1, "ILS"},
		{2, 1, "INR"},
		{3, 1, "IQD"},
		{2, 1, "IRR"},
		{0, 1, "ISK"},
		{2, 1, "JMD"},
		{3, 1, "JOD"},
		{0, 1, "JPY"},
		{2, 1, "KES"},
		{2, 1, "KGS"},
		{2, 1, "KHR"},
		{0, 1, "KMF"},
		{2, 1, "KPW"},
		{0, 1, "KRW"},
		{3, 1, "KWD"},
		{2, 1, "KYD"},
		{2, 1, "KZT"},
		{2, 1, "LAK"},
		{2, 1, "LBP"},
		{2, 1, "LKR"},
		{2, 1, "LRD"},
		{2, 1, "LSL"},
		{3, 1, "LYD"},
		{2, 1, "MAD"},
		{2, 1, "MDL"},
		{1, 1, "MGA"},
		{2, 1, "MKD"},
		{2, 1, "MMK"},
		{2, 1, "MNT"},
		{2, 1, "MOP"},
		{1, 1, "MRO"},
		{2, 1, "MUR"},
		{2, 1, "MVR"},
		{2, 1, "MWK"},
		{2, 1, "MXN"},
		{2, 1, "MXV"},
		{2, 1, "MYR"},
		{2, 1, "MZN"},
		{2, 1, "NAD"},
		{2, 1, "NGN"},
		{2, 1, "NIO"},
		{2, 1, "NOK"},
		{2, 1, "NPR"},
		{2, 1, "NZD"},
		{3, 1, "OMR"},
		{2, 1, "PAB"},
		{2, 1, "PEN"},
		{2, 1, "PGK"},
		{2, 1, "PHP"},
		{2, 1, "PKR"},
		{2, 1, "PLN"},
		{0, 1, "PYG"},
		{2, 1, "QAR"},
		{2, 1, "RON"},
		{2, 1, "RSD"},
		{2, 1, "RUB"},
		{0, 1, "RWF"},
		{2, 1, "SAR"},
		{2, 1, "SBD"},
		{2, 1, "SCR"},
		{2, 1, "SDG"},
		{2, 1, "SEK"},
		{2, 1, "SGD"},
		{2, 1, "SHP"},
		{2, 1, "SLL"},
		{2, 1, "SOS"},
		{2, 1, "SRD"},
		{2, 1, "SSP"},
		{2, 1, "STD"},
		{2, 1, "SVC"},
		{2, 1, "SYP"},
		{2, 1, "SZL"},
		{2, 1, "THB"},
		{2, 1, "TJS"},
		{2, 1, "TMT"},
		{3, 1, "TND"},
		{2, 1, "TOP"},
		{2, 1, "TRY"},
		{2, 1, "TTD"},
		{2, 1, "TWD"},
		{2, 1, "TZS"},
		{2, 1, "UAH"},
		{0, 1, "UGX"},
		{2, 1, "USD"},
		{2, 1, "USN"},
		{0, 1, "UYI"},
		{2, 1, "UYU"},
		{2, 1, "UZS"},
		{2, 1, "VEF"},
		{0, 1, "VND"},
		{0, 1, "VUV"},
		{2, 1, "WST"},
		{0, 1, "XAF"},
		{2, 1, "XCD"},
		{0, 1, "XOF"},
		{0, 1, "XPF"},
		{2, 1, "YER"},
		{2, 1, "ZAR"},
		{2, 1, "ZMW"},
		{2, 1, "ZWL"},
	}
	for _, d := range currencyData {
		Currencies[d.s] = ConcreteCurrency{decimalPlaces: d.d, minimumUnit: d.m, symbol: d.s}
	}
}

// Currency represents a support currency type
type Currency interface {
	Symbol() string
	DecimalPlaces() uint
	MinimumUnit() int64
	ValidateAmount(int64) bool
}

// ConcreteCurrency defines the attributes of a concrete currency type
type ConcreteCurrency struct {
	decimalPlaces uint
	minimumUnit   int64
	symbol        string
}

// Symbol returns the 3-letter ISO 4217 code
func (c ConcreteCurrency) Symbol() string {
	return c.symbol
}

// DecimalPlaces returns the number of decimals of the currency
func (c ConcreteCurrency) DecimalPlaces() uint {
	return c.decimalPlaces
}

// MinimumUnit returns the minimum amount for the currency
func (c ConcreteCurrency) MinimumUnit() int64 {
	return c.minimumUnit
}

// ValidateAmount checks whether an amount is valid for the currency
func (c ConcreteCurrency) ValidateAmount(amount int64) bool {
	return (amount % c.MinimumUnit()) == 0
}
