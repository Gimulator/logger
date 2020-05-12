FROM ubuntu
COPY ./bin/logger /app/logger
WORKDIR /app
CMD ["./logger", "-ip=localhost:3030", "-config-file=/configs/roles.yaml"]
