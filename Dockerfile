FROM golang:1.16-buster AS builder
WORKDIR /project
COPY . ./
RUN cd /project/cmd/kBugger && go build -o /project/bin/ 

FROM registry.access.redhat.com/ubi8/ubi-minimal
EXPOSE 8080
ENV KO_DATA_PATH /kodata
# COPY --from=builder /project/cmd/kBugger/kodata/ ${KO_DATA_PATH}/
COPY --from=builder /project/bin/kBugger /kBugger

ENTRYPOINT ["/kBugger"]

