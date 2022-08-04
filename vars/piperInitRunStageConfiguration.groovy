import com.sap.piper.ConfigurationLoader

import static com.sap.piper.Prerequisites.checkScript

import com.sap.piper.ConfigurationHelper
import com.sap.piper.MapUtils
import groovy.transform.Field

@Field String STEP_NAME = getClass().getName()

@Field Set GENERAL_CONFIG_KEYS = [
    /**
     * Print more detailed information into the log.
     * @possibleValues `true`, `false`
     */
    'verbose'
]

@Field Set STEP_CONFIG_KEYS = GENERAL_CONFIG_KEYS.plus([
    /**
     * Defines the library resource that contains the stage configuration settings
     */
    'stageConfigResource'
])

@Field Set PARAMETER_KEYS = STEP_CONFIG_KEYS

void call(Map parameters = [:]) {

    def script = checkScript(this, parameters) ?: this
    String stageName = parameters.stageName ?: env.STAGE_NAME

    script.commonPipelineEnvironment.configuration.runStage = [:]
    script.commonPipelineEnvironment.configuration.runStep = [:]

    // load default & individual configuration
    Map config = ConfigurationHelper.newInstance(this)
        .loadStepDefaults([:], stageName)
        .mixinGeneralConfig(script.commonPipelineEnvironment, GENERAL_CONFIG_KEYS)
        .mixinStepConfig(script.commonPipelineEnvironment, STEP_CONFIG_KEYS)
        .mixinStageConfig(script.commonPipelineEnvironment, stageName, STEP_CONFIG_KEYS)
        .mixin(parameters, PARAMETER_KEYS)
        .withMandatoryProperty('stageConfigResource')
        .use()

    // Go logic to check if the step is active
    String piperGoPath = parameters?.piperGoPath ?: './piper'
    writeFile(file: ".pipeline/stage_conditions.yaml", text: libraryResource(config.stageConfigResource))
    piperExecuteBin.checkIfStepActive(script,piperGoPath,".pipeline/stage_conditions.yaml",".pipeline/step_out.json",".pipeline/stage_out.json","_","_")

    script.commonPipelineEnvironment.configuration.runStage = script.readJSON file: ".pipeline/stage_out.json"
    script.commonPipelineEnvironment.configuration.runStep = script.readJSON file: ".pipeline/step_out.json"

    // Retaining this groovy code as some additional checks for activating-deactivating a stage seems to be done.
    Map stageConfig = ConfigurationHelper.newInstance(this)
            .loadStepDefaults([:], currentStage)
            .mixinStageConfig(script.commonPipelineEnvironment, currentStage)
            .use()

        boolean runStage
        if (stageConfig.runInAllBranches == false && (config.productiveBranch != env.BRANCH_NAME)) {
            runStage = false
        } else if (ConfigurationLoader.stageConfiguration(script, currentStage)) {
            //activate stage if stage configuration is available
            runStage = true
        } else if (stage.getValue().extensionExists == true) {
            runStage =  checkExtensionExists(script, config, currentStage)
        }

        script.commonPipelineEnvironment.configuration.runStage[currentStage] = runStage

    if (config.verbose) {
        echo "[${STEP_NAME}] Debug - Run Stage Configuration: ${script.commonPipelineEnvironment.configuration.runStage}"
        echo "[${STEP_NAME}] Debug - Run Step Configuration: ${script.commonPipelineEnvironment.configuration.runStep}"
    }
}

private static boolean checkExtensionExists(Script script, Map config, String stageName) {
    if (!script.piperStageWrapper.allowExtensions(script)) {
        return false
    }
    // NOTE: These keys exist in "config" if they are configured in the general section of the project
    // config or the defaults. However, in piperStageWrapper, these keys could also be configured for
    // the step "piperStageWrapper" to be effective. Don't know if this should be considered here for consistency.
    if (!config.globalExtensionsDirectory && !config.projectExtensionsDirectory) {
        return false
    }
    def projectInterceptorFile = "${config.projectExtensionsDirectory}${stageName}.groovy"
    def globalInterceptorFile = "${config.globalExtensionsDirectory}${stageName}.groovy"
    return script.fileExists(projectInterceptorFile) || script.fileExists(globalInterceptorFile)
}