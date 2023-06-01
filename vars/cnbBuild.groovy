import groovy.transform.Field

@Field String STEP_NAME = getClass().getName()
@Field String METADATA_FILE = 'metadata/cnbBuild.yaml'

void call(Map parameters = [:]) {
    List credentials = [[type: 'file', id: 'dockerConfigJsonCredentialsId', env: ['PIPER_dockerConfigJSON']]]
    parameters['securityContext'] = [
        'runAsUser': 0
    ]
    piperExecuteBin(parameters, STEP_NAME, METADATA_FILE, credentials, false, false, true)
}
