package main

import (
	"math/rand"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func getMetrics(w http.ResponseWriter, r *http.Request) {
	initmetrics()
	promhttp.Handler().ServeHTTP(w, r)
}

func initmetrics() {
	tmp := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "sennser",
		Name:      "value",
		Help:      "Random number",
		Subsystem: "temp",
		ConstLabels: prometheus.Labels{
			"instance": "test",
		},
	})
	hum := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "sennser",
		Name:      "value",
		Help:      "Random number",
		Subsystem: "hum",
		ConstLabels: prometheus.Labels{
			"instance": "test",
		},
	})
	var add []prometheus.Collector
	num := rand.Intn(2)

	add = append(add, tmp)
	if num == 0 {
		add = append(add, hum)
	}

	// prometheus.MustRegister(a)
	prometheus.Unregister(tmp)
	prometheus.Unregister(hum)
	prometheus.MustRegister(add...)
	//ランダムな値を取得
	i := rand.Intn(100)
	j := rand.Intn(1000)
	tmp.Set(float64(i) / 10)
	hum.Set(float64(j) / 10)

}
