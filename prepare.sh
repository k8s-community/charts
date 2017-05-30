#!/bin/sh

helm package charts
helm package k8s-community
helm package user-manager
helm package github-integration
mv *.tgz packages/
helm repo index packages --url https://services.k8s.community/charts
