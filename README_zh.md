# macos-sensor-exporter

一个用于 macOS 硬件传感器的 Prometheus 导出器，支持从系统管理控制器（SMC）读取温度、电压、电流、功率、电池和风扇指标。

[English](README.md) | 简体中文

## 功能特性

- 以 Prometheus 格式导出 macOS 硬件传感器指标
- 支持多种传感器类型：
  - 温度（°C）
  - 电压（V）
  - 电流（A）
  - 功率（W）
  - 风扇转速（RPM）
  - 电池信息
- 健康检查端点
- 支持通过 CLI 参数、配置文件或环境变量进行配置
- 可在终端直接显示传感器信息（支持表格、JSON 或 ASCII 格式）

## 系统要求

- macOS（在 macOS 10.13+ 上测试通过）
- Go 1.22+（用于从源码构建）
- 可能需要 root/管理员权限来访问 SMC 传感器

## 安装

### 从源码构建

```bash
git clone https://github.com/xykong/macos-sensor-exporter.git
cd macos-sensor-exporter
go build -o macos-sensor-exporter .
```

### 使用 Go Install

```bash
go install github.com/xykong/macos-sensor-exporter@latest
```

## 使用方法

### 启动 Prometheus 导出器

在默认端口（9101）上启动导出器服务：

```bash
./macos-sensor-exporter start
```

使用自定义端口和指标路径：

```bash
./macos-sensor-exporter start --port 8080 --pattern /custom-metrics
```

启用详细日志：

```bash
./macos-sensor-exporter start -v
```

### 显示传感器信息

直接在终端显示传感器信息：

```bash
# ASCII 格式（默认）
./macos-sensor-exporter show

# 表格格式
./macos-sensor-exporter show -o table

# JSON 格式
./macos-sensor-exporter show -o json
```

## 配置

导出器可以通过以下方式配置（优先级从高到低）：

1. **命令行参数**（最高优先级）
2. **环境变量**（使用 `VIPER_` 前缀）
3. **配置文件**（最低优先级）

### 配置文件

在家目录或当前目录创建 `.macos-sensor-exporter.yaml` 文件：

```yaml
port: 9101
pattern: /metrics
```

或指定自定义配置文件位置：

```bash
./macos-sensor-exporter start --config /path/to/config.yaml
```

## Prometheus 配置

在 `prometheus.yml` 中添加以下配置：

```yaml
scrape_configs:
  - job_name: 'macos-sensors'
    static_configs:
      - targets: ['localhost:9101']
```

## 导出的指标

导出器提供以下格式的指标：

```
sensor_<类别>_<描述>_<单位>{index="<编号>"} <数值>
```

示例指标：

```
sensor_temperature_cpu_die_celsius 45.5
sensor_voltage_cpu_core_volts 1.2
sensor_power_cpu_total_watt 15.3
sensor_fans_speed_rpm{index="0"} 1800
sensor_battery_charge_amperes 2.5
```

### 指标类别

- **Temperature（温度）**：CPU、GPU 和其他组件的温度
- **Voltage（电压）**：CPU 核心、GPU 和系统电压
- **Current（电流）**：电池和电源供应电流
- **Power（功率）**：CPU、GPU 和系统总功耗
- **Fans（风扇）**：所有已安装风扇的转速
- **Battery（电池）**：电池状态和指标

## 端点

- `/metrics` - Prometheus 指标端点（默认，可配置）
- `/healthz` - 健康检查端点（返回 200 OK）

## 开发

### 构建

```bash
make build
```

### 运行测试

```bash
go test ./...
```

### 使用详细日志运行

```bash
./macos-sensor-exporter start -v
```

## 项目结构

```
.
├── main.go              # 入口点
├── cmd/                 # CLI 命令
│   ├── root.go         # 根命令和配置
│   ├── start.go        # 启动导出器服务
│   └── show.go         # 显示传感器信息
└── exporter/           # Prometheus 导出器逻辑
    └── exporter.go     # Collector 实现
```

## 故障排除

### 权限被拒绝

如果遇到访问 SMC 传感器的权限错误，尝试使用 `sudo` 运行：

```bash
sudo ./macos-sensor-exporter start
```

### 没有可用的指标

确保你的 Mac 支持 SMC 传感器访问。某些虚拟化或较旧的 Mac 型号可能只有有限的传感器可用性。

### 连接被拒绝

检查端口是否已被占用：

```bash
lsof -i :9101
```

## 贡献

欢迎贡献！请随时提交 Pull Request。

## 许可证

详见 [LICENSE](LICENSE) 文件。

## 致谢

本项目使用：
- [iSMC](https://github.com/dkorunic/iSMC) - 用于 SMC 传感器访问
- [Prometheus client_golang](https://github.com/prometheus/client_golang) - 用于指标导出
- [Cobra](https://github.com/spf13/cobra) - 用于 CLI
- [Viper](https://github.com/spf13/viper) - 用于配置管理

## 相关项目

- [node_exporter](https://github.com/prometheus/node_exporter) - 用于硬件和操作系统指标的 Prometheus 导出器（Linux）
- [iSMC](https://github.com/dkorunic/iSMC) - macOS SMC 工具和库
