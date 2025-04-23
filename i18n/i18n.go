package i18n

import (
	"os"
	"strings"
)

// Language represents the supported languages
type Language string

const (
	// Chinese language
	Chinese Language = "zh"
	// English language
	English Language = "en"
)

// GetSystemLanguage detects the system language
// Returns Chinese only if the system language starts with "zh", otherwise returns English
func GetSystemLanguage() Language {
	// Get the LANG environment variable
	lang := os.Getenv("LANG")

	// If LANG is not set or doesn't start with "zh", return English
	if lang == "" || !strings.HasPrefix(strings.ToLower(lang), "zh") {
		return English
	}

	return Chinese
}

// Messages contains all the messages for a specific language
type Messages struct {
	// Command descriptions
	CmdShort string
	CmdLong  string

	// Flag descriptions
	FlagHost     string
	FlagPort     string
	FlagUser     string
	FlagPassword string
	FlagDatabase string
	FlagRows     string
	FlagOutput   string
	FlagCompress string

	// User prompts
	PromptPassword string

	// Error messages
	ErrReadPassword     string
	ErrMarkRequiredFlag string

	// Exporter messages
	ExportStart       string
	ExportComplete    string
	ExportFoundTables string
	ExportTableStart  string
	ExportTableRows   string
	ExportCreateZip   string

	// Table structure
	TableStructure string
	ViewStructure  string

	// Table data
	TableData    string
	ViewData     string
	ViewDataNote string

	// Entity types
	EntityTable string
	EntityView  string

	// Error messages for exporter
	ErrConnectDB             string
	ErrPingDB                string
	ErrCreateOutputDir       string
	ErrGetTables             string
	ErrReadTableInfo         string
	ErrCheckTableType        string
	ErrCreateSchemaFile      string
	ErrCreateDataFile        string
	ErrWriteSchemaHeader     string
	ErrWriteDataHeader       string
	ErrGetTableCreateStmt    string
	ErrGetViewCreateStmt     string
	ErrWriteTableStructure   string
	ErrWriteViewStructure    string
	ErrQueryTableData        string
	ErrWriteTableDataComment string
	ErrWriteViewDataComment  string
	ErrGetTableColumns       string
	ErrReadTableColumns      string
	ErrReadTableData         string
	ErrReadViewData          string
	ErrWriteInsertStmt       string
	ErrWriteDataValues       string
	ErrWriteInsertEnd        string
	ErrWriteUnlockTables     string
	ErrCreateZipFile         string
	ErrOpenFile              string
	ErrGetFileInfo           string
	ErrCreateZipHeader       string
	ErrCreateZipWriter       string
	ErrWriteZipContent       string
	ErrWriteSchemaFooter     string
	ErrWriteDataFooter       string
}

// GetMessages returns the messages for the specified language
func GetMessages(lang Language) Messages {
	switch lang {
	case Chinese:
		return chineseMessages
	default:
		return englishMessages
	}
}

