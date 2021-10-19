FROM golang:1.17.1-buster as build
COPY src/httpserver/ /project/
WORKDIR /project/
RUN go build -o /test/httpserver

FROM ubuntu
COPY --from=build /test/httpserver /test/httpserver
RUN chmod +x /test/httpserver
EXPOSE 5055
ENTRYPOINT [ "/test/httpserver" ]