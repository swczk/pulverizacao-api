pipeline {
	agent any

	environment {
		DOCKER_HUB_CREDENTIALS = credentials('dockerhub-credentials') // Reference the Jenkins credential ID
		DOCKER_IMAGE_NAME = 'swczk/pulveriza-nenem-graphql-api'
	}

	stages {
		stage('Checkout') {
			steps {
				git 'https://github.com/swczk/pulverizacao-api.git'
			}
		}
		stage('Build Docker Image') {
			steps {
				script {
					docker.build("${DOCKER_IMAGE_NAME}:${env.BUILD_NUMBER}")
				}
			}
		}
		stage('Push to Docker Hub') {
			steps {
				script {
					docker.withRegistry('https://registry.hub.docker.com', DOCKER_HUB_CREDENTIALS.id) {
						docker.image("${DOCKER_IMAGE_NAME}:${env.BUILD_NUMBER}").push()
						docker.image("${DOCKER_IMAGE_NAME}:latest").push()
					}
				}
			}
		}
	}
}
