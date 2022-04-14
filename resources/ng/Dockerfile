FROM nginx:alpine
WORKDIR /usr/share/nginx/html
RUN ls
COPY dist/* .
ENTRYPOINT ["nginx", "-g", "daemon off;"]