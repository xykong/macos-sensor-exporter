# 如何触发 GitHub Actions Release（包含最新修复）

## 背景

最新的提交 `b9b5ca2`（test: skip SMC-dependent tests when hardware unavailable）已经修复了 GitHub Actions 测试失败的问题。该修复已经合并到 `master` 分支，但还没有发布新版本。

## 当前状态

- **最新提交**: `b9b5ca2` - test: skip SMC-dependent tests when hardware unavailable
- **最新 tag**: `v1.0.3`
- **修复内容**: 测试代码会在 SMC 硬件不可用时自动跳过，不会导致 CI/CD 失败

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

## 版本号建议

根据语义化版本规范（Semantic Versioning）：

- **v1.0.4** (推荐) - 修复 bug，向下兼容
- **v1.1.0** - 如果有新功能
- **v2.0.0** - 如果有破坏性变更

当前修复是 bug fix，建议使用 **v1.0.4**。

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

## 快速执行

如果你已经确定要发布 v1.0.4：

```bash
cd /Users/xykong/workspace/xykong/macos-sensor-exporter-project/macos-sensor-exporter && \
git tag -a v1.0.4 -m "Fix: Skip SMC-dependent tests when hardware unavailable" && \
git push origin v1.0.4 && \
echo "✅ Tag v1.0.4 已推送，GitHub Actions 正在运行..." && \
gh run watch
```
