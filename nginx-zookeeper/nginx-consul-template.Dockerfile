FROM hashicorp/consul-template:0.30.0

# Switch to root user to install packages
USER root

# Install nginx
RUN apk add --no-cache nginx

# Create necessary directories for nginx
RUN mkdir -p /etc/nginx/conf.d /run/nginx

# Make sure nginx runs in the foreground (no daemon)
RUN echo "daemon off;" >> /etc/nginx/nginx.conf

# Switch back to the default user (if needed)
USER consul-template
