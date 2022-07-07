package cmd

import (
	"fmt"
	"path"

	"github.com/SAP/jenkins-library/pkg/kubernetes"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/piperenv"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/versioning"
)

func helmExecute(config helmExecuteOptions, telemetryData *telemetry.CustomData) {
	helmConfig := kubernetes.HelmExecuteOptions{
		AdditionalParameters:      config.AdditionalParameters,
		ChartPath:                 config.ChartPath,
		Image:                     config.Image,
		Namespace:                 config.Namespace,
		KubeContext:               config.KubeContext,
		KeepFailedDeployments:     config.KeepFailedDeployments,
		KubeConfig:                config.KubeConfig,
		HelmDeployWaitSeconds:     config.HelmDeployWaitSeconds,
		DockerConfigJSON:          config.DockerConfigJSON,
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

	helmConfig.DeploymentName = artifactInfo.ArtifactID

	fmt.Printf("\n%v\n", "====== ARTIFACT VERSION ======")
	fmt.Println("artifactInfo.Version", artifactInfo.Version)
	fmt.Println("artifactInfo.ArtifactID", artifactInfo.ArtifactID)
	fmt.Println("artifactInfo.GroupID", artifactInfo.GroupID)
	fmt.Println("artifactInfo.Packaging", artifactInfo.Packaging)

	err = getAndRenderImageInfo(config, GeneralConfig.EnvRootPath, utils)
	if err != nil {
		log.Entry().WithError(err).Fatalf("failed get/render image info: %v", err)
	}

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

func getAndRenderImageInfo(config helmExecuteOptions, rootPath string, utils kubernetes.DeployUtils) error {
	cpe := piperenv.CPEMap{}
	err := cpe.LoadFromDisk(path.Join(rootPath, "commonPipelineEnvironment"))
	if err != nil {
		return fmt.Errorf("failed to load values from commonPipelineEnvironment: %w", err)
	}

	fmt.Println("====== CPE =======")
	fmt.Printf("\n%T, %+v\n\n", cpe, cpe)
	fmt.Printf("\n%T, %+v\n\n", cpe["artifactVersion"], cpe["artifactVersion"])
	fmt.Printf("\n%T, %+v\n\n", cpe["container/imageNames"], cpe["container/imageNames"])
	fmt.Printf("\n%T, %+v\n\n", cpe["custom/nativeBuild"], cpe["custom/nativeBuild"])

	for key, value := range cpe {
		fmt.Printf("\nkey=%v, value=%v, type=%T\n", key, value, value)
	}

	fmt.Println("")

	valuesFiles := []string{}
	defaultValuesFile := fmt.Sprintf("%s/%s", config.ChartPath, "values.yaml")
	defaultValuesFileExists, err := utils.FileExists(defaultValuesFile)
	if err != nil {
		return err
	}

	fmt.Println("====== if statement is started ======")
	if len(config.HelmValues) > 0 {
		fmt.Println("====== case when helmValues > 0 ======")
		fmt.Println("====== defaultValuesFileExists ======", defaultValuesFileExists)
		if defaultValuesFileExists {
			valuesFiles = append(valuesFiles, defaultValuesFile)
		}
		valuesFiles = append(valuesFiles, config.HelmValues...)
	} else {
		if defaultValuesFileExists {
			valuesFiles = append(valuesFiles, defaultValuesFile)
		} else {
			return fmt.Errorf("no one value file is provided, please provide at least one")
			// return fmt.Errorf("no value file to proccess, please provide value file(s)")
		}
	}

	fmt.Println("====== VALUES FILES =======")
	fmt.Printf("\n%+v\n\n", valuesFiles)

	for _, valuesFile := range valuesFiles {
		cpeTemplate, err := utils.FileRead(valuesFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		fmt.Printf("cpeTemplate file: %+v \n", string(cpeTemplate))

		generated, err := cpe.ParseTemplate(string(cpeTemplate))
		if err != nil {
			return fmt.Errorf("failed to parse template: %v", err)
		}

		fmt.Printf("generated: %+v\n", generated.String())

		// tmpl, err := template.New("new").Parse(string(b))
		// if err != nil {
		// 	return fmt.Errorf("failed to parse template: %w", err)
		// }
		// var buf bytes.Buffer
		// err = tmpl.Execute(&buf, params)
		// if err != nil {
		// 	return fmt.Errorf("failed to execute template: %w", err)
		// }
		err = utils.FileWrite(valuesFile, generated.Bytes(), 0700)
		if err != nil {
			return fmt.Errorf("error when updateng file: %w", err)
		}
	}
	return nil
}
