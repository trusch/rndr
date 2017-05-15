// Copyright © 2017 Tino Rusch
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trusch/rndr/renderer"
	"gopkg.in/yaml.v2"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "rndr",
	Short: "render go templates",
	Long:  `rndr renders go templates or folders with templates.`,
	Run: func(cmd *cobra.Command, args []string) {
		src := viper.GetString("src")
		out := viper.GetString("out")
		quiet := viper.GetBool("quiet")
		if quiet {
			f, _ := os.Open(os.DevNull)
			log.SetOutput(f)
		}
		if src == "" {
			log.Fatal("specifiy --src")
		}
		data, err := parseDataFile()
		if err != nil {
			log.Fatal("Error reading data file: ", err)
		}
		renderer := &renderer.Renderer{}
		if err := renderer.Render(src, out, data); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.Flags().StringP("src", "s", "", "source")
	RootCmd.Flags().StringP("out", "o", "", "out")
	RootCmd.Flags().StringP("data", "d", "", "data")
	RootCmd.Flags().BoolP("quiet", "q", false, "quiet")
	viper.BindPFlags(RootCmd.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetEnvPrefix("RNDR_")
	viper.AutomaticEnv() // read in environment variables that match
}

func parseDataFile() (map[string]interface{}, error) {
	dataPath := viper.GetString("data")
	bs, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = yaml.Unmarshal(bs, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
