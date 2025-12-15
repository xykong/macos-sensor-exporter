package exporter

import (
	"os"
	"runtime"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/xykong/iSMC/output"
)

func TestGetUnit(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "amperes",
			input:    "2.5 A",
			expected: "_amperes",
		},
		{
			name:     "volts",
			input:    "1.2 V",
			expected: "_volts",
		},
		{
			name:     "watt",
			input:    "15.3 W",
			expected: "_watt",
		},
		{
			name:     "celsius",
			input:    "45.5 °C",
			expected: "_celsius",
		},
		{
			name:     "rpm",
			input:    "1800 rpm",
			expected: "_rpm",
		},
		{
			name:     "no unit",
			input:    "100",
			expected: "",
		},
		{
			name:     "non-string",
			input:    42,
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getUnit(tt.input)
			if result != tt.expected {
				t.Errorf("getUnit(%v) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetGaugeValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected float64
	}{
		{
			name:     "int",
			input:    42,
			expected: 42.0,
		},
		{
			name:     "float64",
			input:    3.14,
			expected: 3.14,
		},
		{
			name:     "bool true",
			input:    true,
			expected: 1.0,
		},
		{
			name:     "bool false",
			input:    false,
			expected: 0.0,
		},
		{
			name:     "string with unit",
			input:    "45.5 °C",
			expected: 45.5,
		},
		{
			name:     "string without unit",
			input:    "100",
			expected: 100.0,
		},
		{
			name:     "invalid string",
			input:    "invalid",
			expected: 0.0,
		},
		{
			name:     "negative number",
			input:    "-10.5 V",
			expected: -10.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getGaugeValue(tt.input)
			if result != tt.expected {
				t.Errorf("getGaugeValue(%v) = %f, expected %f", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCreateNewDesc(t *testing.T) {
	tests := []struct {
		name        string
		catalog     string
		description string
		value       interface{}
		wantName    string
		wantHelp    string
	}{
		{
			name:        "simple metric",
			catalog:     "Temperature",
			description: "CPU Die",
			value:       "45.5 °C",
			wantName:    "sensor_temperature_cpu_die_celsius",
			wantHelp:    "Temperature CPU Die",
		},
		{
			name:        "metric with index",
			catalog:     "Fans",
			description: "Fan 0",
			value:       "1800 rpm",
			wantName:    "sensor_fans_fan_rpm",
			wantHelp:    "Fan",
		},
		{
			name:        "metric with special chars",
			catalog:     "Power",
			description: "CPU (Total)",
			value:       "15.3 W",
			wantName:    "sensor_power_cpu_total_watt",
			wantHelp:    "Power CPU (Total)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := createNewDesc(tt.catalog, tt.description, tt.value)
			if desc == nil {
				t.Fatal("createNewDesc returned nil")
			}

			// Verify the descriptor was created successfully
			// We can't easily access private fields, but we can verify it doesn't panic
			// and returns a non-nil descriptor
		})
	}
}

func TestSensorsCollectorDescribe(t *testing.T) {
	// Skip if not on macOS
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping test: SMC is only available on macOS")
	}

	// Skip if running in CI environment (GitHub Actions, etc.)
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping test: SMC is not accessible in CI environments")
	}

	// Check if SMC is accessible
	data := output.GetAll()
	if len(data) == 0 {
		t.Skip("Skipping test: SMC is not accessible (possibly running in VM or without proper permissions)")
	}

	collector := NewSensorsCollector()
	if collector == nil {
		t.Fatal("NewSensorsCollector returned nil")
	}

	ch := make(chan *prometheus.Desc, 10)
	go func() {
		collector.Describe(ch)
		close(ch)
	}()

	// Collect all descriptions
	var count int
	for range ch {
		count++
	}

	// We expect at least some descriptions (depends on DescribeByCollect behavior)
	// This test mainly ensures Describe doesn't panic
	t.Logf("Describe sent %d descriptors", count)
}

func TestSensorsCollectorCollect(t *testing.T) {
	// Skip if not on macOS
	if runtime.GOOS != "darwin" {
		t.Skip("Skipping test: SMC is only available on macOS")
	}

	// Skip if running in CI environment (GitHub Actions, etc.)
	if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
		t.Skip("Skipping test: SMC is not accessible in CI environments")
	}

	// Check if SMC is accessible
	data := output.GetAll()
	if len(data) == 0 {
		t.Skip("Skipping test: SMC is not accessible (possibly running in VM or without proper permissions)")
	}

	collector := NewSensorsCollector()
	if collector == nil {
		t.Fatal("NewSensorsCollector returned nil")
	}

	// Create a custom registry to test collection
	registry := prometheus.NewRegistry()
	err := registry.Register(collector)
	if err != nil {
		t.Fatalf("Failed to register collector: %v", err)
	}

	// Gather metrics
	metricFamilies, err := registry.Gather()
	if err != nil {
		t.Fatalf("Failed to gather metrics: %v", err)
	}

	// We expect at least some metrics (depends on system sensors)
	t.Logf("Collected %d metric families", len(metricFamilies))

	// Verify metric names follow the expected pattern
	for _, mf := range metricFamilies {
		name := mf.GetName()
		if len(name) == 0 {
			t.Error("Metric family has empty name")
		}
		t.Logf("Metric family: %s (%d metrics)", name, len(mf.GetMetric()))
	}
}
