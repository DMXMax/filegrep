package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	homedir "github.com/mitchellh/go-homedir"
	"log"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use: "filegrep",
		Short: "A Grep Executor with the abilty to strip out comments",
		Long: `filegrep seaches for patterns in two stages: 
		first: a general grep to find files of interest
		second: an additional scan with comments removed.
		
		The pattern for finding comments is based on the file extension.`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(scanCmd)
}

func initConfig(){
	if cfgFile != "" {
		//use config from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err !=  nil{
			log.Fatal(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using Config File:", viper.ConfigFileUsed())
	}
}
	
var versionCmd = &cobra.Command{
	Use:	"version",
	Short:	"Get Version Information",
	Long: 	"Grep Current Version",

	Run: func(cmd *cobra.Command, args []string){
		log.Println("Version Information!")
	},
}

