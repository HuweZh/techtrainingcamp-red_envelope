FROM centos:7
COPY main /root/main
EXPOSE 8080
ENTRYPOINT ["/root/main"]
