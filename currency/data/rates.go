package data

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
)

type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{log: l, rates: map[string]float64{}}

	err := er.GetRates()

	if err != nil {
		return nil, err
	}

	return er, nil
}

func (er *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := er.rates[base]

	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", base)
	}

	dr, ok := er.rates[dest]

	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", dest)
	}

	return dr / br, nil
}

func (er *ExchangeRates) GetRates() error {
	var url = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml?"
	resp, err := http.DefaultClient.Get(url)

	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected error code 200 got %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	md := &Cubes{}

	xml.NewDecoder(resp.Body).Decode(&md)

	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}

		er.rates[c.Currency] = r
	}

	return nil
}

type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}

type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

// MonitorRates checks the rates in the ECB API every interval and sends a message to the
// returned channel when there are changes
//
// Note: the ECB API only returns data once a day, this function only simulates the changes
// in rates for demonstration purposes
func (e *ExchangeRates) MonitorRates(interval time.Duration) chan struct{} {
	ret := make(chan struct{})

	go func() {
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ticker.C:
				// just add a random difference to the rate and return it
				// this simulates the fluctuations in currency rates
				for k, v := range e.rates {
					// change can be 10% of original value
					change := (rand.Float64() / 10)
					// is this a postive or negative change
					direction := rand.Intn(1)

					if direction == 0 {
						// new value with be min 90% of old
						change = 1 - change
					} else {
						// new value will be 110% of old
						change = 1 + change
					}

					// modify the rate
					e.rates[k] = v * change
				}

				// notify updates, this will block unless there is a listener on the other end
				ret <- struct{}{}
			}
		}
	}()

	return ret
}
