package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
// TODO - define a histogram that will measure the time spent in the handler for "/helloworld"
	timeSpent = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "time_spent"
		Help: "the time spent in the handler for /helloworld" 
	})
// TODO - define a counter vec that will monitor the type of request for "/helloworld"
	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "the type of request",
			Help: "Monitoring the type of request for /helloworld",
		},
		[]string{"counter"},
	)

// TODO - define a counter that will measure the request rate for "/helloworld"
	counter2 = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "request rate"
			Help: "Measuring the request rate"
		}
	)

// TODO - define an error counter for "/sum"
	counterError = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "error counter"
			Help: "Error counter for /sum"
		}
	)

// TODO - define a gauge to measure the uptime
	uptime_mes = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "uptime",
			Help: "Measuring the uptime",
	}))
)

func init() {
	// Metrics have to be registered to be exposed

	// TODO - register the prometheus metrics
	prometheus.MustRegister(timeSpent)
	prometheus.MustRegister(counter)
	prometheus.MustRegister(counter2)
	prometheus.MustRegister(counterError)
	prometheus.MustRegister(uptime_mes)

}

func main() {
	// TODO - add handler for prometheus
	http.Handle("/metrics", promhttp.Handler())

	// TODO - start a go routine that will measure the uptime
	//        create ticker that will trigger each second
	//        and increment the gauge for the uptime by one

	ticker := time.NewTicker(500 * time.Millisecond)
	go func() {
        for t := range <-ticker.C {
            	fmt.Println("Tick at", t)
				uptime_mes.Inc()
        }
    }()

	http.HandleFunc("/helloworld", func(w http.ResponseWriter, r *http.Request) {
		// TODO - measure the time using the histogram that you have defined and a prometheus timer


		// TODO - increment counter for request rate
		counter2.Inc()

		// TODO - increment the counter vec for HTTP methods

		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/sum", func(w http.ResponseWriter, r *http.Request) {
		strA := r.URL.Query().Get("a")
		strB := r.URL.Query().Get("b")

		a, errA := strconv.Atoi(strA)
		b, errB := strconv.Atoi(strB)

		if errA != nil || errB != nil {
			// TODO - increment error counter for "/sum"
			counterError++
			return
		}

		fmt.Fprintf(w, fmt.Sprintf("%d + %d = %d", a, b, a+b))
	})

	// TODO BONUS - using the os package configure the port using an environment variable
	//              you will have to provide a default value if this is empty. Then configure the
	//              the port in docker-compose and try to run it again
	fmt.Printf("Server running (port=8080)\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
