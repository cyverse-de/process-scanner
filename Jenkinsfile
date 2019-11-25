#!groovy
node('docker') {
    slackJobDescription = "job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})"
    try {
        stage "Build"
        buildContainer = "build-${env.BUILD_TAG}"
        sh "docker run --rm --name=${buildContainer} -v $(pwd):/process-scanner -w /process-scanner golang go build ."

        archiveArtifacts artifacts: 'process-scanner', fingerprint: true
    } finally {
        sh returnStatus: true, script: "docker rm ${buildContainer}"
    } catch (InterruptedException e) {
        currentBuild.result = "ABORTED"
        slackSend color: 'warning', message: "ABORTED: ${slackJobDescription}"
        throw e
    } catch (e) {
        currentBuild.result = "FAILED"
        sh "echo ${e}"
        slackSend color: 'danger', message: "FAILED: ${slackJobDescription}"
        throw e
    }
}
