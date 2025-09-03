# fastlink

一个快速高效的工具，用于将文件/目录移动到存储位置，并在原位置创建符号链接，帮助您节省磁盘空间的同时保持文件的可访问性。

## 功能特性

- 🚀 **快速文件/目录操作** 使用 [github.com/spf13/fileflow](https://github.com/spf13/fileflow)
- 🔗 **自动符号链接创建** 移动文件后自动创建符号链接
- 📁 **批量操作** 支持多个文件和目录的批量处理
- ⚙️ **灵活配置** 通过 YAML 配置文件进行设置
- 🛡️ **安全检查** 防止递归操作和数据丢失
- 📋 **复制模式** 可选择保留原文件

## 安装方式

### 从源码安装

```bash
git clone https://github.com/xxnuo/fastlink.git
cd fastlink
make install
```

### 手动安装

```bash
go build -o fastlink main.go
sudo cp fastlink /usr/local/bin/fastlink
chmod +x /usr/local/bin/fastlink
```

## 快速开始

1. **创建配置文件**（可选但推荐）：
   ```bash
   mkdir -p ~/.config/fastlink
   cp config.sample.yaml ~/.config/fastlink/config.yaml
   # 编辑配置文件设置默认目标路径
   ```

2. **移动文件并创建符号链接**：
   ```bash
   fastlink /path/to/large-file.zip
   ```

3. **复制模式（保留原文件）**：
   ```bash
   fastlink -k /path/to/important-file.pdf
   ```

## 使用方法

```bash
fastlink [-k|--keep] <源文件> [<目标位置>]
```

### 选项参数

- `-k, --keep`: 保留原文件模式。将文件复制到目标位置而不是移动，并且不在原位置创建符号链接。

### 命令参数

- `<源文件>`: 要处理的文件或目录（必需）
- `<目标位置>`: 目标位置（如果配置文件中已设置则可选）

## 配置说明

在 `~/.config/fastlink/config.yaml` 创建配置文件：

```yaml
# 默认目标目录
# 当没有指定目标目录时，文件将被移动到这里
destination: "/mnt/storage/fastlink"

# 默认保留原文件
# true:  复制文件，保留原文件（不创建符号链接）
# false: 移动文件，创建符号链接（默认行为）
keep: false
```

### 配置优先级

1. 命令行参数（最高优先级）
2. 配置文件设置
3. 默认值（最低优先级）

## 使用示例

### 基础用法

**移动文件到配置的目标位置：**
```bash
fastlink ~/Downloads/large-video.mp4
# 结果：文件移动到 /mnt/storage/fastlink/large-video.mp4
#       在 ~/Downloads/large-video.mp4 创建符号链接
```

**移动目录到指定位置：**
```bash
fastlink ~/Documents/old-projects /backup/archives/
# 结果：目录移动到 /backup/archives/old-projects
#       在 ~/Documents/old-projects 创建符号链接
```

**复制模式（保留原文件）：**
```bash
fastlink --keep ~/important-document.pdf ~/backup/
# 结果：文件复制到 ~/backup/important-document.pdf
#       原文件 ~/important-document.pdf 保持不变
```

### 高级用法

**归档大型目录：**
```bash
# 将多个大型目录移动到外部存储
fastlink ~/Videos/raw-footage /mnt/external/storage/
fastlink ~/Development/old-projects /mnt/external/storage/
fastlink ~/Downloads/iso-files /mnt/external/storage/
```

**备份重要文件：**
```bash
# 复制重要文件但保留原文件
fastlink -k ~/.ssh /backup/ssh-keys/
fastlink -k ~/Documents/contracts /backup/documents/
```

**释放空间同时保持访问：**
```bash
# 移动大型文件但通过符号链接保持可访问
fastlink ~/Downloads/ubuntu-22.04.iso
fastlink ~/Videos/family-vacation-2023
fastlink ~/.cache/large-app-cache
```

### 错误情况

**防止递归移动：**
```bash
fastlink /home/user/documents /home/user/documents/backup
# 错误：不允许递归移动
```

**缺少配置：**
```bash
fastlink /some/file
# 错误：未提供目标位置且配置文件中未找到
```

**目标已存在：**
```bash
fastlink file.txt /backup/
# 如果 /backup/file.txt 已存在：
# 错误：目标已存在: /backup/file.txt
```

## 工作原理

1. **安全检查**：验证源文件存在并防止递归操作
2. **路径解析**：将所有路径转换为绝对路径以确保可靠性
3. **目标准备**：根据需要创建目标目录
4. **文件操作**：
   - **普通模式**：复制文件 → 删除原文件 → 创建符号链接
   - **保留模式**：仅复制文件
5. **符号链接创建**：创建指向新位置的链接

### 符号链接说明

当 fastlink 移动文件时，它会在原位置创建一个指向新位置的符号链接。这意味着：

- ✅ 应用程序仍可使用原路径访问文件
- ✅ 文件看起来仍在原位置
- ✅ 实际文件数据存储在目标位置
- ✅ 您在原位置节省了磁盘空间

示例：
```bash
# 之前
/home/user/large-file.zip (1GB 文件)

# 之后：fastlink /home/user/large-file.zip
/home/user/large-file.zip -> /mnt/storage/fastlink/large-file.zip
# (符号链接)                    (实际 1GB 文件)
```

## 安全特性

- **递归操作防护**：无法将目录移动到自身内部
- **现有文件保护**：不会覆盖目标位置的现有文件
- **符号链接处理**：跳过符号链接的处理以防止循环
- **路径验证**：确保所有路径都有效且可访问

## 故障排除

### 常见问题

**权限被拒绝错误：**
```bash
# 确保您对源文件和目标位置都有写权限
ls -la /path/to/source
ls -la /path/to/destination
```

**符号链接不工作：**
```bash
# 检查符号链接是否存在并指向正确位置
ls -la /original/path
readlink /original/path
```

**找不到配置：**
```bash
# 验证配置文件存在且可读
ls -la ~/.config/fastlink/config.yaml
cat ~/.config/fastlink/config.yaml
```

### 获取帮助

运行不带参数的 `fastlink` 查看使用信息：
```bash
fastlink
```

## 开发指南

### 从源码构建

```bash
git clone https://github.com/xxnuo/fastlink.git
cd fastlink
go mod download
go build -o fastlink main.go
```

### 运行测试

```bash
go test -v ./...
```

### 可用的 Make 目标

```bash
make build     # 构建二进制文件
make test      # 运行测试
make install   # 构建并安装到 /usr/local/bin
make uninstall # 从 /usr/local/bin 删除
```

## 依赖项

- [github.com/spf13/fileflow](https://github.com/spf13/fileflow) - 快速文件操作
- [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) - YAML 配置解析

## 实际使用场景

### 场景1：清理下载文件夹
```bash
# 将大型下载文件移动到外部存储，但保持在下载文件夹中可访问
fastlink ~/Downloads/large-software.iso /mnt/external/downloads/
fastlink ~/Downloads/video-collection /mnt/external/downloads/
```

### 场景2：开发项目归档
```bash
# 将旧项目移动到归档存储，但在开发目录中保持符号链接以便查阅
fastlink ~/Development/legacy-project-v1 /archive/projects/
fastlink ~/Development/experimental-features /archive/projects/
```

### 场景3：媒体文件管理
```bash
# 将处理过的视频文件移动到 NAS，但在本地保持链接
fastlink ~/Videos/edited-content /nas/media-archive/
fastlink ~/Photos/raw-photos /nas/photo-backup/
```

### 场景4：系统缓存优化
```bash
# 将大型应用缓存移动到更大的存储，释放系统盘空间
fastlink ~/.cache/large-app /mnt/cache-storage/
fastlink ~/.local/share/Steam /mnt/games-storage/
```

## 最佳实践

1. **配置默认目标**：设置配置文件以避免每次指定目标路径
2. **定期检查链接**：偶尔验证符号链接仍然有效
3. **备份重要数据**：对重要文件使用 `-k` 选项进行复制而非移动
4. **监控存储空间**：确保目标存储有足够空间
5. **测试配置**：在重要文件上使用前先用测试文件验证配置

## 许可证

本项目为开源项目。请查看仓库了解许可证详情。

## 贡献

欢迎贡献！请随时提交问题和拉取请求。

---

*注意：符号链接在某些文件系统（如 FAT32）上可能不受支持。确保您的目标文件系统支持符号链接功能。*
