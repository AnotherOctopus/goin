node {
        def app
        tage ('build') {
                checkout scm
                sh 'id'
                app = docker.build("anotheroctopus/goin")
        }
        stage ('push') {
                sh 'docker login -u anotheroctopus -p 44Cobr@'
                app.push()
        }
        stage ('setupswarm') {
                sh 'docker network create --subnet=172.18.0.0/16 stalinnet'
                sh 'docker run -d --hostname OriginNode --network stalinnet goin
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node1 --network stalinnet goin
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node2 --network stalinnet goin
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node3 --network stalinnet goin
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node4 --network stalinnet goin
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node5 --network stalinnet goin
                sh 'docker run -d -e NETNODE='OriginNode' --hostname node6 --network stalinnet goin
        }
                
}
