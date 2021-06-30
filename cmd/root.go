/*
Copyright © 2020 Marc Wickenden <marc@4armed.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/4armed/kubeshot/internal/config"
	"github.com/4armed/kubeshot/internal/k8sapi"
	"github.com/4armed/kubeshot/internal/screenshot"
	"github.com/kubicorn/kubicorn/pkg/logger"
)

var c = &config.Config{}
var url string
var skipScreenshot bool

var rootCmd = &cobra.Command{
	Use:     "kubeshot",
	Version: config.GitVersion,
	Short:   "Takes screenshots of HTTP services inside (and outside) a Kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("starting")
		urls := []string{}
		var err error

		if len(url) > 0 {
			// A URL has been passed on the command line - this takes precedence
			urls = []string{url}
		} else if len(c.InputFile) > 0 {
			// Process URLs from supplied file, one on each line
			logger.Debug("Processing input file %v", c.InputFile)
			file, err := os.Open(c.InputFile)
			if err != nil {
				logger.Critical("Could not open file %v: %v", c.InputFile, err)
			}

			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				url := scanner.Text()
				logger.Debug("Got URL %v from file", url)
				urls = append(urls, url)
			}
		} else {
			// Fetch from Kubernetes API
			urls, err = k8sapi.GetURLs(c)
			if err != nil {
				logger.Critical("Couldn't retrieve URL info: %v", err)
			}
		}

		if len(urls) == 0 {
			logger.Critical("No URLs to process")
		} else {
			if skipScreenshot {
				logger.Info("Skipping screenshot due to --skip-screenshot")
			} else {
				screenshot.Process(urls, c)
			}
		}
	},
}

// Execute runs the rootCmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&c.GetK8sPods, "pods", "p", false, "fetch URLs for Kubernetes Pods")
	rootCmd.Flags().BoolVarP(&c.GetK8sSvcs, "services", "s", false, "fetch URLs for Kubernetes Services")
	rootCmd.Flags().StringVar(&c.KubeConfig, "kubeconfig", "", "Location of kubeconfig file if not in cluster")
	// rootCmd.Flags().StringVarP(&c.ChromeExe, "chrome-binary", "e", "/headless-shell", "Location of google-chrome binary")
	rootCmd.Flags().BoolVarP(&skipScreenshot, "skip-screenshot", "x", false, "Skip screenshots, just build URL list")
	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL to screenshot")
	rootCmd.Flags().StringVarP(&c.InputFile, "file", "f", "", "File containing URLs to process, each on a separate line")
	rootCmd.Flags().StringVarP(&c.OutputDir, "directory", "d", "", "Directory to write output to")
	rootCmd.Flags().IntVarP(&c.NumberOfWorkers, "workers", "w", 3, "number of concurrent workers")

	rootCmd.PersistentFlags().IntVarP(&logger.Level, "verbose", "v", 3, "set log level, use 0 to silence, 4 for debugging")
	rootCmd.PersistentFlags().BoolVarP(&logger.Color, "color", "C", true, "toggle colorised logs")
}
