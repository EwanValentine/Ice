FROM centurylink/ca-certs
WORKDIR /app
COPY Ice /app/
ENTRYPOINT ["./Ice", "-port", "2000"]
