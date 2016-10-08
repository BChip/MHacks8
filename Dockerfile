FROM scratch

COPY bin/main /

EXPOSE 8080

CMD ["/main"]