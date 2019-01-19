FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/
ADD live /latest

EXPOSE 80

ENTRYPOINT ["/live", "80"]
