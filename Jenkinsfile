node {
        def app
        stage ('build') {
                checkout scm
                sh 'id'
                app = docker.build("anotheroctopus/goin")
        }
        stage ('setupswarm') {
                sh 'docker network create --subnet=172.18.0.0/16 stalinnet'
                sh 'docker run -d --hostname OriginNode --network stalinnet anotheroctopus/goin'
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node1 --network stalinnet anotheroctopus/goin'
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node2 --network stalinnet anotheroctopus/goin'
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node3 --network stalinnet anotheroctopus/goin'
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node4 --network stalinnet anotheroctopus/goin'
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node5 --network stalinnet anotheroctopus/goin'
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node6 --network stalinnet anotheroctopus/goin'
        }
        stage('cleanup'){
                sh 'docker stop $(docker ps -a -q)'
                sh 'docker rm $(docker ps -a -q)'
                sh 'docker network rm stalinnet'
        }
        stage ('push') {
                sh 'docker login -u anotheroctopus -p 44Cobr@'
                app.push()
        }

}
