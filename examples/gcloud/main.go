package main

import (
	"fmt"
	"os"

	"github.com/lfritz/cmd"
)

var (
	quiet bool
	info  struct {
		anonymize, runDiagnostics, showLog bool
	}
	run struct {
		platform string
		deploy   struct {
			image   string
			service string
		}
	}
)

// This example implements a subset of the “gcloud” command’s interface
// (https://cloud.google.com/sdk/gcloud/). It shows how to use groups of commands.
func main() {
	top := cmd.NewGroup("gcloud")
	top.Summary = "Manage Google Cloud Platform resources"
	top.Flag("-q --quiet", &quiet, "disable all interactive prompts")

	versionCmd := top.Command("version", func() {
		fmt.Println("quiet: ", quiet)
	})
	versionCmd.Summary = "print version information"

	infoCmd := top.Command("info", func() {
		fmt.Println("quiet:          ", quiet)
		fmt.Println("anonymize:      ", info.anonymize)
		fmt.Println("runDiagnostics: ", info.runDiagnostics)
		fmt.Println("showLog:        ", info.showLog)
	})
	infoCmd.Summary = "display information about the environment"
	infoCmd.Flag("--anonymize", &info.anonymize, "minimize any personal information")
	infoCmd.Flag("--run-diagnostics", &info.runDiagnostics, "run diagnostics")
	infoCmd.Flag("--show-log", &info.showLog, "print the contents of the last log file")

	runGroup := top.Group("run")
	runGroup.Summary = "Manage your Cloud Run applications"
	runGroup.String("--platform", &run.platform, "platform", "target platform for running commands")

	deployCmd := runGroup.Command("deploy", func() {
		fmt.Println("quiet:    ", quiet)
		fmt.Println("platform: ", run.platform)
		fmt.Println("image:    ", run.deploy.image)
		fmt.Println("service:  ", run.deploy.service)
	})
	deployCmd.Summary = "Deploy a container to Cloud Run"
	deployCmd.String("--image", &run.deploy.image, "IMAGE", "name of the image to deploy")
	deployCmd.Arg("SERVICE", &run.deploy.service)

	top.Run(os.Args[1:])
}
