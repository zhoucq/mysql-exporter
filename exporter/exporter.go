package exporter

import (
	"archive/zip"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zhoucq/mysql-exporter/i18n"
)

// Get the messages for the current language
var msgs = i18n.GetCurrentMessages()

// Config stores the exporter's configuration information
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	MaxRows  int
	Output   string
	Compress bool
}

// Exporter represents the database exporter
type Exporter struct {
	config Config
	db     *sql.DB
}

// New creates a new exporter instance
func New(config Config) (*Exporter, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf(msgs.ErrConnectDB, err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf(msgs.ErrPingDB, err)
	}

	return &Exporter{
		config: config,
		db:     db,
	}, nil
}

// Execute performs the export operation
func (e *Exporter) Execute() error {
	fmt.Printf(msgs.ExportStart+"\n", e.config.Database)

	// Ensure the output directory exists
	if err := os.MkdirAll(e.config.Output, 0755); err != nil {
		return fmt.Errorf(msgs.ErrCreateOutputDir, err)
	}

	// Get all tables
	tables, err := e.getTables()
	if err != nil {
		return err
	}

	fmt.Printf(msgs.ExportFoundTables+"\n", len(tables))

	// Create schema.sql file
	schemaPath := filepath.Join(e.config.Output, "schema.sql")
	schemaFile, err := os.Create(schemaPath)
	if err != nil {
		return fmt.Errorf(msgs.ErrCreateSchemaFile, err)
	}
	defer schemaFile.Close()

	// Write schema file header
	headerComment := fmt.Sprintf("-- MySQL导出 表结构导出\n"+
		"-- 数据库: %s\n"+
		"-- 导出时间: %s\n\n"+
		"SET FOREIGN_KEY_CHECKS=0;\n\n",
		e.config.Database, time.Now().Format("2006-01-02 15:04:05"))
	if _, err := schemaFile.WriteString(headerComment); err != nil {
		return fmt.Errorf(msgs.ErrWriteSchemaHeader, err)
	}

	// Create data.sql file
	dataPath := filepath.Join(e.config.Output, "data.sql")
	dataFile, err := os.Create(dataPath)
	if err != nil {
		return fmt.Errorf(msgs.ErrCreateDataFile, err)
	}
	defer dataFile.Close()

	// Write data file header
	dataHeaderComment := fmt.Sprintf("-- MySQL导出 数据导出\n"+
		"-- 数据库: %s\n"+
		"-- 每张表最多导出 %d 行数据\n"+
		"-- 导出时间: %s\n\n"+
		"SET FOREIGN_KEY_CHECKS=0;\n\n",
		e.config.Database, e.config.MaxRows, time.Now().Format("2006-01-02 15:04:05"))
	if _, err := dataFile.WriteString(dataHeaderComment); err != nil {
		return fmt.Errorf(msgs.ErrWriteDataHeader, err)
	}

	// Export structure and data for each table
	for _, table := range tables {
		fmt.Printf(msgs.ExportTableStart+"\n", table)

		// Export table structure
		if err := e.exportTableSchema(table, schemaFile); err != nil {
			return err
		}

		// Export table data
		if err := e.exportTableData(table, dataFile); err != nil {
			return err
		}
	}

	// Write file footer
	footer := "\nSET FOREIGN_KEY_CHECKS=1;\n"
	if _, err := schemaFile.WriteString(footer); err != nil {
		return fmt.Errorf(msgs.ErrWriteSchemaFooter, err)
	}
	if _, err := dataFile.WriteString(footer); err != nil {
		return fmt.Errorf(msgs.ErrWriteDataFooter, err)
	}

	// If compression is needed, create a zip file
	if e.config.Compress {
		zipPath := filepath.Join(e.config.Output, "export.zip")
		if err := e.createZipArchive(zipPath, schemaPath, dataPath); err != nil {
			return err
		}
	}

	fmt.Println(msgs.ExportComplete)
	return nil
}

