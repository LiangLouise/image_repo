# Image Repo

## Run with Docker

```shell
# Pull the Image
docker pull docker.pkg.github.com/lianglouise/image_repo/image-repo:latest

# Run the Image
docker run -it --rm -p 8080:8080 -v "$(pwd)":/app/images --name image-repo image-repo

# Or use docker-compose
docker-compos up
```

## Run directly with `go` command

```shell
# install the packages required
go get -v github.com/gorilla/mux gorm.io/gorm gorm.io/driver/sqlite

git clone https://github.com/LiangLouise/image_repo.git
cd image_repo
go run .
```  

## Test

* SignUp
```shell
curl -X POST -H "Content-Type: application/json" -d '{"name": "testUser"}' http://0.0.0.0:8080/signup
```

* Image-Upload
```shell
curl -X POST -F 'imageName=A_great_image' -F 'file=@<filename>' -F 'isPrivate=false' -F 'userId=1' http://0.0.0.0:8080/image
```

* Image-View

```shell
curl -X GET http://0.0.0.0:8080/image/1?userid=1
```

* Image-Delete

```shell
curl -X DELETE http://0.0.0.0:8080/image/1?userid=1
```

* Search

```shell
curl -X GET http://0.0.0.0:8080/search?userid=1&text=<image_title_keyword>&page=1
```