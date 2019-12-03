#!groovy

@Grab(group='io.github.http-builder-ng', module='http-builder-ng-core', version='1.0.4')

import groovyx.net.http.*

def publish_release(token) {
    github = HttpBuilder.configure {
        request.uri = 'https://api.github.com'
        request.accept = ['application/vnd.github.v3+json']
        request.headers['Authorization'] = "token ${token}"
    }

    // Create the new release.
    releaseName: "build-" + env.BUILD_NUMBER.padLeft(5, 0)
    releaseId = github.post {
        request.uri.path = "/repos/cyverse-de/proces-scanner/releases"
        request.contentType = ContentTypes.JSON[0]
        request.body = [
            tag_name: releaseName,
            target_commitish: "master",
            name: releaseName
        ]
    }['id']

    // Upload the executable.
    f = new File("process-scanner")
    github.post {
        request.uri = "https://uploads.github.com/repos/cyverse-de/process-scanner/releases/${releaseId}/assets"
        request.contentType = 'application/octet-stream'
        request.uri.query = [name: "process-scanner-linux-x86_64"]
        request.body = f.bytes
    }
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
