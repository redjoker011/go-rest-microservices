package data

import (
	"testing"
)

func TestChecksValidateion(t *testing.T) {
	p := &Product{
		Name:  "Cappucino",
		Price: 1.0,
		SKU:   "abc-a-a",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
