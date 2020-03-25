# Copyright 2020 Red Hat, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.13 AS builder

ARG SOURCE_REPOSITORY_URL

RUN git clone $SOURCE_REPOSITORY_URL && \
    cd insights-operator-controller && \
    go build 

FROM registry.access.redhat.com/ubi8-minimal

COPY --from=builder /go/insights-operator-controller/insights-operator-controller .

RUN curl -L https://github.com/Yelp/dumb-init/releases/download/v1.2.2/dumb-init_1.2.2_amd64 -o /usr/bin/dumb-init && \
    chmod a+x /usr/bin/dumb-init && \
    chmod a+x /insights-operator-controller

USER 1001

ENTRYPOINT ["dumb-init", "/insights-operator-controller"]
