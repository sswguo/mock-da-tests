pipeline {
    agent { docker { image 'golang' } }
    stages {
        stage('Running') {
            steps {
                sh 'go run mockdata.go $BUILD_ID'
            }
        }
    }
}
