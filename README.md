# MySQL Exporter

A powerful MySQL database export utility written in Go, designed to efficiently export table structures (including indexes) and a configurable number of data records to compressed files for easy migration and testing.

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

## Use Cases

- Prepare test data for development environments
- Create database backups with a limited number of records
- Clone database structures with sample data
- Generate SQL scripts for database schema version control

## CI/CD

This project uses GitHub Actions for continuous integration and deployment:

- **Automated Testing**: All code changes are automatically tested.
- **Multi-platform Builds**: The application is built for multiple platforms (Linux, macOS, Windows) and architectures (amd64, arm64).
- **Automated Releases**: When a new tag is pushed (e.g., `v1.0.0`), a GitHub release is automatically created with pre-built binaries for all supported platforms.

### Development

#### Creating a Release

To create a new release:

1. Tag the commit with a version number:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically build the binaries and create a release with the built artifacts.

## License

MIT
