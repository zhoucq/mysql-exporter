package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zhoucq/mysql-exporter/exporter"
	"golang.org/x/term"
)

var (
	cfgHost     string
	cfgPort     int
	cfgUser     string
	cfgPassword string
	cfgDatabase string
	cfgRows     int
	cfgOutput   string
	cfgCompress bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mysql-exporter",
	Short: "MySQL数据库导出工具",
	Long: `MySQL Exporter 是一个用于导出MySQL数据库表结构和数据的工具。
可以导出指定数据库的所有表结构（包括索引）以及每张表的指定数量数据记录。
导出的文件可以方便地导入到其他MySQL数据库中。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 如果密码为空，提示用户输入密码
		if cfgPassword == "" {
			fmt.Print("请输入MySQL密码: ")
			passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf("读取密码失败: %w", err)
			}
			fmt.Println() // 添加换行符，因为ReadPassword不会自动添加
			cfgPassword = string(passwordBytes)
		}

		config := exporter.Config{
			Host:     cfgHost,
			Port:     cfgPort,
			User:     cfgUser,
			Password: cfgPassword,
			Database: cfgDatabase,
			MaxRows:  cfgRows,
			Output:   cfgOutput,
			Compress: cfgCompress,
		}

		exp, err := exporter.New(config)
		if err != nil {
			return err
		}

		return exp.Execute()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVar(&cfgHost, "host", "localhost", "MySQL服务器地址")
	rootCmd.Flags().IntVar(&cfgPort, "port", 3306, "MySQL服务器端口")
	rootCmd.Flags().StringVar(&cfgUser, "user", "root", "MySQL用户名")
	rootCmd.Flags().StringVar(&cfgPassword, "password", "", "MySQL密码（如果不提供，将会提示输入）")
	rootCmd.Flags().StringVar(&cfgDatabase, "database", "", "要导出的数据库名")
	rootCmd.Flags().IntVar(&cfgRows, "rows", 1000, "每张表导出的最大行数")
	rootCmd.Flags().StringVar(&cfgOutput, "output", "./output", "输出目录路径")
	rootCmd.Flags().BoolVar(&cfgCompress, "compress", true, "是否压缩输出文件")

	if err := rootCmd.MarkFlagRequired("database"); err != nil {
		fmt.Println("标记必需标志时出错:", err)
		os.Exit(1)
	}
}
