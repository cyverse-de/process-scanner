#!groovy

def publish_release(token) {
    owner = 'cyverse-de'
    repo = 'process-scanner'
    releaseName = "build-" + "${env.BUILD_NUMBER}".padLeft(5, "0")
    echo releaseName

    // Create the release.
    // releaseId = releases.create(token, owner, repo, releaseName)

    // Upload the executable file.
    // artifactName = 'process-scanner-linux-x86_64'
    // artifactFile = new File('process-scanner')
    // releases.uploadArtifact(token, owner, repo, releaseId, artifactName, artifactFile.bytes)
}

node('docker') {
    slackJobDescription = "job '${env.JOB_NAME} [${env.BUILD_NUMBER}]' (${env.BUILD_URL})"
    try {
        stage "Build"
        checkout scm

        buildContainer = "build-${env.BUILD_TAG}"
        sh "docker run --rm --name=${buildContainer} -v \$(pwd):/process-scanner -w /process-scanner golang go build ."

        archiveArtifacts artifacts: 'process-scanner', fingerprint: true

        withCredentials([string(credentialsId: 'github-api-token', variable: 'GITHUB_TOKEN')]) {
            publish_release(env.GITHUB_TOKEN)
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
        sh returnStatus: true, script: "docker rm ${buildContainer}"
    }
}
