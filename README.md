<p align="center">
	<img src="https://user-images.githubusercontent.com/43493852/147871833-8e8f0f21-d6c6-4be9-9e0f-892572bae942.png"></img>
	<p align="center">
		<img src="https://github.com/GoogTech/fileGo/workflows/Go/badge.svg"></img>
		<img src="https://goreportcard.com/badge/github.com/GoogTech/fileGo"></img>
		<img src="https://img.shields.io/github/commit-activity/m/google-golang/fileGo?color=ff69b4"></img>
		<img src="https://img.shields.io/github/repo-size/google-golang/fileGo"></img>
		<img src="https://img.shields.io/github/license/google-golang/fileGo.svg"></img>
	</p>
</p>


# fileGo
> 📁 A bijou file management tool and nothing framework be used in it !


## How To Run
### Window
```shell script
> go run main.go   

# or
> go build
> fileGo.exe
```
### Linux / Mac
```shell script
> go run main.go

# or
> go build
> ./fileGo
```

Then the browser will be opened auto and jump to the link : *https://localhost:8080* if run it successfully .


## Permission( `by session` )
* user : *http://localhost:8080/*
* admin : *http://localhost:8080/?admin*
> For more information please refer to the `fileGo.conf` file bro


## Upload file on linux
```shell script
# it's will rename the file, such as changing filename.txt to filename.1.txt 
curl -F "file=@a.jpg;filename=a.jpg"  http:/ip:port/upload

# it's will rename and mandatorily cover the old file..
curl -F "file=@a.jpg;filename=a.jpg"  http:/ip:port/upload/f
```
