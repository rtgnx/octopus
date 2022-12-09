package octopus

import "fmt"

func (csr *ConsumptionReport) Total() (total float64) {
	for _, v := range csr.Results {
		total += float64(v.Consumption)
	}

	return total
}

func (csr *ConsumptionReport) Avg() (avg float64) {
	avg = csr.Total() / float64(len(csr.Results))
	return
}

func (csr *ConsumptionReport) String() string {
	return fmt.Sprintf("total: %v %s avg: %v %s", csr.Total(), csr.unit, csr.Avg(), csr.unit)
}

func (csr *ConsumptionReport) Count() int {
	return len(csr.Results)
}
