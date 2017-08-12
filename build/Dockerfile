FROM golang:1.8.1-alpine

MAINTAINER Spicer Matthews <spicer@cloudmanic.com>

ARG P_UID=0
ARG P_GID=0

# Gots to go Root to make some changes
USER root

# Install bash
RUN apk add --no-cache bash

# Reset Password and group files
RUN set -x \
      && echo "root:x:0:0:root:/root:/bin/ash" > /etc/passwd \
      && echo "root:::0:::::" > /etc/shadow \
      && echo "root:x:0:root" > /etc/group 
      
# Ensure deploy user exists (as our local user)
RUN set -x \ 
  && addgroup -g ${P_GID} -S deploy \
  && adduser -u ${P_UID} -D -S -G deploy deploy 
  
# Create work directory.
RUN set -x \
  && mkdir /work 
   
# Copy the binary over to the image
ADD cache-warm /work/cache-warm
   
# Setup directory and Set perms for app root
RUN set -x \
  && chmod 755 /work \    
  && chown -R deploy:deploy /work  
  
  
# Expose web socket
EXPOSE 8080
  
# This needs to be at the bottom. (back from root)
USER deploy

# Working directory
WORKDIR /work

# App to start
ENTRYPOINT ["./cache-warm"]