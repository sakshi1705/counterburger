FROM golang:latest 
EXPOSE 3000
RUN mkdir /payment
ADD . /payment/ 
WORKDIR /payment 
ENV GOPATH /payment
RUN cd /payment ; go install payment
CMD ["/payment/bin/payment"]
