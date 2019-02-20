go build
docker rm -f zipkin
killall server

docker run --name=zipkin -d -p 9411:9411 openzipkin/zipkin
UPSTREAM_URI=localhost:8888 DOWNSTREAM_URI=http://localhost:8889 METRICS_PORT=8887 SERVICE_NAME=root ./server &
UPSTREAM_URI=localhost:8889 DOWNSTREAM_URI=http://localhost:8890 METRICS_PORT=8886 SERVICE_NAME=branch ./server &
UPSTREAM_URI=localhost:8890                                      METRICS_PORT=8885 SERVICE_NAME=leaf ./server &
