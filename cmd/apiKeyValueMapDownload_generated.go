// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type apiKeyValueMapDownloadOptions struct {
	APIServiceKey   string `json:"apiServiceKey,omitempty"`
	KeyValueMapName string `json:"keyValueMapName,omitempty"`
	DownloadPath    string `json:"downloadPath,omitempty"`
}

// ApiKeyValueMapDownloadCommand Download a specific Key Value Map from the API Portal
func ApiKeyValueMapDownloadCommand() *cobra.Command {
	const STEP_NAME = "apiKeyValueMapDownload"

	metadata := apiKeyValueMapDownloadMetadata()
	var stepConfig apiKeyValueMapDownloadOptions
	var startTime time.Time
	var logCollector *log.CollectorHook
	var splunkClient *splunk.Splunk
	telemetryClient := &telemetry.Telemetry{}

	var createApiKeyValueMapDownloadCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Download a specific Key Value Map from the API Portal",
		Long: `With this step you can download a specific Key Value Map from the API Portal, which returns a zip file with the Key Value Map contents in to current workspace using the OData API.
Learn more about the SAP API Management API for downloading an Key Value Map artifact [here](https://help.sap.com/viewer/66d066d903c2473f81ec33acfe2ccdb4/Cloud/en-US/e26b3320cd534ae4bc743af8013a8abb.html).`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.APIServiceKey)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunkClient = &splunk.Splunk{}
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			if err = log.RegisterANSHookIfConfigured(GeneralConfig.CorrelationID); err != nil {
				log.Entry().WithError(err).Warn("failed to set up SAP Alert Notification Service log hook")
			}

			validation, err := validation.New(validation.WithJSONNamesForStructFields(), validation.WithPredefinedErrorMessages())
			if err != nil {
				return err
			}
			if err = validation.ValidateStruct(stepConfig); err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			stepTelemetryData := telemetry.CustomData{}
			stepTelemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				stepTelemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				stepTelemetryData.ErrorCategory = log.GetErrorCategory().String()
				stepTelemetryData.PiperCommitHash = GitCommit
				telemetryClient.SetData(&stepTelemetryData)
				telemetryClient.Send()
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {

					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.Dsn,
						GeneralConfig.HookConfig.SplunkConfig.Token,
						GeneralConfig.HookConfig.SplunkConfig.Index,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)

					splunkClient.Send(telemetryClient.GetData(), logCollector)

					log.Entry().Debug("Data is sent to dev-instance")
				}

				if len(GeneralConfig.HookConfig.SplunkConfig.ProdDsn) > 0 {

					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.ProdDsn,
						GeneralConfig.HookConfig.SplunkConfig.ProdToken,
						GeneralConfig.HookConfig.SplunkConfig.ProdIndex,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)

					splunkClient.Send(telemetryClient.GetData(), logCollector)

					log.Entry().Debug("Data is sent to prod-instance")
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetryClient.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)

			apiKeyValueMapDownload(stepConfig, &stepTelemetryData)
			stepTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addApiKeyValueMapDownloadFlags(createApiKeyValueMapDownloadCmd, &stepConfig)
	return createApiKeyValueMapDownloadCmd
}

func addApiKeyValueMapDownloadFlags(cmd *cobra.Command, stepConfig *apiKeyValueMapDownloadOptions) {
	cmd.Flags().StringVar(&stepConfig.APIServiceKey, "apiServiceKey", os.Getenv("PIPER_apiServiceKey"), "Service key JSON string to access the API Management Runtime service instance of plan 'api'")
	cmd.Flags().StringVar(&stepConfig.KeyValueMapName, "keyValueMapName", os.Getenv("PIPER_keyValueMapName"), "Specifies the name of the Key Value Map.")
	cmd.Flags().StringVar(&stepConfig.DownloadPath, "downloadPath", os.Getenv("PIPER_downloadPath"), "Specifies Key Value Map download CSV file location.")

	cmd.MarkFlagRequired("apiServiceKey")
	cmd.MarkFlagRequired("keyValueMapName")
	cmd.MarkFlagRequired("downloadPath")
}

// retrieve step metadata
func apiKeyValueMapDownloadMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "apiKeyValueMapDownload",
			Aliases:     []config.Alias{},
			Description: "Download a specific Key Value Map from the API Portal",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "apimApiServiceKeyCredentialsId", Description: "Jenkins secret text credential ID containing the service key to the API Management Runtime service instance of plan 'api'", Type: "jenkins"},
				},
				Parameters: []config.StepParameters{
					{
						Name: "apiServiceKey",
						ResourceRef: []config.ResourceReference{
							{
								Name:  "apimApiServiceKeyCredentialsId",
								Param: "apiServiceKey",
								Type:  "secret",
							},
						},
						Scope:     []string{"PARAMETERS"},
						Type:      "string",
						Mandatory: true,
						Aliases:   []config.Alias{},
						Default:   os.Getenv("PIPER_apiServiceKey"),
					},
					{
						Name:        "keyValueMapName",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_keyValueMapName"),
					},
					{
						Name:        "downloadPath",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_downloadPath"),
					},
				},
			},
		},
	}
	return theMetaData
}
