# MySQL Exporter

一个用Go编写的MySQL数据库导出工具，可导出表结构和指定数量的数据记录到压缩文件中。

[English Documentation](README.md)

## 功能特点

- 导出MySQL数据库中所有表的结构定义，包括表索引
- 支持每张表导出指定数量的数据记录
- 生成的SQL文件可以方便地导入到新的数据库中
- 支持将导出的文件打包为压缩格式
- 可以通过命令行参数指定数据库连接信息

## 安装

```bash
go install github.com/zhoucq/mysql-exporter@latest
```

或者从源码构建：

```bash
git clone https://github.com/zhoucq/mysql-exporter.git
cd mysql-exporter
go build
```

## 使用方法

```bash
mysql-exporter --host localhost --port 3306 --user root --password your_password --database your_db --rows 1000 --output ./export
```

### 参数说明

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--host` | MySQL服务器地址 | localhost |
| `--port` | MySQL服务器端口 | 3306 |
| `--user` | 用户名 | root |
| `--password` | 密码 | - |
| `--database` | 要导出的数据库名 | - |
| `--rows` | 每张表导出的最大行数 | 1000 |
| `--output` | 输出目录路径 | ./output |
| `--compress` | 是否压缩输出文件 | true |

## 导出格式

导出的文件将包含以下内容：

- `schema.sql` - 包含所有表结构和索引的定义
- `data.sql` - 包含所有表的数据INSERT语句
- `export.zip` - 包含以上文件的压缩包（当启用压缩时）

## CI/CD

本项目使用GitHub Actions进行持续集成和部署：

- **自动代码检查**：使用golangci-lint检查代码质量。
- **自动测试**：所有代码更改都会自动进行测试。
- **多平台构建**：应用程序会为多个平台（Linux、macOS、Windows）和架构（amd64、arm64）构建。
- **自动发布**：当推送新标签（例如`v1.0.0`）时，GitHub会自动创建一个包含所有支持平台预构建二进制文件的发布版本。

### 创建发布版本

要创建新的发布版本：

1. 使用版本号标记提交：
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions将自动构建二进制文件并创建包含构建产物的发布版本。

## 许可证

MIT
