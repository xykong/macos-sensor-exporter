# macos-sensor-exporter

A Prometheus exporter for macOS hardware sensors, including temperature, voltage, current, power, battery, and fan metrics from the System Management Controller (SMC).

English | [简体中文](README_zh.md)

## Features

- Exports macOS hardware sensor metrics in Prometheus format
- Support for multiple sensor types:
  - Temperature (°C)
  - Voltage (V)
  - Current (A)
  - Power (W)
  - Fan speed (RPM)
  - Battery information
- Health check endpoint
- Configurable via CLI flags, config file, or environment variables
- Display sensor information directly in terminal (table, JSON, or ASCII format)

## Requirements

- macOS (tested on macOS 10.13+)
- Go 1.22+ (for building from source)
- Root/Administrator privileges may be required to access SMC sensors

## Installation

### From Source

```bash
git clone https://github.com/xykong/macos-sensor-exporter.git
cd macos-sensor-exporter
go build -o macos-sensor-exporter .
```

### Using Go Install

```bash
go install github.com/xykong/macos-sensor-exporter@latest
```

## Usage

### Start the Prometheus Exporter

Start the exporter server on the default port (9101):

```bash
./macos-sensor-exporter start
```

With custom port and metrics path:

```bash
./macos-sensor-exporter start --port 8080 --pattern /custom-metrics
```

With verbose logging:

```bash
./macos-sensor-exporter start -v
```

### Show Sensor Information

Display sensor information directly in the terminal:

```bash
# ASCII format (default)
./macos-sensor-exporter show

# Table format
./macos-sensor-exporter show -o table

# JSON format
./macos-sensor-exporter show -o json
```

## Configuration

The exporter can be configured using:

1. **Command-line flags** (highest priority)
2. **Environment variables** (with `VIPER_` prefix)
3. **Configuration file** (lowest priority)

### Configuration File

Create a `.macos-sensor-exporter.yaml` file in your home directory or current directory:

```yaml
port: 9101
pattern: /metrics
```

Or specify a custom config file location:

```bash
./macos-sensor-exporter start --config /path/to/config.yaml
```

## Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'macos-sensors'
    static_configs:
      - targets: ['localhost:9101']
```

## Exported Metrics

The exporter provides metrics in the following format:

```
sensor_<category>_<description>_<unit>{index="<number>"} <value>
```

Example metrics:

```
sensor_temperature_cpu_die_celsius 45.5
sensor_voltage_cpu_core_volts 1.2
sensor_power_cpu_total_watt 15.3
sensor_fans_speed_rpm{index="0"} 1800
sensor_battery_charge_amperes 2.5
```

### Metric Categories

- **Temperature**: CPU, GPU, and other component temperatures
- **Voltage**: CPU core, GPU, and system voltages
- **Current**: Battery and power supply current
- **Power**: CPU, GPU, and total system power consumption
- **Fans**: Fan speeds for all installed fans
- **Battery**: Battery status and metrics

## Endpoints

- `/metrics` - Prometheus metrics endpoint (default, configurable)
- `/healthz` - Health check endpoint (returns 200 OK)

## Development

### Build

```bash
make build
```

### Run Tests

```bash
go test ./...
```

### Run with Verbose Logging

```bash
./macos-sensor-exporter start -v
```

## Architecture

The project structure:

```
.
├── main.go              # Entry point
├── cmd/                 # CLI commands
│   ├── root.go         # Root command and config
│   ├── start.go        # Start exporter server
│   └── show.go         # Show sensor info
└── exporter/           # Prometheus exporter logic
    └── exporter.go     # Collector implementation
```

## Troubleshooting

### Permission Denied

If you encounter permission errors accessing SMC sensors, try running with `sudo`:

```bash
sudo ./macos-sensor-exporter start
```

### No Metrics Available

Ensure your Mac supports SMC sensor access. Some virtualized or older Mac models may have limited sensor availability.

### Connection Refused

Check if the port is already in use:

```bash
lsof -i :9101
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

See [LICENSE](LICENSE) file for details.

## Credits

This project uses:
- [iSMC](https://github.com/dkorunic/iSMC) for SMC sensor access
- [Prometheus client_golang](https://github.com/prometheus/client_golang) for metrics export
- [Cobra](https://github.com/spf13/cobra) for CLI
- [Viper](https://github.com/spf13/viper) for configuration management

## Related Projects

- [node_exporter](https://github.com/prometheus/node_exporter) - Prometheus exporter for hardware and OS metrics (Linux)
- [iSMC](https://github.com/dkorunic/iSMC) - macOS SMC tool and library
