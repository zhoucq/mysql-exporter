package cmd

import (
	"fmt"
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/zhoucq/mysql-exporter/exporter"
	"github.com/zhoucq/mysql-exporter/i18n"
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

// Get the messages for the current language
var msgs = i18n.GetCurrentMessages()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mysql-exporter",
	Short: msgs.CmdShort,
	Long:  msgs.CmdLong,
	RunE: func(cmd *cobra.Command, args []string) error {
		// If the password is empty, prompt the user to enter a password
		if cfgPassword == "" {
			fmt.Print(msgs.PromptPassword)
			passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
			if err != nil {
				return fmt.Errorf(msgs.ErrReadPassword, err)
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
	rootCmd.Flags().StringVar(&cfgHost, "host", "localhost", msgs.FlagHost)
	rootCmd.Flags().IntVar(&cfgPort, "port", 3306, msgs.FlagPort)
	rootCmd.Flags().StringVar(&cfgUser, "user", "root", msgs.FlagUser)
	rootCmd.Flags().StringVar(&cfgPassword, "password", "", msgs.FlagPassword)
	rootCmd.Flags().StringVar(&cfgDatabase, "database", "", msgs.FlagDatabase)
	rootCmd.Flags().IntVar(&cfgRows, "rows", 1000, msgs.FlagRows)
	rootCmd.Flags().StringVar(&cfgOutput, "output", "./output", msgs.FlagOutput)
	rootCmd.Flags().BoolVar(&cfgCompress, "compress", true, msgs.FlagCompress)

	if err := rootCmd.MarkFlagRequired("database"); err != nil {
		fmt.Printf(msgs.ErrMarkRequiredFlag, err)
		os.Exit(1)
	}
}
