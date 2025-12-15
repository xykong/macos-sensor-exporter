package exporter

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/xykong/iSMC/output"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type (
	SensorsCollector struct{}
)

func NewSensorsCollector() *SensorsCollector {
	return &SensorsCollector{}
}

func createNewDesc(catelog, description string, value interface{}) *prometheus.Desc {
	var unit = getUnit(value)

	var help = catelog + " " + description
	var variableLabels []string = []string{}
	var constLabels prometheus.Labels = prometheus.Labels{}

	re := regexp.MustCompile("[0-9]+")
	var idx = re.FindAllStringIndex(description, -1)
	if len(idx) == 1 && idx[0][0] > 0 {
		// fmt.Println(description, description[idx[0][0]:idx[0][1]], description[:idx[0][0]]+description[idx[0][1]:])
		help = description[:idx[0][0]] + description[idx[0][1]:]

		// variableLabels = append(variableLabels, "index")
		constLabels["index"] = description[idx[0][0]:idx[0][1]]
	}

	help = strings.TrimSpace(help)
	help = strings.Replace(help, "  ", " ", -1)

	var fqName = strings.ToLower(help)
	fqName = strings.Replace(fqName, " ", "_", -1)
	fqName = strings.Replace(fqName, ".", "_", -1)
	fqName = strings.Replace(fqName, "-", "_", -1)
	fqName = strings.Replace(fqName, "(", "", -1)
	fqName = strings.Replace(fqName, ")", "", -1)
	fqName = "sensor_" + fqName + unit

	log.Debug(fqName, help, variableLabels, constLabels)

	return prometheus.NewDesc(
		fqName,
		help,
		variableLabels,
		constLabels)
}

func getUnit(value interface{}) string {
	if v, ok := value.(string); ok {
		if idx := strings.Index(v, " "); idx != -1 {

			switch v[idx:] {
			case " A":
				return "_amperes"
			case " V":
				return "_volts"
			case " W":
				return "_watt"
			case " °C":
				return "_celsius"
			case " rpm":
				return "_rpm"
			default:
				log.WithField("value", value).WithField("unit", v[idx:]).Warn("unknown unit type")
			}
		}
	}

	return ""
}

func getGaugeValue(value interface{}) float64 {

	switch v := value.(type) {
	case int:
		return float64(v)
	case float64:
		return v
	case bool:
		if v {
			return 1
		} else {
			return 0
		}
	case string:
		if idx := strings.Index(v, " "); idx != -1 {
			v = v[:idx]
		}

		if s, err := strconv.ParseFloat(v, 64); err == nil {
			return s
		} else {
			log.WithError(err).WithField("value", value).Warn("failed to parse sensor value")
		}
	default:
		log.WithField("type", fmt.Sprintf("%T", v)).WithField("value", v).Warn("unknown value type")
	}

	return 0
}

// Describe implements prometheus.Collector.
func (l *SensorsCollector) Describe(ch chan<- *prometheus.Desc) {
	// Use DescribeByCollect for dynamic metrics
	prometheus.DescribeByCollect(l, ch)
}

// Collect implements prometheus.Collector.
func (l *SensorsCollector) Collect(ch chan<- prometheus.Metric) {

	for catelog, catelogValue := range output.GetAll() {
		log.Debugf("Catelog: %s\n", catelog)

		if catelogValue, ok := catelogValue.(map[string]interface{}); ok {

			for description, details := range catelogValue {
				// log.Infof("description: %s\n", description)
				// log.Infof("details: %s\n", details.(map[string]interface{})["value"])

				var value = details.(map[string]interface{})["value"]

				ch <- prometheus.MustNewConstMetric(createNewDesc(catelog, description, value),
					prometheus.GaugeValue,
					getGaugeValue(value),
				)
			}
		}
	}
}

func Start() {
	collector := NewSensorsCollector()
	prometheus.MustRegister(collector)

	pattern := viper.GetString("pattern")
	port := viper.GetInt("port")

	addr := fmt.Sprintf(":%d", port)

	// Create explicit mux instead of using default
	mux := http.NewServeMux()
	mux.Handle(pattern, promhttp.Handler())
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	log.Infof("starting server listening at %s%s", addr, pattern)
	log.Infof("health check endpoint available at %s/healthz", addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Fatal("server exited with error")
	}
}
