# http_resource_get
一个命令行的下载http资源的小工具
支持下载单个url资源、支持下载文件中的多行纯url资源
按url中的path_info自动创建目录

#交叉编译命令(64bit)

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build get.go

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build get.go

CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build get.go

CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build get.go