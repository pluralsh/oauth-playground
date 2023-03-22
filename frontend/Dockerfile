# Name the node stage "builder"
FROM node:18.12.1-alpine AS builder
# Set working directory
WORKDIR /app
# Copy all files from current directory to working dir in image
COPY . .
# install node modules and build assets
RUN yarn install && yarn build

# nginx state for serving content
FROM nginxinc/nginx-unprivileged:1.21.6-alpine
# Set working directory to nginx asset directory
WORKDIR /usr/share/nginx/html
# Copy static assets from builder stage
COPY --from=builder /app/build .
# Copy nginx.conf
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
# Containers run nginx with global directives and daemon off
ENTRYPOINT ["sh", "-c", "cd /usr/share/nginx/html/ && nginx -g 'daemon off;'"]
# CMD ["sh", "-c", "cd /usr/share/nginx/html/ && nginx -g 'daemon off;'"]