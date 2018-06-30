pipeline {
        agent { docker { image 'go:1.8.3' } }
        stages {
                stage('build') {
                        steps {
                                sh 'go version'
                        }
                }
        }
}
