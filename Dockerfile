FROM centos:7
COPY main2 /root/server
EXPOSE 8080
CMD /root/server