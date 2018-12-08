# The dockerfile for salmon
FROM alpine:3.8

# copy file
COPY bin/salmon /salmon

# work dir
WORKDIR /

# start salmon
ENTRYPOINT ["/salmon"]
