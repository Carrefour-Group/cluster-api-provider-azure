#!/bin/bash

# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

###############################################################################

# To run locally, set AZURE_CLIENT_ID, AZURE_CLIENT_SECRET, AZURE_SUBSCRIPTION_ID, AZURE_TENANT_ID

set -o errexit
set -o nounset
set -o pipefail

REPO_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
cd "${REPO_ROOT}" || exit 1

# shellcheck source=../hack/ensure-go.sh
source "${REPO_ROOT}/hack/ensure-go.sh"
# shellcheck source=../hack/ensure-kind.sh
source "${REPO_ROOT}/hack/ensure-kind.sh"
# shellcheck source=../hack/ensure-kubectl.sh
source "${REPO_ROOT}/hack/ensure-kubectl.sh"
# shellcheck source=../hack/ensure-kustomize.sh
source "${REPO_ROOT}/hack/ensure-kustomize.sh"
# shellcheck source=../hack/ensure-tags.sh
source "${REPO_ROOT}/hack/ensure-tags.sh"
# shellcheck source=../hack/parse-prow-creds.sh
source "${REPO_ROOT}/hack/parse-prow-creds.sh"

create_cluster() {
    # export cluster template which contains the manifests needed for creating the Azure cluster to run the tests
    if [[ -n ${CI_VERSION:-} || -n ${USE_CI_ARTIFACTS:-} ]]; then
        KUBERNETES_BRANCH="$(cd $(go env GOPATH)/src/k8s.io/kubernetes && git rev-parse --abbrev-ref HEAD)"
        if [[ "${KUBERNETES_BRANCH:-}" =~ "release-" ]]; then
            CI_VERSION_URL="https://dl.k8s.io/ci/latest-${KUBERNETES_BRANCH/release-}.txt"
        else
            CI_VERSION_URL="https://dl.k8s.io/ci/latest.txt"
        fi
        export CLUSTER_TEMPLATE="test/cluster-template-prow-ci-version.yaml"
        export CI_VERSION="${CI_VERSION:-$(curl -sSL ${CI_VERSION_URL})}"
        export KUBERNETES_VERSION="${CI_VERSION}"
    else
        export CLUSTER_TEMPLATE="test/cluster-template-prow.yaml"
    fi

    if [[ "${EXP_MACHINE_POOL:-}" == "true" ]]; then
        export CLUSTER_TEMPLATE="${CLUSTER_TEMPLATE/prow/prow-machine-pool}"
    fi

    export CLUSTER_NAME="capz-$(head /dev/urandom | LC_ALL=C tr -dc a-z0-9 | head -c 6 ; echo '')"
    export AZURE_RESOURCE_GROUP="${CLUSTER_NAME}"
    # Need a cluster with at least 2 nodes
    export CONTROL_PLANE_MACHINE_COUNT="${CONTROL_PLANE_MACHINE_COUNT:-1}"
    export WORKER_MACHINE_COUNT="${WORKER_MACHINE_COUNT:-2}"
    export REGISTRY="${REGISTRY:-capz}"
    export EXP_CLUSTER_RESOURCE_SET="true"
    ${REPO_ROOT}/hack/create-dev-cluster.sh
}

wait_for_nodes() {
    echo "Waiting for ${CONTROL_PLANE_MACHINE_COUNT} control plane machine(s) and ${WORKER_MACHINE_COUNT} worker machine(s) to become Ready"

    # Ensure that all nodes are registered with the API server before checking for readiness
    local total_nodes="$((${CONTROL_PLANE_MACHINE_COUNT} + ${WORKER_MACHINE_COUNT}))"
    while [[ $(kubectl get nodes -ojson | jq '.items | length') -ne "${total_nodes}" ]]; do
        sleep 10
    done

    kubectl wait --for=condition=Ready node --all --timeout=5m
    kubectl get nodes -owide
}

# cleanup all resources we use
cleanup() {
    timeout 1800 kubectl delete cluster "${CLUSTER_NAME}" || true
    make kind-reset || true
}

on_exit() {
    unset KUBECONFIG
    ${REPO_ROOT}/hack/log/log-dump.sh || true
    # cleanup
    if [[ -z "${SKIP_CLEANUP:-}" ]]; then
        cleanup
    fi
}

trap on_exit EXIT
export ARTIFACTS="${ARTIFACTS:-${PWD}/_artifacts}"

# create cluster
if [[ -z "${SKIP_CREATE_WORKLOAD_CLUSTER:-}" ]]; then
    create_cluster
fi

# export the target cluster KUBECONFIG if not already set
export KUBECONFIG="${KUBECONFIG:-${PWD}/kubeconfig}"

export -f wait_for_nodes
timeout --foreground 1800 bash -c wait_for_nodes

if [[ "${#}" -gt 0 ]]; then
    # disable error exit so we can run post-command cleanup
    set +o errexit
    "${@}"
    EXIT_VALUE="${?}"
    exit ${EXIT_VALUE}
fi