// Chinese messages
var chineseMessages = Messages{
	// Command descriptions
	CmdShort: "MySQL数据库导出工具",
	CmdLong:  "MySQL Exporter 是一个用于导出MySQL数据库表结构和数据的工具。\n可以导出指定数据库的所有表结构（包括索引）以及每张表的指定数量数据记录。\n导出的文件可以方便地导入到其他MySQL数据库中。",

	// Flag descriptions
	FlagHost:     "MySQL服务器地址",
	FlagPort:     "MySQL服务器端口",
	FlagUser:     "MySQL用户名",
	FlagPassword: "MySQL密码（如果不提供，将会提示输入）",
	FlagDatabase: "要导出的数据库名",
	FlagRows:     "每张表导出的最大行数",
	FlagOutput:   "输出目录路径",
	FlagCompress: "是否压缩输出文件",

	// User prompts
	PromptPassword: "请输入MySQL密码: ",

	// Error messages
	ErrReadPassword:     "读取密码失败: %w",
	ErrMarkRequiredFlag: "标记必需标志时出错: %v",

	// Exporter messages
	ExportStart:       "开始导出数据库 %s...",
	ExportComplete:    "导出完成!",
	ExportFoundTables: "找到 %d 张表",
	ExportTableStart:  "导出表 %s...",
	ExportTableRows:   "  导出了%s %s 的 %d 行数据",
	ExportCreateZip:   "创建压缩文件 %s...",

	// Table structure
	TableStructure: "-- 表结构 `%s`\nDROP TABLE IF EXISTS `%s`;\n%s;\n\n",
	ViewStructure:  "-- 视图结构 `%s`\nDROP VIEW IF EXISTS `%s`;\n%s;\n\n",

	// Table data
	TableData:    "\n-- 表数据 `%s`\nLOCK TABLES `%s` WRITE;\n",
	ViewData:     "\n-- 视图数据 `%s`\n-- 注意：视图数据仅供参考，不会被导入\n",
	ViewDataNote: "-- 注意：视图数据仅供参考，不会被导入",

	// Entity types
	EntityTable: "表",
	EntityView:  "视图",

	// Error messages for exporter
	ErrConnectDB:             "连接数据库失败: %w",
	ErrPingDB:                "无法连接到数据库: %w",
	ErrCreateOutputDir:       "创建输出目录失败: %w",
	ErrGetTables:             "获取表列表失败: %w",
	ErrReadTableInfo:         "读取表信息失败: %w",
	ErrCheckTableType:        "检查表类型失败: %w",
	ErrCreateSchemaFile:      "创建schema文件失败: %w",
	ErrCreateDataFile:        "创建data文件失败: %w",
	ErrWriteSchemaHeader:     "写入schema文件头部失败: %w",
	ErrWriteDataHeader:       "写入data文件头部失败: %w",
	ErrGetTableCreateStmt:    "获取表 %s 的创建语句失败: %w",
	ErrGetViewCreateStmt:     "获取视图 %s 的创建语句失败: %w",
	ErrWriteTableStructure:   "写入表 %s 的结构失败: %w",
	ErrWriteViewStructure:    "写入视图 %s 的结构失败: %w",
	ErrQueryTableData:        "查询表 %s 的数据失败: %w",
	ErrWriteTableDataComment: "写入表 %s 的数据注释失败: %w",
	ErrWriteViewDataComment:  "写入视图 %s 的数据注释失败: %w",
	ErrGetTableColumns:       "获取%s %s 的列信息失败: %w",
	ErrReadTableColumns:      "读取%s %s 的列信息失败: %w",
	ErrReadTableData:         "读取表 %s 的行数据失败: %w",
	ErrReadViewData:          "读取视图 %s 的行数据失败: %v",
	ErrWriteInsertStmt:       "写入%s %s 的INSERT语句失败: %w",
	ErrWriteDataValues:       "写入%s %s 的数据值失败: %w",
	ErrWriteInsertEnd:        "写入%s %s 的INSERT语句结束符失败: %w",
	ErrWriteUnlockTables:     "写入表 %s 的解锁语句失败: %w",
	ErrCreateZipFile:         "创建zip文件失败: %w",
	ErrOpenFile:              "打开文件 %s 失败: %w",
	ErrGetFileInfo:           "获取文件 %s 信息失败: %w",
	ErrCreateZipHeader:       "创建zip文件头失败: %w",
	ErrCreateZipWriter:       "创建zip writer失败: %w",
	ErrWriteZipContent:       "写入zip文件内容失败: %w",
	ErrWriteSchemaFooter:     "写入schema文件尾部失败: %w",
	ErrWriteDataFooter:       "写入data文件尾部失败: %w",
}

