Image Repo
----

## Run with Docker

```shell
# Pull the Image
docker pull docker.pkg.github.com/lianglouise/image_repo/image-repo:latest

#Run the Image
docker run -it --rm -p 8080:8080 -v "$(pwd)":/app/images --name my-running-app my-golang-app
```

## Run directly with `go` command

```shell
# install the packages required
go get -v github.com/gorilla/mux

git clone https://github.com/LiangLouise/image_repo.git
cd image_repo
go run .
```  
