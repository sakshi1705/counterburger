FROM golang:latest 
EXPOSE 8000
RUN mkdir /menu
ADD . /menu/ 
WORKDIR /menu 
ENV GOPATH /menu
RUN cd /menu ; go install menu
CMD ["/menu/bin/menu"]
