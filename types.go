package octopus

import (
	"fmt"
	"time"
)

const (
	API     = "v1"
	BaseURL = "https://api.octopus.energy"
)

type Meter struct {
	MPAN string
	SN   string
}

func (m Meter) Endpoint() string {
	return fmt.Sprintf("%s/meters/%s", m.MPAN, m.SN)
}

type Opts struct {
	Gas      Meter
	Electric Meter
}

type Octopus struct {
	Gas, Electric Meter
	account, key  string
}

type CSP struct {
	Consumption   float32   `json:"consumption"`
	IntervalStart time.Time `json:"interval_start"`
	IntervalEnd   time.Time `json:"interval_end"`
}

type ConsumptionReport struct {
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
	Results  []CSP  `json:"results"`
	unit     string
}

type ConsumptionOpts struct {
	PageSize   int
	GroupBy    string
	PeriodFrom time.Time
	PeriodTo   time.Time
}

func (copts *ConsumptionOpts) Map() (m map[string]string) {
	m = make(map[string]string)
	m["page_size"] = fmt.Sprintf("%d", copts.PageSize)
	m["group_by"] = copts.GroupBy
	if !copts.PeriodFrom.IsZero() {
		m["period_from"] = copts.PeriodFrom.String()
	}
	if !copts.PeriodTo.IsZero() {
		m["period_to"] = copts.PeriodTo.String()
	}
	return m
}
