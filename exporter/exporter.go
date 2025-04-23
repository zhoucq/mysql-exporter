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
)

// Config 存储导出器的配置信息
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

// Exporter 表示数据库导出器
type Exporter struct {
	config Config
	db     *sql.DB
}

// New 创建一个新的导出器实例
func New(config Config) (*Exporter, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("无法连接到数据库: %w", err)
	}

	return &Exporter{
		config: config,
		db:     db,
	}, nil
}

// Execute 执行导出操作
func (e *Exporter) Execute() error {
	fmt.Printf("开始导出数据库 %s...\n", e.config.Database)

	// 确保输出目录存在
	if err := os.MkdirAll(e.config.Output, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	// 获取所有表
	tables, err := e.getTables()
	if err != nil {
		return err
	}

	fmt.Printf("找到 %d 张表\n", len(tables))

	// 创建schema.sql文件
	schemaPath := filepath.Join(e.config.Output, "schema.sql")
	schemaFile, err := os.Create(schemaPath)
	if err != nil {
		return fmt.Errorf("创建schema文件失败: %w", err)
	}
	defer schemaFile.Close()

	// 写入schema文件头部
	headerComment := fmt.Sprintf("-- MySQL导出 表结构导出\n"+
		"-- 数据库: %s\n"+
		"-- 导出时间: %s\n\n"+
		"SET FOREIGN_KEY_CHECKS=0;\n\n",
		e.config.Database, time.Now().Format("2006-01-02 15:04:05"))
	if _, err := schemaFile.WriteString(headerComment); err != nil {
		return fmt.Errorf("写入schema文件头部失败: %w", err)
	}

	// 创建data.sql文件
	dataPath := filepath.Join(e.config.Output, "data.sql")
	dataFile, err := os.Create(dataPath)
	if err != nil {
		return fmt.Errorf("创建data文件失败: %w", err)
	}
	defer dataFile.Close()

	// 写入data文件头部
	dataHeaderComment := fmt.Sprintf("-- MySQL导出 数据导出\n"+
		"-- 数据库: %s\n"+
		"-- 每张表最多导出 %d 行数据\n"+
		"-- 导出时间: %s\n\n"+
		"SET FOREIGN_KEY_CHECKS=0;\n\n",
		e.config.Database, e.config.MaxRows, time.Now().Format("2006-01-02 15:04:05"))
	if _, err := dataFile.WriteString(dataHeaderComment); err != nil {
		return fmt.Errorf("写入data文件头部失败: %w", err)
	}

	// 导出每张表的结构和数据
	for _, table := range tables {
		fmt.Printf("导出表 %s...\n", table)

		// 导出表结构
		if err := e.exportTableSchema(table, schemaFile); err != nil {
			return err
		}

		// 导出表数据
		if err := e.exportTableData(table, dataFile); err != nil {
			return err
		}
	}

	// 写入文件尾部
	footer := "\nSET FOREIGN_KEY_CHECKS=1;\n"
	if _, err := schemaFile.WriteString(footer); err != nil {
		return fmt.Errorf("写入schema文件尾部失败: %w", err)
	}
	if _, err := dataFile.WriteString(footer); err != nil {
		return fmt.Errorf("写入data文件尾部失败: %w", err)
	}

	// 如果需要压缩，创建zip文件
	if e.config.Compress {
		zipPath := filepath.Join(e.config.Output, "export.zip")
		if err := e.createZipArchive(zipPath, schemaPath, dataPath); err != nil {
			return err
		}
	}

	fmt.Println("导出完成!")
	return nil
}

// TableInfo 存储表的信息
type TableInfo struct {
	Name   string
	IsView bool
}

// getTables 获取数据库中的所有表
func (e *Exporter) getTables() ([]string, error) {
	// 使用information_schema来区分表和视图
	query := fmt.Sprintf("SELECT TABLE_NAME, TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = '%s'", e.config.Database)
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("获取表列表失败: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName, tableType string
		if err := rows.Scan(&tableName, &tableType); err != nil {
			return nil, fmt.Errorf("读取表信息失败: %w", err)
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// isView 检查表是否为视图
func (e *Exporter) isView(tableName string) (bool, error) {
	query := fmt.Sprintf("SELECT TABLE_TYPE FROM information_schema.TABLES WHERE TABLE_SCHEMA = '%s' AND TABLE_NAME = '%s'",
		e.config.Database, tableName)
	var tableType string
	err := e.db.QueryRow(query).Scan(&tableType)
	if err != nil {
		return false, fmt.Errorf("检查表类型失败: %w", err)
	}
	return tableType == "VIEW", nil
}

// exportTableSchema 导出表结构
func (e *Exporter) exportTableSchema(table string, file *os.File) error {
	// 检查是否为视图
	isView, err := e.isView(table)
	if err != nil {
		return err
	}

	// 获取表或视图的CREATE语句
	var tableSchema string
	var query string

	if isView {
		query = fmt.Sprintf("SHOW CREATE VIEW `%s`", table)
		var viewName, characterSet, collation string
		if err := e.db.QueryRow(query).Scan(&viewName, &tableSchema, &characterSet, &collation); err != nil {
			return fmt.Errorf("获取视图 %s 的创建语句失败: %w", table, err)
		}
		// 写入视图结构到文件
		content := fmt.Sprintf("-- 视图结构 `%s`\nDROP VIEW IF EXISTS `%s`;\n%s;\n\n", table, table, tableSchema)
		if _, err := file.WriteString(content); err != nil {
			return fmt.Errorf("写入视图 %s 的结构失败: %w", table, err)
		}
	} else {
		query = fmt.Sprintf("SHOW CREATE TABLE `%s`", table)
		var tableName string
		if err := e.db.QueryRow(query).Scan(&tableName, &tableSchema); err != nil {
			return fmt.Errorf("获取表 %s 的创建语句失败: %w", table, err)
		}
		// 重置自增ID
		tableSchema = resetAutoIncrement(tableSchema)

		// 写入表结构到文件
		content := fmt.Sprintf("-- 表结构 `%s`\nDROP TABLE IF EXISTS `%s`;\n%s;\n\n", table, table, tableSchema)
		if _, err := file.WriteString(content); err != nil {
			return fmt.Errorf("写入表 %s 的结构失败: %w", table, err)
		}
	}

	return nil
}

// exportTableData 导出表数据
func (e *Exporter) exportTableData(table string, file *os.File) error {
	// 检查是否为视图
	isView, err := e.isView(table)
	if err != nil {
		return err
	}

	// 根据是否为视图，使用不同的注释和处理方式
	if isView {
		// 对于视图，只添加注释，不锁定表
		comment := fmt.Sprintf("\n-- 视图数据 `%s`\n-- 注意：视图数据仅供参考，不会被导入\n", table)
		if _, err := file.WriteString(comment); err != nil {
			return fmt.Errorf("写入视图 %s 的数据注释失败: %w", table, err)
		}
	} else {
		// 对于普通表，添加注释并锁定表
		comment := fmt.Sprintf("\n-- 表数据 `%s`\nLOCK TABLES `%s` WRITE;\n", table, table)
		if _, err := file.WriteString(comment); err != nil {
			return fmt.Errorf("写入表 %s 的数据注释失败: %w", table, err)
		}
	}

	// 获取表的所有列
	columns, err := e.getTableColumns(table)
	if err != nil {
		return err
	}

	// 如果没有列，直接返回
	if len(columns) == 0 {
		if !isView {
			// 只有普通表需要解锁
			endComment := fmt.Sprintf("UNLOCK TABLES;\n")
			if _, err := file.WriteString(endComment); err != nil {
				return fmt.Errorf("写入表 %s 的解锁语句失败: %w", table, err)
			}
		}
		return nil
	}

	// 获取表数据
	query := fmt.Sprintf("SELECT * FROM `%s` LIMIT %d", table, e.config.MaxRows)
	rows, err := e.db.Query(query)
	if err != nil {
		// 如果是视图查询失败，记录错误但继续执行
		if isView {
			fmt.Printf("  警告: 无法查询视图 %s 的数据: %v\n", table, err)
			return nil
		}
		return fmt.Errorf("查询表 %s 的数据失败: %w", table, err)
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
				fmt.Printf("  警告: 读取视图 %s 的行数据失败: %v\n", table, err)
				continue
			}
			return fmt.Errorf("读取表 %s 的行数据失败: %w", table, err)
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
			entityType = "视图"
		} else {
			entityType = "表"
		}

		// 如果是新批次的开始，写入完整的INSERT语句
		if batchSize == 0 {
			insertStmt := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)",
				table, columnsList, strings.Join(valueStrings, ", "))
			if _, err := file.WriteString(insertStmt); err != nil {
				return fmt.Errorf("写入%s %s 的INSERT语句失败: %w", entityType, table, err)
			}
		} else {
			// 后续行，只写入值部分
			insertValues := fmt.Sprintf(",\n(%s)", strings.Join(valueStrings, ", "))
			if _, err := file.WriteString(insertValues); err != nil {
				return fmt.Errorf("写入%s %s 的数据值失败: %w", entityType, table, err)
			}
		}

		rowCount++
		batchSize++

		// 如果当前批次已满或者是最后一行，结束当前INSERT语句并开始新批次
		if batchSize >= batchLimit {
			if _, err := file.WriteString(";\n"); err != nil {
				return fmt.Errorf("写入%s %s 的INSERT语句结束符失败: %w", entityType, table, err)
			}
			batchSize = 0 // 重置批次大小
		}
	}

	// 如果有未完成的批次，添加分号结束INSERT语句
	if batchSize > 0 {
		if _, err := file.WriteString(";\n"); err != nil {
			var entityType string
			if isView {
				entityType = "视图"
			} else {
				entityType = "表"
			}
			return fmt.Errorf("写入%s %s 的INSERT语句结束符失败: %w", entityType, table, err)
		}
	}

	// 只有普通表需要解锁
	if !isView {
		endComment := fmt.Sprintf("UNLOCK TABLES;\n")
		if _, err := file.WriteString(endComment); err != nil {
			return fmt.Errorf("写入表 %s 的解锁语句失败: %w", table, err)
		}
	}

	var entityType string
	if isView {
		entityType = "视图"
	} else {
		entityType = "表"
	}
	fmt.Printf("  导出了%s %s 的 %d 行数据\n", entityType, table, rowCount)
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
	entityType := "表"
	if isView {
		entityType = "视图"
	}

	// 查询表或视图的列信息
	query := fmt.Sprintf("SHOW COLUMNS FROM `%s`", table)
	rows, err := e.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("获取%s %s 的列信息失败: %w", entityType, table, err)
	}
	defer rows.Close()

	// 提取列名
	var columns []string
	for rows.Next() {
		var field, typ, null, key, def, extra sql.NullString
		if err := rows.Scan(&field, &typ, &null, &key, &def, &extra); err != nil {
			return nil, fmt.Errorf("读取%s %s 的列信息失败: %w", entityType, table, err)
		}
		columns = append(columns, field.String)
	}

	return columns, nil
}

// createZipArchive 创建zip压缩文件
func (e *Exporter) createZipArchive(zipPath, schemaPath, dataPath string) error {
	fmt.Printf("创建压缩文件 %s...\n", zipPath)

	// 创建zip文件
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("创建zip文件失败: %w", err)
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
		return fmt.Errorf("打开文件 %s 失败: %w", filePath, err)
	}
	defer file.Close()

	// 获取文件信息
	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("获取文件 %s 信息失败: %w", filePath, err)
	}

	// 创建zip文件头
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("创建zip文件头失败: %w", err)
	}

	// 设置压缩方法和文件名
	header.Method = zip.Deflate
	header.Name = zipPath

	// 创建writer
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("创建zip writer失败: %w", err)
	}

	// 复制文件内容到zip
	if _, err := io.Copy(writer, file); err != nil {
		return fmt.Errorf("写入zip文件内容失败: %w", err)
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
