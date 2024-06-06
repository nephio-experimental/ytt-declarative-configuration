# Copyright 2022 Ericsson AB
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

FROM armdocker.rnd.ericsson.se/dockerhub-ericsson-remote/golang:1.21.5-alpine3.19
ENV CGO_ENABLED=0
WORKDIR /
RUN mkdir local-bin
RUN apk update && apk add bash curl perl-utils
RUN curl -L https://carvel.dev/install.sh | K14SIO_INSTALL_BIN_DIR=local-bin bash
ENV PATH=$PWD/local-bin/:$PATH
COPY ytt_render/src/go.mod ytt_render/src/go.sum ./
WORKDIR /go/src/
RUN go mod download
COPY ytt_render/src .
WORKDIR /go/
COPY parsing/ ./parsing/
RUN go work init
RUN go work use parsing/src
RUN go work use src/
WORKDIR /go/src/
RUN go build -o /usr/local/bin/function ./
# FROM armdocker.rnd.ericsson.se/dockerhub-ericsson-remote/gcr.io/distroless/static:latest
# COPY --from=0 /usr/local/bin/function /usr/local/bin/function
ENTRYPOINT ["function"]