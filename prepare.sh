#!/bin/sh

helm package charts
helm package oauth-proxy
mv *.tgz packages/
helm repo index packages --url https://services.k8s.community/charts
