package types

import "testing"

func TestSymbol(t *testing.T) {
	testCases := []struct {
		dec     uint
		min     int64
		sym     string
		symRepr string
	}{
		{2, 1, "USD", "USD"},
		{2, 1, "GBP", "GBP"},
		{2, 1, "EUR", "EUR"},
	}
	for i, tc := range testCases {
		c := ConcreteCurrency{tc.dec, tc.min, tc.sym}
		symRepr := c.Symbol()
		if symRepr != tc.symRepr {
			t.Errorf("%d:%v: Symbol() return %q; should be %v",
				i, tc, symRepr, tc.symRepr)
		}
	}
}

func TestValidateAmount(t *testing.T) {
	var amount int64
	amount = 100
	c := ConcreteCurrency{2, 1, "EUR"}
	if !c.ValidateAmount(amount) {
		t.Errorf("c (%s):%d should be valid", c.Symbol(), amount)
	}
	c = ConcreteCurrency{2, 3, "XXX"}
	if c.ValidateAmount(amount) {
		t.Errorf("c (%s):%d should not be valid (modulo: %d)",
			c.Symbol(), amount, amount%c.MinimumUnit())
	}
}
