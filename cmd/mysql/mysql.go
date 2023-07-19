/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/19 2:50 下午
 */
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	source string
	destination string
)

// mysqlCmd represents the mysql command
var mysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Realize data migration commands between isomorphic mysql",
	Long: "Realize data migration commands between isomorphic mysql, support single-threaded mode and multi-threaded mode",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mysql called")
		thread := viper.GetBool("toggle")
		if thread {
			fmt.Println("thread is true!")
		} else {
			fmt.Println("thread is false")
		}

		fmt.Println("username is: " + viper.GetString("user"))
	},
}

func init() {

	rootCmd.AddCommand(mysqlCmd)

	mysqlCmd.Flags().StringVarP(&source, "source", "src", "root:root@tcp(localhost:3306)", "source mysql database connnecti")
}