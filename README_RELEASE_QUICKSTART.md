# 快速发布指南 Quick Release Guide

## 🚀 快速发布新版本

```bash
# 1. 更新 CHANGELOG.md 并提交
git add CHANGELOG.md
git commit -m "docs: update changelog for v1.0.0"

# 2. 创建并推送 tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin main
git push origin v1.0.0

# 3. 等待自动构建完成
# 访问 https://github.com/xykong/macos-sensor-exporter/actions
```

## 📦 发布产物

自动构建将创建以下文件：

- `macos-sensor-exporter_v1.0.0_Darwin_x86_64.tar.gz` - Intel Mac
- `macos-sensor-exporter_v1.0.0_Darwin_arm64.tar.gz` - Apple Silicon
- `checksums.txt` - SHA256 校验和

## 🔍 本地测试（可选）

```bash
# 安装 GoReleaser
brew install goreleaser

# 测试构建（不会发布）
goreleaser release --snapshot --clean
```

## 📝 版本号规范

- `v1.0.0` - 正式版本
- `v1.0.0-beta.1` - Beta 版本
- `v1.0.0-rc.1` - Release Candidate

## 📚 完整文档

详细说明请参考 [RELEASE.md](RELEASE.md)
