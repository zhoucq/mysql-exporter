# MySQL Exporter

A MySQL database export tool written in Go that can export table structures and a specified number of data records to compressed files.

[中文文档](README_zh.md)

## Features

- Export structure definitions of all tables in a MySQL database, including table indexes
- Support exporting a specified number of data records for each table
- Generated SQL files can be easily imported into a new database
- Support packaging exported files into compressed format
- Database connection information can be specified via command line parameters

## Installation

```bash
go install github.com/zhoucq/mysql-exporter@latest
```

Or build from source:

```bash
git clone https://github.com/zhoucq/mysql-exporter.git
cd mysql-exporter
go build
```

## Usage

```bash
mysql-exporter --host localhost --port 3306 --user root --password your_password --database your_db --rows 1000 --output ./export
```

### Parameters

| Parameter | Description | Default Value |
|-----------|-------------|---------------|
| `--host` | MySQL server address | localhost |
| `--port` | MySQL server port | 3306 |
| `--user` | Username | root |
| `--password` | Password | - |
| `--database` | Database name to export | - |
| `--rows` | Maximum number of rows to export per table | 1000 |
| `--output` | Output directory path | ./output |
| `--compress` | Whether to compress output files | true |

## Export Format

The exported files will contain the following:

- `schema.sql` - Contains all table structure and index definitions
- `data.sql` - Contains INSERT statements for all table data
- `export.zip` - Contains the above files in a compressed package (when compression is enabled)

## License

MIT
