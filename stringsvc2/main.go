package main

import (
	"net/http"
	"os"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log" // NB: we replace std log package with this one
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
)

func main() {
	// initialise all dependencies required for middleware chain
	// just the main method in this file such that all focus is on control flow as opposed to abstraction effort
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	// service initialisation and middleware chain
	var svc StringService
	svc = stringService{}                                                         // barebones service implementation - JUST business logic
	svc = loggingMiddleware{logger, svc}                                          // chain middleare onto service to decorate service with logging
	svc = instrumentingMiddleware{requestCount, requestLatency, countResult, svc} // chain more middleware onto updated service to decorate with prometheus observability
	// notice onion effect of core service being built out to include additional services
	// the updated service is passed into each middleware, thereby maintaining all added middleware and creating the onion
	// refer to README for an explanation of how middleware adds logic
	// essentially, the middleware also comes in the form of an interface that lists the same methods as the service interface
	// therefore, one can implement the middleware interface by adding the logic that necessitates the middleware interfacen and then calling the service interface's implementation of business logic

	// handlers link endpoints to transport
	uppercaseHandler := httptransport.NewServer(
		makeUppercaseEndpoint(svc),
		decodeUppercaseRequest,
		encodeResponse,
	)

	countHandler := httptransport.NewServer(
		makeCountEndpoint(svc),
		decodeCountRequest,
		encodeResponse,
	)

	// router
	http.Handle("/uppercase", uppercaseHandler)
	http.Handle("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", ":8080")

	// serve
	logger.Log("err", http.ListenAndServe(":8080", nil))
}
