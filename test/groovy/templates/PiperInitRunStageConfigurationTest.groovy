package templates

import org.junit.Before
import org.junit.Rule
import org.junit.Test
import org.junit.rules.ExpectedException
import org.junit.rules.RuleChain
import util.BasePiperTest
import util.JenkinsLoggingRule
import util.JenkinsReadYamlRule
import util.JenkinsStepRule
import util.Rules

import static org.hamcrest.Matchers.*
import static org.junit.Assert.assertThat

class PiperInitRunStageConfigurationTest extends BasePiperTest {
    private JenkinsStepRule jsr = new JenkinsStepRule(this)
    private JenkinsLoggingRule jlr = new JenkinsLoggingRule(this)
    private JenkinsReadYamlRule jryr = new JenkinsReadYamlRule(this)
    private ExpectedException thrown = new ExpectedException()

    @Rule
    public RuleChain rules = Rules
        .getCommonRules(this)
        .around(jryr)
        .around(thrown)
        .around(jlr)
        .around(jsr)

    @Before
    void init()  {

        binding.variables.env.STAGE_NAME = 'Test'

        helper.registerAllowedMethod("findFiles", [Map.class], { map ->
            switch (map.glob) {
                case '**/conf.js':
                    return [new File('conf.js')].toArray()
                case 'myCollection.json':
                    return [new File('myCollection.json')].toArray()
                default:
                    return [].toArray()
            }
        })
    }

 
    @Test
    void testVerboseOption() {
        nullScript.commonPipelineEnvironment.configuration = [
            general: [verbose: true],
            steps: [:],
            stages: [
                Test: [:],
                Integration: [test: 'test'],
                Acceptance: [test: 'test']
            ]
        ]

        jsr.step.piperInitRunStageConfiguration(
            script: nullScript,
            juStabUtils: utils,
            stageConfigResource: 'com.sap.piper/pipeline/stageDefaults.yml'
        )

        assertThat(jlr.log, allOf(
            containsString('[piperInitRunStageConfiguration] Debug - Run Stage Configuration:'),
            containsString('[piperInitRunStageConfiguration] Debug - Run Step Configuration:')
        ))
    }

    @Test
    void testPiperInitDefault() {

        helper.registerAllowedMethod("findFiles", [Map.class], { map -> [].toArray() })

        nullScript.commonPipelineEnvironment.configuration = [
            general: [:],
            steps: [:],
            stages: [
                Test: [:],
                Integration: [test: 'test'],
                Acceptance: [test: 'test']
            ]
        ]

        jsr.step.piperInitRunStageConfiguration(
            script: nullScript,
            juStabUtils: utils,
            stageConfigResource: 'com.sap.piper/pipeline/stageDefaults.yml'
        )

        assertThat(nullScript.commonPipelineEnvironment.configuration.runStage.Acceptance, is(true))
        assertThat(nullScript.commonPipelineEnvironment.configuration.runStage.Integration, is(true))

    }

    @Test
    void testConditionOnlyProductiveBranchOnNonProductiveBranch() {
        helper.registerAllowedMethod('libraryResource', [String.class], {s ->
            if(s == 'testDefault.yml') {
                return '''
stages:
  testStage1:
    stepConditions:
      firstStep:
        filePattern: \'**/conf.js\'
'''
            } else {
                return '''
general: {}
steps: {}
stages:
  testStage1:
    runInAllBranches: false
'''
            }
        })

        binding.variables.env.BRANCH_NAME = 'test'

        jsr.step.piperInitRunStageConfiguration(
            script: nullScript,
            juStabUtils: utils,
            stageConfigResource: 'testDefault.yml',
            productiveBranch: 'master'
        )

        assertThat(nullScript.commonPipelineEnvironment.configuration.runStage.testStage1, is(false))
    }

    @Test
    void testConditionOnlyProductiveBranchOnProductiveBranch() {
        helper.registerAllowedMethod("writeFile", [Map.class], null)
        helper.registerAllowedMethod('libraryResource', [String.class], {s ->
            if(s == 'testDefault.yml') {
                return '''
stages:
  testStage1:
    stepConditions:
      firstStep:
        filePattern: \'**/conf.js\'
'''
            } else {
                return '''
general: {}
steps: {}
stages:
  testStage1:
    runInAllBranches: false
'''
            }
        })

        binding.variables.env.BRANCH_NAME = 'test'

        jsr.step.piperInitRunStageConfiguration(
            script: nullScript,
            juStabUtils: utils,
            stageConfigResource: 'testDefault.yml',
            productiveBranch: 'test'
        )

        assertThat(nullScript.commonPipelineEnvironment.configuration.runStage.testStage1, is(true))
    }

    @Test
    void testStageExtensionExists() {
        helper.registerAllowedMethod('libraryResource', [String.class], {s ->
            if(s == 'testDefault.yml') {
                return '''
stages:
  testStage1:
    extensionExists: true
  testStage2:
    extensionExists: true
  testStage3:
    extensionExists: false
  testStage4:
    extensionExists: 'false'
  testStage5:
    dummy: true
'''
            } else {
                return '''
general:
  projectExtensionsDirectory: './extensions/'
steps: {}
'''
            }
        })

        helper.registerAllowedMethod('fileExists', [String], {path ->
            switch (path) {
                case './extensions/testStage1.groovy':
                    return true
                case './extensions/testStage2.groovy':
                    return false
                case './extensions/testStage3.groovy':
                    return true
                case './extensions/testStage4.groovy':
                    return true
                case './extensions/testStage5.groovy':
                    return true
                default:
                    return false
            }
        })

        nullScript.piperStageWrapper = [:]
        nullScript.piperStageWrapper.allowExtensions = {script -> return true}

        jsr.step.piperInitRunStageConfiguration(
            script: nullScript,
            juStabUtils: utils,
            stageConfigResource: 'testDefault.yml'
        )

        assertThat(nullScript.commonPipelineEnvironment.configuration.runStage.testStage1, is(true))
        assertThat(nullScript.commonPipelineEnvironment.configuration.runStage.testStage2, is(false))
        assertThat(nullScript.commonPipelineEnvironment.configuration.runStage.testStage3, is(false))
        assertThat(nullScript.commonPipelineEnvironment.configuration.runStage.testStage4, is(false))
        assertThat(nullScript.commonPipelineEnvironment.configuration.runStage.testStage5, is(false))
    }
}
