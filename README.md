# charts


[![Build Status](https://travis-ci.org/k8s-community/charts.png?branch=master)](https://travis-ci.org/k8s-community/charts)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](https://github.com/k8s-community/charts/issues)
[![Go Report Card](https://goreportcard.com/badge/github.com/k8s-community/charts)](https://goreportcard.com/report/github.com/k8s-community/charts)

A simple Helm charts server and service template for Kubernetes environment

## Configure Helm

Sets up local helm configuration by reading config (default '~/.kube/config') and using the default context.
Assumed `~/.kube/config` already configured properly and using correct default context.
```sh
helm init
```

## Add charts repository

This command registers specified charts repository served on `https://services.k8s.community/charts`
```sh
helm repo add community-charts https://services.k8s.community/charts
```

## Get charts repository information

Update information on available charts in all registered chart repositories
```sh
helm repo up
```

Show charts repository list
```sh
helm repo list
```

Search any chart in repository
```sh
helm search [keywords]
```

## Release chart

First or manually release of the chart in `default` kubernetes namespace 
```sh
helm install <your-release-name> <absolute path to the chart> --wait
```

Upgrade/install the release `k8sapp-stage` in repository `community-charts` on `stage` kubernetes namespace 
```sh
helm upgrade k8sapp-stage community-charts/k8sapp -i --wait --namespace=stage
```

## Releases information

See all releases with revisions
```sh
helm list
```

## Contributors (unsorted)

- [Igor Dolzhikov](https://github.com/takama)
- [Elena Grahovac](https://github.com/rumyantseva)

All the contributors are welcome. If you would like to be the contributor please accept some rules.
- The pull requests will be accepted only in "develop" branch
- All modifications or additions should be tested

Thank you for your understanding!

## License

[MIT Public License](https://github.com/k8s-community/charts/blob/master/LICENSE)
