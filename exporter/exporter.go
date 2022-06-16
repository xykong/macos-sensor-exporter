package exporter

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/dkorunic/iSMC/output"
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
				fmt.Print(value, v[idx:])
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
		}
	default:
		fmt.Printf("I don't know about type %T!\n", v)
	}

	return 0
}

// Describe implements prometheus.Collector.
func (l *SensorsCollector) Describe(ch chan<- *prometheus.Desc) {
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
	// output.OutputFactory("ascii").All()

	log.Debug(output.GetAll())

	collector := NewSensorsCollector()
	prometheus.MustRegister(collector)

	pattern := viper.GetString("pattern")
	port := viper.GetInt("port")

	addr := fmt.Sprintf(":%d", port)
	log.Infof("staring server listening at %s%s", addr, pattern)

	http.Handle(pattern, promhttp.Handler())
	err := http.ListenAndServe(addr, nil)
	log.Error(err)
}
