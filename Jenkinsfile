pipeline {
    agent { docker { image 'golang' } }
    stages {
        stage('Build') {
            steps {
                sh 'go build cmd/mockdatests/main.go'
            }
        }
        stage('Running') {
            steps {
                sh './main $BUILD_ID'
            }
        }
    }
}
