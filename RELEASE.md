# Release 发布指南

本项目使用 GitHub Actions 和 GoReleaser 自动化发布流程。

## 自动发布流程

当你推送一个以 `v` 开头的 git tag 时，GitHub Actions 会自动：

1. 运行所有测试
2. 为 macOS (Intel 和 Apple Silicon) 构建二进制文件
3. 创建 GitHub Release
4. 上传构建的二进制文件和 checksums
5. 自动生成 changelog

## 如何发布新版本

### 1. 更新 CHANGELOG.md

在发布前，确保更新 `CHANGELOG.md` 文件，记录新版本的变更：

```markdown
## [1.0.0] - 2025-12-15

### Added
- 新增功能描述

### Fixed
- 修复的问题描述

### Changed
- 变更的内容描述
```

### 2. 创建并推送 tag

```bash
# 确保所有更改已提交
git add .
git commit -m "chore: prepare for release v1.0.0"

# 创建 tag（使用语义化版本号）
git tag -a v1.0.0 -m "Release v1.0.0"

# 推送代码和 tag
git push origin main
git push origin v1.0.0
```

### 3. 查看发布进度

1. 访问 GitHub 仓库的 [Actions](https://github.com/xykong/macos-sensor-exporter/actions) 页面
2. 查看 "Release" workflow 的运行状态
3. 等待构建完成（通常需要 5-10 分钟）

### 4. 验证发布

构建完成后：

1. 访问 [Releases](https://github.com/xykong/macos-sensor-exporter/releases) 页面
2. 验证新版本已发布
3. 下载并测试构建的二进制文件

```bash
# 下载并测试 (以 arm64 为例)
curl -LO https://github.com/xykong/macos-sensor-exporter/releases/download/v1.0.0/macos-sensor-exporter_v1.0.0_Darwin_arm64.tar.gz
tar -xzf macos-sensor-exporter_v1.0.0_Darwin_arm64.tar.gz
./macos-sensor-exporter --version
```

## 版本号规范

本项目遵循 [语义化版本 2.0.0](https://semver.org/lang/zh-CN/)：

- **MAJOR 主版本号**：当你做了不兼容的 API 修改
- **MINOR 次版本号**：当你做了向下兼容的功能性新增
- **PATCH 修订号**：当你做了向下兼容的问题修正

示例：
- `v1.0.0` - 首个稳定版本
- `v1.1.0` - 新增功能
- `v1.1.1` - 修复 bug
- `v2.0.0` - 破坏性更改

## 预发布版本

对于 beta 或 rc 版本，使用以下格式：

```bash
# Beta 版本
git tag -a v1.0.0-beta.1 -m "Release v1.0.0-beta.1"

# Release Candidate
git tag -a v1.0.0-rc.1 -m "Release v1.0.0-rc.1"

git push origin v1.0.0-beta.1
```

GoReleaser 会自动将这些版本标记为 "pre-release"。

## 本地测试发布流程

在推送 tag 之前，你可以在本地测试 GoReleaser 配置：

```bash
# 安装 GoReleaser (如果还没有安装)
brew install goreleaser

# 测试配置文件
goreleaser check

# 构建快照版本（不会发布）
goreleaser release --snapshot --clean

# 查看构建产物
ls -la dist/
```

## 文件说明

- `.github/workflows/release.yml` - GitHub Actions workflow 配置
- `.goreleaser.yml` - GoReleaser 构建和发布配置
- `CHANGELOG.md` - 版本变更日志

## 常见问题

### Q: 如何删除错误的 release？

1. 在 GitHub 上删除 release
2. 删除本地 tag: `git tag -d v1.0.0`
3. 删除远程 tag: `git push origin :refs/tags/v1.0.0`
4. 重新创建正确的 tag

### Q: 构建失败怎么办？

1. 检查 GitHub Actions 日志查看错误信息
2. 在本地运行 `goreleaser release --snapshot --clean` 测试
3. 修复问题后，删除旧 tag，重新创建并推送

### Q: 如何添加更多构建目标？

编辑 `.goreleaser.yml` 文件的 `builds` 部分：

```yaml
builds:
  - goos:
      - darwin
      - linux  # 添加 Linux 支持
    goarch:
      - amd64
      - arm64
```

但请注意，本项目依赖于 macOS SMC，只能在 macOS 上运行。

## 自动化改进建议

未来可以考虑添加：

1. **自动版本号更新**：使用工具自动更新版本号
2. **Homebrew 发布**：取消注释 `.goreleaser.yml` 中的 `brews` 部分
3. **自动化测试**：在多个 macOS 版本上测试
4. **Docker 镜像**：虽然本项目不适合容器化，但可以考虑创建安装脚本镜像

## 参考资料

- [GoReleaser 官方文档](https://goreleaser.com/)
- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [语义化版本规范](https://semver.org/lang/zh-CN/)
