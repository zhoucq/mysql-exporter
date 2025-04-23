# MySQL Exporter

一个用Go编写的MySQL数据库导出工具，可导出表结构和指定数量的数据记录到压缩文件中。

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

## 许可证

MIT
