#!/bin/sh

helm package charts
mv *.tgz packages/
helm repo index packages --url https://services.k8s.community/charts
