package cmd

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/SAP/jenkins-library/pkg/kubernetes"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/versioning"
)

func helmExecute(config helmExecuteOptions, telemetryData *telemetry.CustomData) {
	helmConfig := kubernetes.HelmExecuteOptions{
		ChartPath:                 config.ChartPath,
		Image:                     config.Image,
		Namespace:                 config.Namespace,
		KubeContext:               config.KubeContext,
		KubeConfig:                config.KubeConfig,
		HelmDeployWaitSeconds:     config.HelmDeployWaitSeconds,
		AppVersion:                config.AppVersion,
		Dependency:                config.Dependency,
		PackageDependencyUpdate:   config.PackageDependencyUpdate,
		HelmValues:                config.HelmValues,
		FilterTest:                config.FilterTest,
		DumpLogs:                  config.DumpLogs,
		TargetRepositoryURL:       config.TargetRepositoryURL,
		TargetRepositoryName:      config.TargetRepositoryName,
		TargetRepositoryUser:      config.TargetRepositoryUser,
		TargetRepositoryPassword:  config.TargetRepositoryPassword,
		HelmCommand:               config.HelmCommand,
		CustomTLSCertificateLinks: config.CustomTLSCertificateLinks,
		Version:                   config.Version,
		PublishVersion:            config.Version,
	}

	utils := kubernetes.NewDeployUtilsBundle(helmConfig.CustomTLSCertificateLinks)

	artifactOpts := versioning.Options{
		VersioningScheme: "library",
	}

	artifact, err := versioning.GetArtifact("helm", "", &artifactOpts, utils)
	if err != nil {
		log.Entry().WithError(err).Fatalf("getting artifact information failed: %v", err)
	}
	artifactInfo, err := artifact.GetCoordinates()
	fmt.Printf("\n%v\n\n", artifactInfo)

	fmt.Println("====== Artifact Info ======")
	fmt.Printf("\n%v\n\n", artifactInfo.ArtifactID)
	fmt.Printf("\n%v\n\n", artifactInfo.GroupID)
	fmt.Printf("\n%v\n\n", artifactInfo.Packaging)
	fmt.Printf("\n%v\n\n", artifactInfo.Version)

	helmConfig.DeploymentName = artifactInfo.ArtifactID

	cpe := piperenv.CPEMap{}
	err = cpe.LoadFromDisk(path.Join(GeneralConfig.EnvRootPath, "commonPipelineEnvironment"))
	if err != nil {
		log.Entry().Warning("failed to load values from commonPipelineEnvironment")
	}

	fmt.Println("====== CPE =======")
	fmt.Printf("\n%+v\n\n", cpe)

	file, err := utils.FileRead(config.HelmValues[0])
	if err != nil {
		log.Entry().WithError(err).Fatalf("Error when reading appTemplate '%v'", config.HelmValues[0])
	}

	tmpl, err := template.New("new").Parse(string(file))
	if err != nil {
		log.Entry().Warning("failed to parse template")
	}

	params := struct {
		CPE map[string]interface{}
	}{
		// CPE: map[string]interface{}(cpe),
		CPE: cpe,
	}
	var generatedLdflags bytes.Buffer
	err = tmpl.Execute(&generatedLdflags, params)
	if err != nil {
		log.Entry().Warning("failed to execute template")
	}

	fmt.Println("====== generatedLdflags =======")
	fmt.Printf("\n%v\n\n", generatedLdflags.String())

	err = utils.FileWrite(config.HelmValues[0], generatedLdflags.Bytes(), 0700)
	if err != nil {
		log.Entry().Warning("Error when updating appTemplate")
	}

	valuesFile, err := utils.FileRead(config.HelmValues[0])
	if err != nil {
		log.Entry().WithError(err).Fatalf("Error when reading appTemplate '%v'", config.HelmValues[0])
	}

	fmt.Println("====== valuesFile ======")
	fmt.Printf("\n%v\n\n", string(valuesFile))

	if len(helmConfig.PublishVersion) == 0 {
		helmConfig.PublishVersion = artifactInfo.Version
	}

	helmExecutor := kubernetes.NewHelmExecutor(helmConfig, utils, GeneralConfig.Verbose, log.Writer())

	// error situations should stop execution through log.Entry().Fatal() call which leads to an os.Exit(1) in the end
	if err := runHelmExecute(config, helmExecutor); err != nil {
		log.Entry().WithError(err).Fatalf("step execution failed: %v", err)
	}
}

func runHelmExecute(config helmExecuteOptions, helmExecutor kubernetes.HelmExecutor) error {
	switch config.HelmCommand {
	case "upgrade":
		if err := helmExecutor.RunHelmUpgrade(); err != nil {
			return fmt.Errorf("failed to execute upgrade: %v", err)
		}
	case "lint":
		if err := helmExecutor.RunHelmLint(); err != nil {
			return fmt.Errorf("failed to execute helm lint: %v", err)
		}
	case "install":
		if err := helmExecutor.RunHelmInstall(); err != nil {
			return fmt.Errorf("failed to execute helm install: %v", err)
		}
	case "test":
		if err := helmExecutor.RunHelmTest(); err != nil {
			return fmt.Errorf("failed to execute helm test: %v", err)
		}
	case "uninstall":
		if err := helmExecutor.RunHelmUninstall(); err != nil {
			return fmt.Errorf("failed to execute helm uninstall: %v", err)
		}
	case "dependency":
		if err := helmExecutor.RunHelmDependency(); err != nil {
			return fmt.Errorf("failed to execute helm dependency: %v", err)
		}
	case "publish":
		if err := helmExecutor.RunHelmPublish(); err != nil {
			return fmt.Errorf("failed to execute helm publish: %v", err)
		}
	default:
		if err := runHelmExecuteDefault(config, helmExecutor); err != nil {
			return err
		}
	}

	return nil
}

func runHelmExecuteDefault(config helmExecuteOptions, helmExecutor kubernetes.HelmExecutor) error {
	if err := helmExecutor.RunHelmLint(); err != nil {
		return fmt.Errorf("failed to execute helm lint: %v", err)
	}

	if len(config.Dependency) > 0 {
		if err := helmExecutor.RunHelmDependency(); err != nil {
			return fmt.Errorf("failed to execute helm dependency: %v", err)
		}
	}

	if config.Publish {
		if err := helmExecutor.RunHelmPublish(); err != nil {
			return fmt.Errorf("failed to execute helm publish: %v", err)
		}
	}

	return nil
}
