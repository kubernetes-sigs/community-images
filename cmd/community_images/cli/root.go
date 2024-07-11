/*
Copyright 2023 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cli

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/kubernetes-sigs/community-images/pkg/community_images"
	"github.com/kubernetes-sigs/community-images/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	spin "github.com/tj/go-spin"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

var (
	KubernetesConfigFlags *genericclioptions.ConfigFlags
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "community-images",
		Short:         "",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()
			if v.GetBool("mirror") {
				// regexp for all community registries
				re := regexp.MustCompile(`^k8s\.gcr\.io/|^gcr\.io/google-containers|^registry\.k8s\.io/`)
				return runParsableCommand(v, re)
			} else {
				if v.GetBool("plain") {
					// regexp for only the old registries that we need to move folks off of
					re := regexp.MustCompile(`^k8s\.gcr\.io/|^gcr\.io/google-containers`)
					return runParsableCommand(v, re)
				}
				return runPrettyCommand(v)
			}
		},
	}

	cobra.OnInitialize(initConfig)

	KubernetesConfigFlags = genericclioptions.NewConfigFlags(false)
	KubernetesConfigFlags.AddFlags(cmd.Flags())

	cmd.Flags().StringSlice("include-ns", []string{}, "optional list of namespaces to include in the searching")
	cmd.Flags().StringSlice("ignore-ns", []string{}, "optional list of namespaces to exclude from searching")
	cmd.Flags().Bool("plain", false, "machine parsable output (list of images from older registries ONLY)")
	cmd.Flags().Bool("mirror", false, "list of images that should be mirrored from all community registries (assumes --plain is true)")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetEnvPrefix("OUTDATED")
	viper.AutomaticEnv()
}

func runParsableCommand(v *viper.Viper, re *regexp.Regexp) error {
	imagesList, err := community_images.ListImages(KubernetesConfigFlags, nil, v.GetStringSlice("ignore-ns"), v.GetStringSlice("include-ns"))
	if err != nil {
		os.Exit(1)
		return nil
	}

	for _, runningImage := range imagesList {
		image := imageWithTagPlain(runningImage)
		if re.MatchString(image) {
			fmt.Println(image)
		}
	}
	return nil
}

func runPrettyCommand(v *viper.Viper) error {
	log := logger.NewLogger()
	log.Info("")

	s := spin.New()
	finishedCh := make(chan bool, 1)
	foundImageName := make(chan string, 1)
	go func() {
		lastImageName := ""
		for {
			select {
			case <-finishedCh:
				fmt.Printf("\r")
				return
			case i := <-foundImageName:
				lastImageName = i
			case <-time.After(time.Millisecond * 100):
				if lastImageName == "" {
					fmt.Printf("\r  \033[36mSearching for images\033[m %s", s.Next())
				} else {
					fmt.Printf("\r  \033[36mSearching for images\033[m %s (%s)", s.Next(), lastImageName)
				}
			}
		}
	}()
	defer func() {
		finishedCh <- true
	}()

	imagesList, err := community_images.ListImages(KubernetesConfigFlags, foundImageName, v.GetStringSlice("ignore-ns"), v.GetStringSlice("include-ns"))
	if err != nil {
		log.Error(err)
		log.Info("")
		os.Exit(1)
		return nil
	}
	finishedCh <- true

	config, _ := KubernetesConfigFlags.ToRESTConfig()
	log.Header(headerLine(config.Host))
	re := regexp.MustCompile(`^k8s\.gcr\.io/|^gcr\.io/google[-_]containers`)
	for _, runningImage := range imagesList {
		image := imageWithTag(runningImage)
		log.StartImageLine(image)
		if re.MatchString(image) {
			log.ImageRedLine(image)
		} else {
			log.ImageGreenLine(image)
		}
	}

	fmt.Printf("\nImages in \033[91mred ❌ \033[mare being pulled from \033[1m*outdated*\033[0m Kubernetes community registries.\n" +
		"The others marked in \033[92mgreen ✅ \u001B[mare good as they do not use the outdated registries.\n" +
		"Please copy these images to your own registry and change your manifest(s)\nto point to the new location.\n\n")
	fmt.Printf(
		"If you are unable to do so, as a short term fix please use \033[92m`registry.k8s.io`\033[m " +
			"\ninstead of \033[91m`k8s.gcr.io`\033[m until you have your own registry.\n\n")
	fmt.Printf("This simple change on your part will help the Kubernetes community immensely as it\n" +
		"reduces the cost of us serving these container images.\n")

	fmt.Printf("\n\033[1mWhy you should do this as soon as possible? Read more in the following blog\n" +
		"posts by the Kubernetes community:\033[m\n" +
		"- https://kubernetes.io/blog/2022/11/28/registry-k8s-io-faster-cheaper-ga/\n" +
		"- https://kubernetes.io/blog/2023/02/06/k8s-gcr-io-freeze-announcement/\n")

	log.Info("")
	return nil
}
