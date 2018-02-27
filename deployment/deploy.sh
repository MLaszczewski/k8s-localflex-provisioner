#!/bin/bash

set -euo pipefail

function main() {
        kubectl create -f $(dirname $(readlink --canonicalize-existing "$0"))/clusterrole.yaml
        kubectl create -f $(dirname $(readlink --canonicalize-existing "$0"))/clusterrolebinding.yaml
        kubectl create -f $(dirname $(readlink --canonicalize-existing "$0"))/role.yaml
        kubectl create -f $(dirname $(readlink --canonicalize-existing "$0"))/rolebinding.yaml
        kubectl create -f $(dirname $(readlink --canonicalize-existing "$0"))/serviceaccount.yaml
	kubectl create -f $(dirname $(readlink --canonicalize-existing "$0"))/daemonset.yaml
	kubectl create -f $(dirname $(readlink --canonicalize-existing "$0"))/deployment.yaml
	kubectl create -f $(dirname $(readlink --canonicalize-existing "$0"))/storageclass.yaml
	wait localflex-deploy
	echo -n "deleting \"localflex-deploy\" daemonset"
	echo
	kubectl delete -f $(dirname $(readlink --canonicalize-existing "$0"))/daemonset.yaml
}

function wait() {
	# need some time for pods to show up
	#sleep 10

	PODS=$(kubectl get pods -n kube-system | grep $1 | awk '{print $1}')

	echo -n "waiting for \"$1\" pods to run"

	for POD in ${PODS}; do
		while [[ $(kubectl get pod ${POD} -n kube-system -o go-template --template "{{.status.phase}}") != "Running" ]]; do
			echo -n "."
			sleep 1
		done
	done

	echo
	echo -n "waiting for \"$1\" daemonset to complete"

	for POD in ${PODS}; do
		while [[ $(kubectl logs ${POD} -n kube-system --tail 1) != "done" ]]; do
			echo -n "."
			sleep 1
		done
	done

	echo
}

main
