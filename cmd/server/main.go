package main

import (
	"net/http"
	"strconv"
	"strings"
)

type metricValue interface {
	changeValue(int)
}

func changeValue(m metricValue, val int) {
	m.changeValue(val)
}

type gauge float64
func (g *gauge) changeValue(val int) {
	*g = gauge(val)
}

type counter int64
func (c *counter) changeValue(val int) {
	*c += counter(val)
}

type MemStorage struct {
	metrics map[string]metricValue
}

var memStorage MemStorage = MemStorage{metrics: make(map[string]metricValue)}

func metricsPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	r.Header.Set("Content-Type", "text/plain")
	//split URL Metrics
	splitMetrics := strings.Split(r.RequestURI[1:], "/")
	typeMetrics, nameMetrics := splitMetrics[1], splitMetrics[2]
	valueMetrics, err := strconv.Atoi(splitMetrics[3])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if val, ok := memStorage.metrics[nameMetrics]; ok {
		changeValue(val,valueMetrics)
		return
	} else {
		switch typeMetrics {
		case "counter":
			c := counter(valueMetrics)
			memStorage.metrics[nameMetrics] = &c
		case "gauge":
			g := gauge(valueMetrics)
			memStorage.metrics[nameMetrics] = &g
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

//disable redirect and split URL
func CheckURLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := strings.Split(r.RequestURI[1:], "/")
		//checking method, type and name metrics
		if url[0] != "update" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if url[1] != "counter" && url[1] != "gauge" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if len(url[2]) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	rtr := http.NewServeMux()
	rtr.HandleFunc(`/`, metricsPost)
	middle := CheckURLMiddleware(rtr)

	if err := http.ListenAndServe(`:8080`, middle); err != nil {
		panic(err)
	}
}
