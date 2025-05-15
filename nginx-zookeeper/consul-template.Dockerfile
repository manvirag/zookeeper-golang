FROM hashicorp/consul-template:0.30.0

USER root

# Install nginx
RUN apk add --no-cache nginx

# Create directory for nginx configuration
RUN mkdir -p /etc/nginx

# Set the working directory
WORKDIR /etc/consul-template 