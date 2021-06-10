pipeline {
    agent any
    tools {
        go
    }
    stages {
        stage('Running') {
            steps {
                sh 'go run mockdata.go $BUILD_ID'
            }
        }
    }
}