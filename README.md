Image Repo
----

## Run with Docker

```bash
# Pull the Image
docker pull docker.pkg.github.com/lianglouise/image_repo/image-repo:latest

# Run the Image
docker run -it --rm -p 8080:8080 -v "$(pwd)":/app/images --name image-repo image-repo

# Or use docker-compose
docker-compos up
```

## Run directly with `go` command

```bash
# install the packages required
go get -v github.com/gorilla/mux github.com/mattn/go-sqlite3

git clone https://github.com/LiangLouise/image_repo.git
cd image_repo
go run .
```  
