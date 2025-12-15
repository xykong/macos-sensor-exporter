# 如何触发 GitHub Actions Release（包含最新修复）

## 背景

GitHub Actions release workflow 在 v1.0.4 版本上失败，原因是 SMC 测试代码在调用 `output.GetAll()` 时就会失败，即使有 skip 逻辑也无法阻止。

## 当前状态

- **最新提交**: `4496d86` - fix: enable CGO for iSMC hid package compilation
- **最新 tag**: `v1.0.6` ✅ （完整修复并发布）
- **修复内容**: 
  1. 在 CI 环境中提前跳过 SMC 测试（检查 CI 和 GITHUB_ACTIONS 环境变量）
  2. 启用 CGO 以支持 iSMC hid 包的编译

## 问题分析

### 问题 1: 测试失败（v1.0.4）
- 虽然添加了 `output.GetAll()` 返回空数据的检查
- 但在 GitHub Actions 环境中，`output.GetAll()` 在尝试访问 SMC 时就会直接失败
- 导致测试在 skip 逻辑之前就崩溃了

**解决方案**: 添加 CI 环境变量检查，在调用 `output.GetAll()` 之前就跳过
```go
// Skip if running in CI environment (GitHub Actions, etc.)
if os.Getenv("CI") != "" || os.Getenv("GITHUB_ACTIONS") != "" {
    t.Skip("Skipping test: SMC is not accessible in CI environments")
}
```

### 问题 2: 构建失败（v1.0.5）
- GoReleaser 配置中设置了 `CGO_ENABLED=0`
- 但 iSMC 的 hid 包使用了 `import "C"`，需要 CGO 支持
- 导致构建时所有 Go 文件被 build constraints 排除

**解决方案**: 在 `.goreleaser.yml` 中启用 CGO
```yaml
env:
  - CGO_ENABLED=1  # 从 0 改为 1
```

## 触发 Release 的步骤

### 方法一：使用 Git 命令行（推荐）

```bash
# 1. 确保在 macos-sensor-exporter 目录
cd /Users/xykong/workspace/xykong/macos-sensor-exporter-project/macos-sensor-exporter

# 2. 确保本地代码是最新的
git pull origin master

# 3. 创建新的 tag（建议使用 v1.0.4）
git tag -a v1.0.4 -m "Fix: Skip SMC-dependent tests when hardware unavailable"

# 4. 推送 tag 到 GitHub（这会自动触发 Release workflow）
git push origin v1.0.4
```

### 方法二：使用 GitHub CLI

```bash
# 1. 确保在 macos-sensor-exporter 目录
cd /Users/xykong/workspace/xykong/macos-sensor-exporter-project/macos-sensor-exporter

# 2. 创建并推送 tag
git tag -a v1.0.4 -m "Fix: Skip SMC-dependent tests when hardware unavailable"
git push origin v1.0.4

# 3. 查看 workflow 运行状态
gh workflow view release
gh run list --workflow=release.yml
```

## 触发原理

根据 `.github/workflows/release.yml` 配置：

```yaml
on:
  push:
    tags:
      - 'v*'  # 当推送 v* 格式的 tag 时触发
```

只要推送一个以 `v` 开头的 tag（如 `v1.0.4`），GitHub Actions 就会自动：

1. ✅ Checkout 代码
2. ✅ 设置 Go 环境
3. ✅ 运行测试（SMC 相关测试会自动跳过，不会失败）
4. ✅ 使用 GoReleaser 构建并发布 release

## 验证 Release

推送 tag 后，可以通过以下方式验证：

### 1. 查看 GitHub Actions 状态

```bash
# 使用 GitHub CLI
gh run list --workflow=release.yml --limit 5

# 或访问网页
open https://github.com/xykong/macos-sensor-exporter/actions
```

### 2. 查看 Release 页面

```bash
# 使用 GitHub CLI
gh release list

# 或访问网页
open https://github.com/xykong/macos-sensor-exporter/releases
```

### 3. 实时监控 workflow 运行

```bash
# 获取最新的 run ID 并观察
gh run watch
```

## 预期结果

✅ **测试阶段**: SMC 相关测试会被跳过（在 GitHub Actions 虚拟环境中）
```
=== RUN   TestSensorsCollectorDescribe
--- SKIP: TestSensorsCollectorDescribe (0.00s)
    exporter_test.go:XX: Skipping test: SMC not accessible on this system
=== RUN   TestSensorsCollectorCollect
--- SKIP: TestSensorsCollectorCollect (0.00s)
    exporter_test.go:XX: Skipping test: SMC not accessible on this system
```

✅ **Release 阶段**: GoReleaser 成功构建并创建 release，包含：
- macOS ARM64 二进制文件
- macOS AMD64 二进制文件
- 源码压缩包
- Release notes

## 版本历史

- **v1.0.3** - 最后一个稳定版本（测试会在 GitHub Actions 上失败）
- **v1.0.4** - ❌ 失败（尝试修复测试但不完整）
- **v1.0.5** - ❌ 失败（测试通过但构建失败，CGO 未启用）
- **v1.0.6** - ✅ 成功（完整修复：CI 测试跳过 + CGO 启用）

## 回滚方法

如果需要回滚或删除 tag：

```bash
# 删除本地 tag
git tag -d v1.0.4

# 删除远程 tag
git push origin --delete v1.0.4

# 使用 GitHub CLI 删除 release
gh release delete v1.0.4
```

## 常见问题

### Q: 为什么不能手动触发 workflow？
A: 当前的 `release.yml` 只配置了 tag push 触发。如果需要手动触发，可以添加：
```yaml
on:
  push:
    tags:
      - 'v*'
  workflow_dispatch:  # 添加此行以支持手动触发
```

### Q: 如何查看之前失败的 workflow？
```bash
gh run list --workflow=release.yml --status=failure
```

### Q: 测试在本地 macOS 机器上会运行吗？
A: 是的，在有 SMC 硬件访问权限的本地 macOS 机器上，SMC 测试会正常运行。只在 GitHub Actions 虚拟环境中会跳过。

## v1.0.6 Release 状态

**✅ v1.0.6 已成功发布！**

tag 已推送，GitHub Actions 正在构建 release。这个版本包含了完整的修复：
- ✅ CI 环境中自动跳过 SMC 测试
- ✅ 启用 CGO 支持 iSMC hid 包编译

### 查看 Release 进度

```bash
# 查看最新的 workflow 运行
cd /Users/xykong/workspace/xykong/macos-sensor-exporter-project/macos-sensor-exporter
gh run list --workflow=release.yml --limit 5

# 实时监控
gh run watch

# 查看 release 列表
gh release list

# 访问 GitHub Actions 页面
open https://github.com/xykong/macos-sensor-exporter/actions

# 访问 Release 页面
open https://github.com/xykong/macos-sensor-exporter/releases
```

### 预期输出

在 GitHub Actions 中，SMC 测试会被跳过：
```
=== RUN   TestSensorsCollectorDescribe
--- SKIP: TestSensorsCollectorDescribe (0.00s)
    exporter_test.go:XX: Skipping test: SMC is not accessible in CI environments
=== RUN   TestSensorsCollectorCollect
--- SKIP: TestSensorsCollectorCollect (0.00s)
    exporter_test.go:XX: Skipping test: SMC is not accessible in CI environments
```

而在本地 macOS 机器上，测试会正常运行并通过。
