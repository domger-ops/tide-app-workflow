# Build and Push Docker Image Workflow

This GitHub Actions workflow automates the process of building a Docker image from the main branch of a repository and pushing it to Docker Hub. Additionally, it updates a Kubernetes manifest file in a separate repository with the new image version.

## Workflow Overview

The workflow is triggered on every push to the `main` branch. It consists of two jobs:

1. **build-and-push**: This job builds the Docker image and pushes it to Docker Hub using the short SHA of the commit as the image tag.
2. **modify-manifest**: This job updates the [`cronjob.yaml`] file in a separate GitHub repository with the new Docker image version.

## Jobs and Steps

### build-and-push

- **Checkout Repository**: Checks out the source code.
- **Set Env**: Sets environment variables, including the short SHA of the commit.
- **Set up Docker Buildx**: Sets up Docker Buildx for building the Docker image.
- **Login to Docker Hub**: Logs into Docker Hub using credentials stored in GitHub secrets.
- **Build and Push Docker Image**: Builds the Docker image with the tag as the short SHA of the commit and pushes it to Docker Hub.

### modify-manifest

- **Checkout Repository**: Checks out the manifest repository specified.
- **Update yaml files**: Updates the [`cronjob.yaml`] file in the checked-out repository with the new Docker image version and pushes the changes.

## Configuration

To use this workflow, you need to configure the following secrets in your GitHub repository:

- `DOCKER_PASS`: Your Docker Hub password.
- `DOCKER_USER`: Your Docker Hub username.
- `GIT_TOKEN`: A GitHub token with permissions to push to the manifest repository.
- `GIT_EMAIL`: The email address to use for commits.

## Usage

To use this workflow, add it to your [`.github/workflows`] directory in your GitHub repository. Ensure you have the required secrets configured in your repository settings.

This workflow automates the build and deployment process, making it easier to maintain and update your Dockerized applications and their Kubernetes manifests.