// TableInfo stores table information
type TableInfo struct {
	Name   string
	IsView bool
}

// getTables gets all tables in the database
func (e *Exporter) getTables() ([]string, error) {
	// Use information_schema to distinguish between tables and views
	query := fmt.Sprintf("SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = '%s'", e.config.Database)
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(msgs.ErrGetTables, err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName, tableType string
		if err := rows.Scan(&tableName, &tableType); err != nil {
			return nil, fmt.Errorf(msgs.ErrReadTableInfo, err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// isView checks if a table is a view
func (e *Exporter) isView(tableName string) (bool, error) {
	query := fmt.Sprintf("SELECT TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'",
		e.config.Database, tableName)
	var tableType string
	err := e.db.QueryRow(query).Scan(&tableType)
	if err != nil {
		return false, fmt.Errorf(msgs.ErrCheckTableType, err)
	}
	return tableType == "VIEW", nil
}

// exportTableSchema exports the table structure
func (e *Exporter) exportTableSchema(table string, file *os.File) error {
	// Check if it's a view
	isView, err := e.isView(table)
	if err != nil {
		return err
	}

	// Get the CREATE statement for the table or view
	var tableSchema string
	var query string

	if isView {
		query = fmt.Sprintf("SHOW CREATE VIEW `%s`", table)
		var viewName, characterSet, collation string
		if err := e.db.QueryRow(query).Scan(&viewName, &tableSchema, &characterSet, &collation); err != nil {
			return fmt.Errorf(msgs.ErrGetViewCreateStmt, table, err)
		}
		// Write view structure to file
		content := fmt.Sprintf(msgs.ViewStructure, table, table, tableSchema)
		if _, err := file.WriteString(content); err != nil {
			return fmt.Errorf(msgs.ErrWriteViewStructure, table, err)
		}
	} else {
		query = fmt.Sprintf("SHOW CREATE TABLE `%s`", table)
		var tableName string
		if err := e.db.QueryRow(query).Scan(&tableName, &tableSchema); err != nil {
			return fmt.Errorf(msgs.ErrGetTableCreateStmt, table, err)
		}
		// Reset auto-increment ID
		tableSchema = resetAutoIncrement(tableSchema)

		// Write table structure to file
		content := fmt.Sprintf(msgs.TableStructure, table, table, tableSchema)
		if _, err := file.WriteString(content); err != nil {
			return fmt.Errorf(msgs.ErrWriteTableStructure, table, err)
		}
	}

	return nil
}

// exportTableData exports table data
func (e *Exporter) exportTableData(table string, file *os.File) error {
	// Check if it's a view
	isView, err := e.isView(table)
	if err != nil {
		return err
	}

	// Use different comments and processing methods based on whether it's a view
	if isView {
		// For views, only add comments, don't lock the table
		comment := fmt.Sprintf(msgs.ViewData, table)
		if _, err := file.WriteString(comment); err != nil {
			return fmt.Errorf(msgs.ErrWriteViewDataComment, table, err)
		}
	} else {
		// For regular tables, add comments and lock the table
		comment := fmt.Sprintf(msgs.TableData, table, table)
		if _, err := file.WriteString(comment); err != nil {
			return fmt.Errorf(msgs.ErrWriteTableDataComment, table, err)
		}
	}

	// Get all columns of the table
	columns, err := e.getTableColumns(table)
	if err != nil {
		return err
	}

	// If there are no columns, return directly
	if len(columns) == 0 {
		if !isView {
			// Only regular tables need to be unlocked
			endComment := fmt.Sprintf("UNLOCK TABLES;\n")
			if _, err := file.WriteString(endComment); err != nil {
				return fmt.Errorf(msgs.ErrWriteUnlockTables, table, err)
			}
		}
		return nil
	}

	// Get table data
	query := fmt.Sprintf("SELECT * FROM `%s` LIMIT %d", table, e.config.MaxRows)
	rows, err := e.db.Query(query)
	if err != nil {
		// 如果是视图查询失败，记录错误但继续执行
		if isView {
			fmt.Printf("  Warning: %v\n", fmt.Errorf(msgs.ErrReadViewData, table, err))
			return nil
		}
		return fmt.Errorf(msgs.ErrQueryTableData, table, err)
	}
	defer rows.Close()

	// 准备列列表
	columnsList := "`" + strings.Join(columns, "`, `") + "`"

	// 遍历每一行数据
	rowCount := 0
	batchSize := 0
	batchLimit := 1000 // 每批最多1000行

	for rows.Next() {
		// 创建一个动态大小的值切片
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		// 扫描行数据
		if err := rows.Scan(valuePtrs...); err != nil {
			if isView {
				// 如果是视图数据读取失败，记录警告并继续
				fmt.Printf("  Warning: %v\n", fmt.Errorf(msgs.ErrReadViewData, table, err))
				continue
			}
			return fmt.Errorf(msgs.ErrReadTableData, table, err)
		}

		// 处理每个值
		valueStrings := make([]string, len(columns))
		for i, v := range values {
			if v == nil {
				valueStrings[i] = "NULL"
			} else {
				switch value := v.(type) {
				case []byte:
					// 对字符串进行转义
					valueStrings[i] = "'" + escapeString(string(value)) + "'"
				case string:
					valueStrings[i] = "'" + escapeString(value) + "'"
				case time.Time:
					valueStrings[i] = "'" + value.Format("2006-01-02 15:04:05") + "'"
				default:
					valueStrings[i] = fmt.Sprintf("%v", v)
				}
			}
		}

		// 确定实体类型（表或视图）
		var entityType string
		if isView {
			entityType = msgs.EntityView
		} else {
			entityType = msgs.EntityTable
		}

		// 如果是新批次的开始，写入完整的INSERT语句
		if batchSize == 0 {
			insertStmt := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
				table, columnsList, strings.Join(valueStrings, ", "))
			if _, err := file.WriteString(insertStmt); err != nil {
				return fmt.Errorf(msgs.ErrWriteInsertStmt, entityType, table, err)
			}
		} else {
			// 后续行，只写入值部分
			insertValues := fmt.Sprintf(",\n(%s)", strings.Join(valueStrings, ", "))
			if _, err := file.WriteString(insertValues); err != nil {
				return fmt.Errorf(msgs.ErrWriteDataValues, entityType, table, err)
			}
		}

		rowCount++
		batchSize++

		// 如果当前批次已满或者是最后一行，结束当前INSERT语句并开始新批次
		if batchSize >= batchLimit {
			if _, err := file.WriteString(";\n"); err != nil {
				return fmt.Errorf(msgs.ErrWriteInsertEnd, entityType, table, err)
			}
			batchSize = 0 // 重置批次大小
		}
	}

	// 如果有未完成的批次，添加分号结束INSERT语句
	if batchSize > 0 {
		if _, err := file.WriteString(";\n"); err != nil {
			var entityType string
			if isView {
				entityType = msgs.EntityView
			} else {
				entityType = msgs.EntityTable
			}
			return fmt.Errorf(msgs.ErrWriteInsertEnd, entityType, table, err)
		}
	}

	// 只有普通表需要解锁
	if !isView {
		endComment := fmt.Sprintf("UNLOCK TABLES;\n")
		if _, err := file.WriteString(endComment); err != nil {
			return fmt.Errorf(msgs.ErrWriteUnlockTables, table, err)
		}
	}

	var entityType string
	if isView {
		entityType = msgs.EntityView
	} else {
		entityType = msgs.EntityTable
	}
	fmt.Printf(msgs.ExportTableRows+"\n", rowCount, entityType, table)
	return nil
}

// getTableColumns 获取表或视图的所有列名
func (e *Exporter) getTableColumns(table string) ([]string, error) {
	// 检查是否为视图
	isView, err := e.isView(table)
	if err != nil {
		return nil, err
	}

	// 确定实体类型（表或视图）
	entityType := msgs.EntityTable
	if isView {
		entityType = msgs.EntityView
	}

	// 查询表或视图的列信息
	query := fmt.Sprintf("SHOW COLUMNS FROM `%s`", table)
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf(msgs.ErrGetTableColumns, entityType, table, err)
	}
	defer rows.Close()

	// 提取列名
	var columns []string
	for rows.Next() {
		var field, typ, null, key, def, extra sql.NullString
		if err := rows.Scan(&field, &typ, &null, &key, &def, &extra); err != nil {
			return nil, fmt.Errorf(msgs.ErrReadTableColumns, entityType, table, err)
		}
		columns = append(columns, field.String)
	}

	return columns, nil
}

// createZipArchive 创建zip压缩文件
func (e *Exporter) createZipArchive(zipPath, schemaPath, dataPath string) error {
	fmt.Printf(msgs.ExportCreateZip+"\n", zipPath)

	// 创建zip文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf(msgs.ErrCreateZipFile, err)
	}
	defer zipFile.Close()

	// 创建zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 添加schema.sql到zip
	if err := addFileToZip(zipWriter, schemaPath, "schema.sql"); err != nil {
		return err
	}

	// 添加data.sql到zip
	if err := addFileToZip(zipWriter, dataPath, "data.sql"); err != nil {
		return err
	}

	return nil
}

// addFileToZip 将文件添加到zip
func addFileToZip(zipWriter *zip.Writer, filePath, zipPath string) error {
	// 打开源文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf(msgs.ErrOpenFile, filePath, err)
	}
	defer file.Close()

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf(msgs.ErrGetFileInfo, filePath, err)
	}

	// 创建zip文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf(msgs.ErrCreateZipHeader, err)
	}

	// 设置压缩方法和文件名
	header.Method = zip.Deflate
	header.Name = zipPath

	// 创建writer
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf(msgs.ErrCreateZipWriter, err)
	}

	// 复制文件内容到zip
	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf(msgs.ErrWriteZipContent, err)
	}

	return nil
}

