# 实施选项 A：Fork iSMC 并更新依赖的详细步骤

## 第一步：修复 gosmc 仓库

### 1.1 检查是否已经克隆了 gosmc
```bash
# 如果还没有克隆，执行：
cd ~/workspace/xykong  # 或其他你存放项目的目录
git clone https://github.com/xykong/gosmc.git
cd gosmc
```

### 1.2 修改 go.mod 文件
```bash
# 将 go.mod 第一行从 "module github.com/panotza/gosmc" 改为：
# module github.com/xykong/gosmc

# 可以使用 sed 命令自动修改：
sed -i '' 's|module github.com/panotza/gosmc|module github.com/xykong/gosmc|' go.mod
```

### 1.3 提交并推送更改
```bash
git add go.mod
git commit -m "Fix: Update module path to github.com/xykong/gosmc"
git push origin master  # 或 main，取决于你的默认分支

# 打新标签
git tag v1.0.3
git push origin v1.0.3
```

## 第二步：Fork 和修改 iSMC 项目

### 2.1 Fork iSMC 项目
访问 https://github.com/dkorunic/iSMC 并点击 Fork 按钮，fork 到你的账号下。

### 2.2 克隆你 fork 的 iSMC 项目
```bash
cd ~/workspace/xykong  # 或其他你存放项目的目录
git clone https://github.com/xykong/iSMC.git
cd iSMC
```

### 2.3 检查当前的 go.mod
```bash
cat go.mod | grep gosmc
# 应该看到类似：github.com/panotza/gosmc v1.0.0
```

### 2.4 更新依赖到你的 gosmc 版本
```bash
# 方法 1：直接编辑 go.mod，将 github.com/panotza/gosmc 改为 github.com/xykong/gosmc v1.0.3
# 然后运行：
go mod tidy

# 方法 2：使用 go get（推荐）
go get github.com/xykong/gosmc@v1.0.3
go mod tidy
```

### 2.5 测试编译
```bash
go build ./...
go test ./...
```

### 2.6 提交并推送
```bash
git add go.mod go.sum
git commit -m "Use github.com/xykong/gosmc v1.0.3 (fixes warnings)"
git push origin master  # 或 main

# 打标签
git tag v0.7.1-xykong
git push origin v0.7.1-xykong
```

## 第三步：更新 macos-sensor-exporter 项目

### 3.1 回到 macos-sensor-exporter 项目
```bash
cd ~/workspace/xykong/macos-sensor-exporter
```

### 3.2 更新 go.mod，使用你 fork 的 iSMC
```bash
# 编辑 go.mod，将：
# github.com/dkorunic/iSMC v0.7.0
# 改为：
# github.com/xykong/iSMC v0.7.1-xykong

# 或使用 sed：
sed -i '' 's|github.com/dkorunic/iSMC v0.7.0|github.com/xykong/iSMC v0.7.1-xykong|' go.mod
```

### 3.3 更新导入路径
```bash
# 需要更新代码中的导入路径
# 将所有的 "github.com/dkorunic/iSMC" 改为 "github.com/xykong/iSMC"

# 检查哪些文件需要修改：
grep -r "dkorunic/iSMC" --include="*.go" .

# 使用 sed 批量替换：
find . -name "*.go" -type f -exec sed -i '' 's|github.com/dkorunic/iSMC|github.com/xykong/iSMC|g' {} +
```

### 3.4 清理并重新构建
```bash
rm go.sum
go clean -modcache
go mod tidy
go build
```

### 3.5 验证
```bash
# 检查 go.mod 中没有 replace 指令
grep "replace" go.mod
# 应该没有输出

# 检查依赖
go list -m all | grep -E "(gosmc|iSMC)"
# 应该看到：
# github.com/xykong/iSMC v0.7.1-xykong
# github.com/xykong/gosmc v1.0.3
```

### 3.6 提交更改
```bash
git add .
git commit -m "Switch to forked iSMC with fixed gosmc (removes warnings)"
git push origin master  # 或你的分支
```

## 验证最终结果

### 模拟其他用户安装
```bash
# 清理本地环境
cd /tmp
go clean -modcache

# 尝试 go install（需要等你推送并打 tag 后）
go install github.com/xykong/macos-sensor-exporter@latest
```

## 注意事项

1. 确保所有的标签都已经推送到 GitHub
2. 如果有 CI/CD 配置，确保它们能正确工作
3. 更新 README.md，说明项目现在使用了 fork 的依赖
4. 考虑向原 iSMC 项目提交 PR，建议他们也使用修复后的 gosmc

## 回滚方案

如果遇到问题，可以执行：
```bash
cd ~/workspace/xykong/macos-sensor-exporter
git revert HEAD
# 或恢复到之前的提交
```
