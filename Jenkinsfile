#!groovy

@NonCPS
def publish_release(token) {
    def owner = 'cyverse-de'
    def repo = 'process-scanner'
    def releaseName = "build-" + "${env.BUILD_NUMBER}".padLeft(5, "0")

    // Create the release.
    def releaseId = releases.create(token, owner, repo, releaseName)

    // Upload the executable file.
    def artifactName = 'process-scanner-linux-x86_64'
    def fileName = 'process-scanner'
    releases.uploadArtifact(token, owner, repo, releaseId, artifactName, fileName)
}

node('docker') {
    def slackJobDescription = "job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})"
    def container
    try {
        stage "build" {
            checkout scm

            container = "build-${env.BUILD_TAG}"
            sh "docker run --rm --name=${container} -v \$(pwd):/process-scanner -w /process-scanner golang go build ."

            archiveArtifacts artifacts: 'process-scanner', fingerprint: true

            withCredentials([string(credentialsId: 'github-api-token', variable: 'GITHUB_TOKEN')]) {
                publish_release(env.GITHUB_TOKEN)
            }
        }
    } catch (InterruptedException e) {
        currentBuild.result = "ABORTED"
        slackSend color: 'warning', message: "ABORTED: ${slackJobDescription}"
        throw e
    } catch (e) {
        currentBuild.result = "FAILED"
        sh "echo ${e}"
        slackSend color: 'danger', message: "FAILED: ${slackJobDescription}"
        throw e
    } finally {
        sh returnStatus: true, script: "docker rm ${container}"
    }
}
