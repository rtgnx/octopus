package octopus

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
)

const (
	Daily   = "day"
	Weekly  = "week"
	Hourly  = "hour"
	Monthly = "month"
)

var (
	DefaultOpts = &ConsumptionOpts{
		PageSize: 25000,
		GroupBy:  Daily,
	}
)

func New(account string, key string, opts *Opts) (*Octopus, error) {
	if opts == nil {
		return nil, fmt.Errorf("meter options required")
	}

	return &Octopus{opts.Gas, opts.Electric, account, key}, nil
}

func (o *Octopus) request(endpoint string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, BaseURL, nil)
	req.SetBasicAuth(o.key, "")

	req.URL.Path = endpoint
	return req
}

func (o *Octopus) ElectricityConsuption(opts *ConsumptionOpts) (ConsumptionReport, error) {

	endpoint := path.Join(
		"/", API, "electricity-meter-points", o.Electric.Endpoint(), "consumption/",
	)

	req := o.request(endpoint)
	if opts != nil {
		query := req.URL.Query()
		for k, v := range opts.Map() {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()
		log.Println(req.URL.String())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ConsumptionReport{}, err
	}

	csr := &ConsumptionReport{Results: make([]CSP, 0), unit: "kWh"}

	return *csr, json.NewDecoder(res.Body).Decode(csr)
}

func (o *Octopus) GasConsuption(opts *ConsumptionOpts) (ConsumptionReport, error) {

	endpoint := path.Join(
		"/", API, "gas-meter-points", o.Gas.Endpoint(), "consumption/",
	)

	req := o.request(endpoint)
	if opts != nil {
		query := req.URL.Query()
		for k, v := range opts.Map() {
			query.Add(k, v)
		}
		req.URL.RawQuery = query.Encode()
		log.Println(req.URL.String())
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ConsumptionReport{}, err
	}

	csr := &ConsumptionReport{Results: make([]CSP, 0), unit: "m^3"}

	return *csr, json.NewDecoder(res.Body).Decode(csr)
}
