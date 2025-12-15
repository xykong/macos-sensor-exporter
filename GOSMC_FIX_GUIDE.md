# 如何正确使用 github.com/xykong/gosmc 修复版本

## 问题说明
当前 `github.com/xykong/gosmc` 的 `go.mod` 文件中模块路径仍然是 `github.com/panotza/gosmc`，导致 Go 无法正确识别。

## 解决方案

### 步骤 1: 修复 gosmc 仓库的 go.mod

在你的 `github.com/xykong/gosmc` 仓库中：

1. 克隆仓库（如果还没有）：
```bash
git clone https://github.com/xykong/gosmc.git
cd gosmc
```

2. 修改 `go.mod` 文件，将第一行改为：
```go
module github.com/xykong/gosmc
```

3. 提交更改：
```bash
git add go.mod
git commit -m "Fix module path to github.com/xykong/gosmc"
```

4. 打一个新的标签（例如 v1.0.3）：
```bash
git tag v1.0.3
git push origin v1.0.3
git push origin master  # 或者 main
```

### 步骤 2: 更新 macos-sensor-exporter 项目

回到 `macos-sensor-exporter` 项目，使用 replace 指令：

```bash
cd /path/to/macos-sensor-exporter
```

在 `go.mod` 文件中添加：
```go
replace github.com/panotza/gosmc => github.com/xykong/gosmc v1.0.3
```

然后运行：
```bash
go mod tidy
go build
```

### 步骤 3: 发布版本前处理

**重要：** 当你准备发布 `macos-sensor-exporter` 新版本时，有两种选择：

#### 选项 A（推荐）：Fork iSMC 项目并更新其依赖

1. Fork `github.com/dkorunic/iSMC` 到 `github.com/xykong/iSMC`
2. 在你 fork 的 iSMC 项目中，修改 `go.mod`，将 `github.com/panotza/gosmc` 的依赖改为 `github.com/xykong/gosmc v1.0.3`
3. 提交并打 tag（例如 v0.7.1）
4. 在 `macos-sensor-exporter` 中使用你 fork 的 iSMC：
   ```go
   require (
       github.com/xykong/iSMC v0.7.1
       // ... 其他依赖
   )
   ```
5. 移除所有 replace 指令
6. 这样其他用户就可以通过 `go install` 安装了

#### 选项 B：保持当前方案，文档说明

如果你不想 fork iSMC，可以：

1. 保持 `replace` 指令在 go.mod 中
2. 在 README.md 中说明用户不能使用 `go install`，需要：
   ```bash
   git clone https://github.com/xykong/macos-sensor-exporter.git
   cd macos-sensor-exporter
   go build
   ```
3. 在 Release 中提供编译好的二进制文件供用户下载

## 推荐方案

**最佳实践是选项 A**：Fork iSMC 并更新其依赖，这样可以：
- ✅ 保持 Go modules 最佳实践
- ✅ 允许用户使用 `go install` 安装
- ✅ 避免 replace 指令的问题
- ✅ 更容易维护和分发

## 当前状态

目前项目已经移除了 replace 指令，使用原版 `github.com/panotza/gosmc v1.0.0`。如果要切换到你的修复版本，请按照上述步骤操作。
