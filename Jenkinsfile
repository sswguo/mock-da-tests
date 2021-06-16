/*
 *  Copyright (C) 2011-2021 Red Hat, Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *          http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

def artifact_glob="build/*"
// def build_image="registry.redhat.io/rhel8/go-toolset:1.15"
def build_image="quay.io/factory2/spmm-jenkins-agent-go-centos7:latest"
// backup build image
// def build_image = "quay.io/app-sre/ubi8-go-toolset:1.15.7"
pipeline {
  agent {
    kubernetes {
      cloud params.JENKINS_AGENT_CLOUD_NAME
      label "jenkins-slave-${UUID.randomUUID().toString()}"
      serviceAccount "jenkins"
      defaultContainer 'jnlp'
      yaml """
      apiVersion: v1
      kind: Pod
      metadata:
        labels:
          app: "jenkins-${env.JOB_BASE_NAME}"
          indy-pipeline-build-number: "${env.BUILD_NUMBER}"
      spec:
        containers:
        - name: jnlp
          image: ${build_image}
          imagePullPolicy: Always
          tty: true
          env:
          - name: HOME
            value: /home/jenkins
          - name: GOROOT
            value: /usr/lib/golang
          - name: GOPATH
            value: /home/jenkins/gopath
          - name: GOPROXY
            value: https://proxy.golang.org
          resources:
            requests:
              memory: 4Gi
              cpu: 2000m
            limits:
              memory: 8Gi
              cpu: 4000m
          workingDir: /home/jenkins
      """
    }
  }
  options {
    //timestamps()
    timeout(time: 120, unit: 'MINUTES')
  }
  environment {
    PIPELINE_NAMESPACE = readFile('/run/secrets/kubernetes.io/serviceaccount/namespace').trim()
    PIPELINE_USERNAME = sh(returnStdout: true, script: 'id -un').trim()
  }
  parameters {
    string(name: 'PNC_REST', defaultValue: '', description: 'Enter the pnc rest url.')
    string(name: 'INDY_URL', defaultValue: '', description: 'Enter the indy url.')
    string(name: 'DA_GROUP', defaultValue: 'DA', description: 'Enter the name of da group.')
    string(name: 'Build_ID', defaultValue: '', description: 'Enter the build id.')
    string(name: 'Concurrent_Goroutines', defaultValue: '9', description: 'Enter the max number of concurrent goroutines.')
    }
  stages {
    stage('Prepare') {
      steps {
        sh 'printenv'
      }
    }

    stage('git checkout') {
      steps{
        script{
          checkout([$class      : 'GitSCM', branches: [[name: 'main']], doGenerateSubmoduleConfigurations: false,
                    extensions  : [[$class: 'CleanCheckout']], submoduleCfg: [],
                    userRemoteConfigs: [[url: 'https://github.com/sswguo/mock-da-tests.git', refspec: '+refs/heads/*:refs/remotes/origin/* +refs/pull/*/head:refs/remotes/origin/pull/*/head']]])
          env.GIT_COMMIT = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()

          echo "Building main commit: ${env.GIT_COMMIT}"
        }
      }
    }

    stage('Build') {
      steps {
        sh 'go build -o main .'
      }
    }

    stage('Run') {
      steps {
        sh "./main ${params.PNC_REST} ${params.INDY_URL} ${params.DA_GROUP} ${params.Build_ID} ${params.Concurrent_Goroutines}"
      }
    }
  }
  post {
    success {
      script {
        echo "SUCCEED"
      }
    }
    failure {
      script {
        echo "FAILED"
      }
    }
  }
}
