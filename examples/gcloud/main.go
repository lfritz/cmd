package main

import (
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
	top := cmd.NewGroup("gcloud", "Manage Google Cloud Platform resources", "")
	top.Flag("-q --quiet", &quiet, "disable all interactive prompts")

	versionCmd := cmd.New("version", "print version information", "")
	top.Command(versionCmd, func() {
		// ...
	})

	infoCmd := cmd.New("gcloud info", "display information about the environment", "")
	infoCmd.Flag("--anonymize", &info.anonymize, "minimize any personal information")
	infoCmd.Flag("--run-diagnostics", &info.runDiagnostics, "run diagnostics")
	infoCmd.Flag("--show-log", &info.showLog, "print the contents of the last log file")
	top.Command(infoCmd, func() {
		// ...
	})

	runGroup := cmd.NewGroup("gcloud run", "Manage your Cloud Run applications", "")
	runGroup.String("--platform", &run.platform, "platform", "target platform for running commands")
	top.Group(runGroup)

	deployCmd := cmd.New("gcloud run deploy", "Deploy a container to Cloud Run", "")
	deployCmd.String("--image", &run.deploy.image, "IMAGE", "name of the image to deploy")
	deployCmd.Arg("SERVICE", &run.deploy.service)
	runGroup.Command(deployCmd, func() {
		// ...
	})

	top.Run()
}
