node {
        def app
        stage ('build') {
                checkout scm
                sh 'id'
                app = docker.build("goin")
        }
        stage ('setupswarm') {
                sh 'docker network create --subnet=172.18.0.0/24 stalinnet'
                sh 'docker run -d --ip 172.18.0.2 --network stalinnet goin'
                sh 'docker run -d --ip 172.18.0.3 --network stalinnet --env NETNODE=\'172.18.0.2\' goin'
                sh 'docker run -d --ip 172.18.0.4 --network stalinnet --env NETNODE=\'172.18.0.2\' goin'
                sh 'docker run -d --ip 172.18.0.5 --network stalinnet --env NETNODE=\'172.18.0.2\' goin'
                sh 'docker run -d --ip 172.18.0.6 --network stalinnet --env NETNODE=\'172.18.0.2\' goin'
                sh 'docker run -d --ip 172.18.0.7 --network stalinnet --env NETNODE=\'172.18.0.2\' goin'
                sh 'docker run -d --ip 172.18.0.8 --network stalinnet --env NETNODE=\'172.18.0.2\' goin'
        }
        stage('tests') {
                sh 'python tests/test1.py'
        }
        stage('cleanup'){
                sh 'docker stop $(docker ps -a -q)'
                sh 'docker rm $(docker ps -a -q)'
                sh 'docker network rm stalinnet'
        }
        stage ('push') {
                sh 'docker login -u anotheroctopus -p 44Cobr@'
                topush = docker.build("anotheroctopus/goin")
                topush.push()
        }

}
