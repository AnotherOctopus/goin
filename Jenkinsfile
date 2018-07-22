node {
        def app
        stage ('build') {
                checkout scm
                sh 'id'
                app = docker.build("goin")
        }
        stage ('setupswarm') {
                sh 'docker network create --subnet=172.18.0.0/16 stalinnet'
                sh 'docker run -d --hostname OriginNode --network stalinnet goin'
                sh 'docker run -d --hostname node1 --network stalinnet --env NETNODE=\'OriginNode\' goin'
                sh 'docker run -d --hostname node2 --network stalinnet --env NETNODE=\'OriginNode\' goin'
                sh 'docker run -d --hostname node3 --network stalinnet --env NETNODE=\'OriginNode\' goin'
                sh 'docker run -d --hostname node4 --network stalinnet --env NETNODE=\'OriginNode\' goin'
                sh 'docker run -d --hostname node5 --network stalinnet --env NETNODE=\'OriginNode\' goin'
                sh 'docker run -d --hostname node6 --network stalinnet --env NETNODE=\'OriginNode\' goin'
        }
        stage('tests') {
                sh 'docker build -f testDockerfile -t testbench .'
                sh 'docker run -d --hostname test --network stalinnet testbench'
        }
        stage('cleanup'){
                sh 'docker stop $(docker ps -a -q)'
                sh 'docker rm $(docker ps -a -q)'
                sh 'docker network rm stalinnet'
        }
        stage ('push') {
                sh 'docker login -u anotheroctopus -p 44Cobr@'
                topush = app.build("anotheroctopus/goin")
                topush.push()
        }

}
