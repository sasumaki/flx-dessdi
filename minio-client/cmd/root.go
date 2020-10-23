package cmd

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"github.com/theherk/viper"
)

var (
	colorReset = "\033[0m"

	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"

	cfgFile     string
	userLicense string
	rootCmd     = &cobra.Command{
		Use:   "aiga",
		Short: "CLI for deploying ML models.",
		Long:  `Deploy your machine learning models to storage for ezpz deployment`,
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
)

// Execute root command
func Execute() {
	fmt.Println("yolo")
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	loginCmd.PersistentFlags().StringP("endpoint", "e", "", "bucket endpoint")

	deployCmd.PersistentFlags().StringP("file", "f", "", "file path for model")
	deployCmd.MarkPersistentFlagRequired("file")
	viper.BindPFlag("file", deployCmd.PersistentFlags().Lookup("file"))

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(deployCmd)

}

func er(msg interface{}) {
	fmt.Print(string(colorRed))
	fmt.Println("Error:", msg)
	os.Exit(1)
}
func initConfig() {
	configHome, _ := os.UserHomeDir()
	configName := ".aiga"
	configType := "yaml"
	// configPath := filepath.Join(configHome, configName+"."+configType)
	// ----

	viper.AddConfigPath(configHome)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {

	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

}
