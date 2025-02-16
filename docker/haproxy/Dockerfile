# Use official HAProxy image
FROM haproxy:2.8-alpine

# Switch to root to install packages
USER root

# Install additional tools for health checking
RUN apk add --no-cache curl

# Create required directories with proper permissions
RUN mkdir -p /var/lib/haproxy && \
    chown -R haproxy:haproxy /var/lib/haproxy

# Switch back to haproxy user
USER haproxy

# Expose ports
EXPOSE 80 8404

# Default command with debug mode
CMD ["haproxy", "-f", "/usr/local/etc/haproxy/haproxy.cfg", "-d"]
