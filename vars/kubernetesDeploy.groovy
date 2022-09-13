import groovy.transform.Field
// import com.sap.piper.Utils
// import static com.sap.piper.Prerequisites.checkScript

@Field String STEP_NAME = getClass().getName()
@Field String METADATA_FILE = 'metadata/kubernetesDeploy.yaml'

void call(Map parameters = [:]) {

    // final script = checkScript(this, parameters) ?: this
    // String stageName = parameters.stageName ?: env.STAGE_NAME

    // def utils = parameters.juStabUtils ?: new Utils()
    // utils.unstashAll(["deployDescriptor"])

    // def utils = parameters.juStabUtils ?: new Utils()
    // utils.unstashAll(["deployDescriptor", "buildDescriptor"])

    List credentials = [
        [type: 'file', id: 'kubeConfigFileCredentialsId', env: ['PIPER_kubeConfig']],
        [type: 'file', id: 'dockerConfigJsonCredentialsId', env: ['PIPER_dockerConfigJSON']],
        [type: 'token', id: 'kubeTokenCredentialsId', env: ['PIPER_kubeToken']],
        [type: 'usernamePassword', id: 'dockerCredentialsId', env: ['PIPER_containerRegistryUser', 'PIPER_containerRegistryPassword']],
        [type: 'token', id: 'githubTokenCredentialsId', env: ['PIPER_githubToken']],
    ]
    piperExecuteBin(parameters, STEP_NAME, METADATA_FILE, credentials)
}
