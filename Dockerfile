FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD latest-linux /latest

EXPOSE 80

ENTRYPOINT ["/latest", "80"]
