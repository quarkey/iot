FROM nginx:alpine
WORKDIR /usr/share/nginx/html
RUN ls
COPY dist/* .
COPY default.conf /etc/nginx/conf.d
ENTRYPOINT ["nginx", "-g", "daemon off;"]