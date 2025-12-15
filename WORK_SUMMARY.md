# 工作总结：修复 go install 问题并使用修复版 gosmc

## 问题背景
原项目使用了 `replace` 指令来替换 gosmc 依赖，导致其他用户无法使用 `go install` 安装：
```
replace github.com/panotza/gosmc v1.0.0 => github.com/xykong/gosmc v1.0.2
```

## 已完成的工作
✅ 移除了 `replace` 指令，使项目可以通过 `go install` 安装
✅ 暂时使用原版 `github.com/panotza/gosmc v1.0.0`
✅ 项目已经可以正常编译和运行

## 选定的解决方案：选项 A
Fork iSMC 项目并更新其依赖，这样可以：
- 使用修复版的 gosmc（去除警告）
- 保持 Go modules 最佳实践
- 允许用户使用 `go install` 安装
- 避免 replace 指令

## 需要操作的三个仓库

### 1. github.com/xykong/gosmc
**目的**：修复模块路径声明
**当前状态**：go.mod 中仍声明为 `module github.com/panotza/gosmc`
**需要做的**：
- 将 go.mod 第一行改为 `module github.com/xykong/gosmc`
- 提交并打 tag v1.0.3

### 2. github.com/xykong/iSMC（需要先 Fork）
**目的**：使用修复版的 gosmc
**当前状态**：还未 fork
**需要做的**：
- Fork https://github.com/dkorunic/iSMC
- 将依赖从 `github.com/panotza/gosmc v1.0.0` 改为 `github.com/xykong/gosmc v1.0.3`
- 提交并打 tag v0.7.1-xykong

### 3. github.com/xykong/macos-sensor-exporter
**目的**：使用 fork 的 iSMC
**当前状态**：使用 `github.com/dkorunic/iSMC v0.7.0` 和原版 gosmc
**需要做的**：
- 将依赖改为 `github.com/xykong/iSMC v0.7.1-xykong`
- 更新代码中的导入路径（从 dkorunic/iSMC 改为 xykong/iSMC）
- 验证没有 replace 指令
- 提交更改

## 快速操作命令（按顺序执行）

### 步骤 1：修复 gosmc（在 gosmc 目录）
```bash
cd ~/workspace/xykong/gosmc  # 假设你会把仓库放这里
sed -i '' 's|module github.com/panotza/gosmc|module github.com/xykong/gosmc|' go.mod
git add go.mod
git commit -m "Fix: Update module path to github.com/xykong/gosmc"
git push origin master
git tag v1.0.3
git push origin v1.0.3
```

### 步骤 2：Fork 并修改 iSMC（需要先在 GitHub 上 fork）
```bash
cd ~/workspace/xykong
git clone https://github.com/xykong/iSMC.git
cd iSMC
go get github.com/xykong/gosmc@v1.0.3
go mod tidy
go build ./...
git add go.mod go.sum
git commit -m "Use github.com/xykong/gosmc v1.0.3 (fixes warnings)"
git push origin master
git tag v0.7.1-xykong
git push origin v0.7.1-xykong
```

### 步骤 3：更新 macos-sensor-exporter（在 macos-sensor-exporter 目录）
```bash
cd ~/workspace/xykong/macos-sensor-exporter
sed -i '' 's|github.com/dkorunic/iSMC v0.7.0|github.com/xykong/iSMC v0.7.1-xykong|' go.mod
find . -name "*.go" -type f -exec sed -i '' 's|github.com/dkorunic/iSMC|github.com/xykong/iSMC|g' {} +
rm go.sum
go clean -modcache
go mod tidy
go build
git add .
git commit -m "Switch to forked iSMC with fixed gosmc (removes warnings)"
git push origin master
```

## 验证清单
- [ ] gosmc 的 go.mod 已更新并打了 v1.0.3 tag
- [ ] iSMC 已 fork 并更新依赖，打了 v0.7.1-xykong tag
- [ ] macos-sensor-exporter 已更新为使用 fork 的 iSMC
- [ ] go.mod 中没有 replace 指令
- [ ] 代码可以正常编译
- [ ] 测试 `go install github.com/xykong/macos-sensor-exporter@latest` 能否工作

## 相关文档
- `GOSMC_FIX_GUIDE.md` - 详细的问题说明和解决方案对比
- `IMPLEMENTATION_STEPS.md` - 完整的分步实施指南

## 关键信息
- **gosmc 修复原因**：消除编译警告（IOMasterPort deprecated）
- **为什么要 fork iSMC**：因为 macos-sensor-exporter 依赖 iSMC，而 iSMC 依赖 gosmc。要使用修复版 gosmc 又要避免 replace 指令，就需要 fork iSMC
- **标签命名建议**：
  - gosmc: v1.0.3（新版本）
  - iSMC: v0.7.1-xykong（表示是 xykong 维护的版本）
  - macos-sensor-exporter: 根据你的版本规划

## 注意事项
1. 在新窗口中工作时，确保三个仓库都在同一个父目录下便于管理
2. 按照顺序操作：gosmc → iSMC → macos-sensor-exporter
3. 每个仓库都要确保 tag 已推送到远程
4. 最后可以清理缓存测试 go install 是否正常

## 下一步
1. 在 GitHub 上 fork https://github.com/dkorunic/iSMC
2. 将三个仓库克隆到统一目录
3. 按照上述命令依次执行
4. 有问题随时寻求帮助
