mkdir build
cp ../cmd/tambola/.env build/.env
go build -o build/server ../cmd/tambola/main.go
docker network create --subnet=172.18.0.0/16 tambola-net
docker build . -t tambola-go
