pipeline {
	agent any

	environment {
		DOCKER_HUB_CREDENTIALS = credentials(env.DOCKER_HUB_CREDENTIALS_ID ?: 'dockerhub-credentials')
		DOCKER_IMAGE_NAME = env.DOCKER_IMAGE_NAME ?: 'swczk/spray-graphql-api'
		GIT_BRANCH = env.GIT_BRANCH ?: 'main'
		GIT_REPO_URL = env.GIT_REPO_URL ?: 'https://github.com/swczk/spray-graphql-api.git'
		DOCKER_REGISTRY_URL = env.DOCKER_REGISTRY_URL ?: 'https://registry.hub.docker.com'
	}

	stages {
		stage('Checkout') {
			steps {
				git branch: "${GIT_BRANCH}", url: "${GIT_REPO_URL}"
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
					docker.withRegistry("${DOCKER_REGISTRY_URL}", env.DOCKER_HUB_CREDENTIALS_ID ?: 'dockerhub-credentials') {
						docker.image("${DOCKER_IMAGE_NAME}:${env.BUILD_NUMBER}").push()
						docker.image("${DOCKER_IMAGE_NAME}:${env.BUILD_NUMBER}").tag("${DOCKER_IMAGE_NAME}:latest")
						docker.image("${DOCKER_IMAGE_NAME}:latest").push()
					}
				}
			}
		}
	}
}
