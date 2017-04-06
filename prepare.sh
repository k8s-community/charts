#!/bin/sh

helm package charts
mv *.tgz packages/
helm repo index packages --url https://containers.golang.services/charts
