pipeline {
    agent any
    stages {
        stage('Running') {
            steps {
                sh 'go run mockdata.go $BUILD_ID'
            }
        }
    }
}