// resetAutoIncrement 重置CREATE TABLE语句中的AUTO_INCREMENT值
func resetAutoIncrement(createTableStmt string) string {
	// 使用正则表达式查找并替换AUTO_INCREMENT=数字
	// 这里使用简单的字符串替换方法
	autoIncrIndex := strings.Index(createTableStmt, "AUTO_INCREMENT=")
	if autoIncrIndex == -1 {
		return createTableStmt // 没有找到AUTO_INCREMENT
	}

	// 找到AUTO_INCREMENT=后面的数字结束位置
	endIndex := autoIncrIndex + len("AUTO_INCREMENT=")
	for endIndex < len(createTableStmt) && (createTableStmt[endIndex] >= '0' && createTableStmt[endIndex] <= '9') {
		endIndex++
	}

	// 替换AUTO_INCREMENT值为1
	return createTableStmt[:autoIncrIndex] + "AUTO_INCREMENT=1" + createTableStmt[endIndex:]
}

// escapeString 转义SQL字符串中的特殊字符
func escapeString(s string) string {
	var result strings.Builder
	for _, c := range s {
		switch c {
		case '\'':
			result.WriteString("\\'")
		case '"':
			result.WriteString("\\\"")
		case '\\':
			result.WriteString("\\\\")
		case '\n':
			result.WriteString("\\n")
		case '\r':
			result.WriteString("\\r")
		case '\t':
			result.WriteString("\\t")
		case '\b':
			result.WriteString("\\b")
		case '\f':
			result.WriteString("\\f")
		case '\x00':
			result.WriteString("\\0")
		default:
			result.WriteRune(c)
		}
	}
	return result.String()
}