// English messages
var englishMessages = Messages{
	// Command descriptions
	CmdShort: "MySQL Database Export Tool",
	CmdLong:  "MySQL Exporter is a tool for exporting MySQL database table structures and data.\nIt can export all table structures (including indexes) of a specified database and a specified number of data records for each table.\nThe exported files can be easily imported into other MySQL databases.",

	// Flag descriptions
	FlagHost:     "MySQL server address",
	FlagPort:     "MySQL server port",
	FlagUser:     "MySQL username",
	FlagPassword: "MySQL password (if not provided, will prompt for input)",
	FlagDatabase: "Database name to export",
	FlagRows:     "Maximum number of rows to export per table",
	FlagOutput:   "Output directory path",
	FlagCompress: "Whether to compress output files",

	// User prompts
	PromptPassword: "Enter MySQL password: ",

	// Error messages
	ErrReadPassword:     "Failed to read password: %w",
	ErrMarkRequiredFlag: "Error marking required flag: %v",

	// Exporter messages
	ExportStart:       "Starting export of database %s...",
	ExportComplete:    "Export completed!",
	ExportFoundTables: "Found %d tables",
	ExportTableStart:  "Exporting table %s...",
	ExportTableRows:   "  Exported %d rows from %s %s",
	ExportCreateZip:   "Creating zip file %s...",

	// Table structure
	TableStructure: "-- Table structure for `%s`\nDROP TABLE IF EXISTS `%s`;\n%s;\n\n",
	ViewStructure:  "-- View structure for `%s`\nDROP VIEW IF EXISTS `%s`;\n%s;\n\n",

	// Table data
	TableData:    "\n-- Data for table `%s`\nLOCK TABLES `%s` WRITE;\n",
	ViewData:     "\n-- Data for view `%s`\n-- Note: View data is for reference only and will not be imported\n",
	ViewDataNote: "-- Note: View data is for reference only and will not be imported",

	// Entity types
	EntityTable: "table",
	EntityView:  "view",

	// Error messages for exporter
	ErrConnectDB:             "Failed to connect to database: %w",
	ErrPingDB:                "Unable to connect to database: %w",
	ErrCreateOutputDir:       "Failed to create output directory: %w",
	ErrGetTables:             "Failed to get table list: %w",
	ErrReadTableInfo:         "Failed to read table information: %w",
	ErrCheckTableType:        "Failed to check table type: %w",
	ErrCreateSchemaFile:      "Failed to create schema file: %w",
	ErrCreateDataFile:        "Failed to create data file: %w",
	ErrWriteSchemaHeader:     "Failed to write schema file header: %w",
	ErrWriteDataHeader:       "Failed to write data file header: %w",
	ErrGetTableCreateStmt:    "Failed to get CREATE statement for table %s: %w",
	ErrGetViewCreateStmt:     "Failed to get CREATE statement for view %s: %w",
	ErrWriteTableStructure:   "Failed to write structure for table %s: %w",
	ErrWriteViewStructure:    "Failed to write structure for view %s: %w",
	ErrQueryTableData:        "Failed to query data for table %s: %w",
	ErrWriteTableDataComment: "Failed to write data comment for table %s: %w",
	ErrWriteViewDataComment:  "Failed to write data comment for view %s: %w",
	ErrGetTableColumns:       "Failed to get column information for %s %s: %w",
	ErrReadTableColumns:      "Failed to read column information for %s %s: %w",
	ErrReadTableData:         "Failed to read row data for table %s: %w",
	ErrReadViewData:          "Failed to read row data for view %s: %v",
	ErrWriteInsertStmt:       "Failed to write INSERT statement for %s %s: %w",
	ErrWriteDataValues:       "Failed to write data values for %s %s: %w",
	ErrWriteInsertEnd:        "Failed to write INSERT statement end for %s %s: %w",
	ErrWriteUnlockTables:     "Failed to write UNLOCK TABLES statement for table %s: %w",
	ErrCreateZipFile:         "Failed to create zip file: %w",
	ErrOpenFile:              "Failed to open file %s: %w",
	ErrGetFileInfo:           "Failed to get file information for %s: %w",
	ErrCreateZipHeader:       "Failed to create zip file header: %w",
	ErrCreateZipWriter:       "Failed to create zip writer: %w",
	ErrWriteZipContent:       "Failed to write zip file content: %w",
	ErrWriteSchemaFooter:     "Failed to write schema file footer: %w",
	ErrWriteDataFooter:       "Failed to write data file footer: %w",
}

// Current language based on system settings
var currentLanguage = GetSystemLanguage()

// GetCurrentMessages returns the messages for the current system language
func GetCurrentMessages() Messages {
	return GetMessages(currentLanguage)
}
