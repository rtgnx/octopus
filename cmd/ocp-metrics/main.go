package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"octopus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ocp           *octopus.Octopus
	electricTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   "ocp",
		Subsystem:   "electric",
		Name:        "total",
		Help:        "total consumed kWh",
		ConstLabels: map[string]string{},
	})
	gasTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace:   "ocp",
		Subsystem:   "gas",
		Name:        "total",
		Help:        "total consumed m^3",
		ConstLabels: map[string]string{},
	})
)

func init() {
	var err error
	ocp, err = octopus.New(os.Getenv("OCP_ID"), os.Getenv("OCP_KEY"), &octopus.Opts{
		Gas:      octopus.Meter{MPAN: os.Getenv("GAS_MPAN"), SN: os.Getenv("GAS_SN")},
		Electric: octopus.Meter{MPAN: os.Getenv("ELECTRIC_MPAN"), SN: os.Getenv("ELECTRIC_SN")},
	})
	if err != nil {
		log.Fatal(err)
	}
	record()
}

func record() {
	go func() {
		for {
			csr, _ := ocp.ElectricityConsuption(&octopus.ConsumptionOpts{
				GroupBy:  "hour",
				PageSize: 25000,
			})
			electricTotal.Set(csr.Total())
			log.Println(csr.String())

			csr, _ = ocp.GasConsuption(&octopus.ConsumptionOpts{
				GroupBy:  "hour",
				PageSize: 25000,
			})
			gasTotal.Set(csr.Total())
			log.Println(csr.String())

			time.Sleep(time.Hour)
		}
	}()
}

func main() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
