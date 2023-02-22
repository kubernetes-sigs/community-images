# `kubectl community-images`

`kubectl` `community-images` is a `kubectl` plugin that displays images running in a Kubernetes cluster that were pulled from community owned repositories and warn the user to switch repositories if needed

## How it Works

The plugin will iterate through readable namespaces, and look for pods. For every pod it can read, the plugin will read the `podspec` for the container images, and any `init` container images. Additionally, it collects the content sha of the image, so that it can be used to disambiguate between different versions pushed with the same tag.

Once there is a list of images, the plugin will print those images that come from a community owned repository and specifically point out those whose repository path have to be updated  

## Quickstart

### Prerequisites

<mark>Note:</mark> You will need [git](https://git-scm.com/downloads) to install the `krew` plugin.

the `community-images` plugin is installed using the `krew` plugin manager for Kubernetes CLI. Installation instructions for `krew` can be found [here](https://krew.sigs.k8s.io/docs/user-guide/setup/install/).

### Installation

After installing & configuring the k8s `krew` plugin, install `community-images` using the following command:

````
$ kubectl krew install community-images
````

### Usage

````
kubectl community-images
````

The community-images is a list of all community owned images, with the most out-of-date images in red.

### Contributing to `community-images`

Find a bug? Want to add a new feature? Want to write docs? Send a [pull request](https://docs.github.com/en/github/collaborating-with-issues-and-pull-requests/about-pull-requests) & we'll review it! 