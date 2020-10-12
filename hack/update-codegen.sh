#!/usr/bin/env bash

# Copyright 2017 The Kubernetes Authors.
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

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname "${BASH_SOURCE[0]}")/..
CODEGEN_PKG=${CODEGEN_PKG:-$(cd "${SCRIPT_ROOT}"; ls -d -1 ./vendor/k8s.io/code-generator 2>/dev/null || echo ../code-generator)}
#OUTPUT_BASE="$(dirname "${BASH_SOURCE[0]}")/../../.."
OUTPUT_BASE="$(dirname "${BASH_SOURCE[0]}")/.."

# IMPORTANT
# all code needs to be generated with the following PROJECT_MODULE module name
# and will thus be generated in the directory named the same
# move stuff from there manually or modify the script
# cp -r ${PROJECT_MODULE}/pkg/generated pkg/
# cp ${PROJECT_MODULE}/pkg/apis/ndbcontroller/v1alpha1/zz_generated.deepcopy.go \
#     pkg/apis/ndbcontroller/v1alpha1/zz_generated.deepcopy.go
# rm -rf ${PROJECT_MODULE}

# mind the missing slash 
# - listers and informers are not generated 
#   without 100% clean paths (e.g. // in path)
PROJECT_MODULE="github.com/ocklin/ndb-operator"
#PROJECT_MODULE=""


echo "output base is ${OUTPUT_BASE}"
echo "code gen pkg is ${CODEGEN_PKG}"
# generate the code with:
# --output-base    because this script should also be able to run inside the vendor dir of
#                  k8s.io/kubernetes. The output-base is needed for the generators to output into the vendor dir
#                  instead of the $GOPATH directly. For normal projects this can be dropped.
# "deepcopy,client,informer,lister" \
bash "${CODEGEN_PKG}"/generate-groups.sh \
  "deepcopy,client,informer,lister" \
  "${PROJECT_MODULE}/pkg/generated" \
  "${PROJECT_MODULE}/pkg/apis" \
  ndbcontroller:v1alpha1 \
  --output-base  "${OUTPUT_BASE}" \
  --go-header-file "${SCRIPT_ROOT}"/hack/ndb-boilerplate.go.txt

#  "${PROJECT_MODULE}/pkg/apis" \
