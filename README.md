# kube-control-api
<b> 描述：</b> 对kubernetes 集群一些通用的接口都可封装到这里，后续可以搭建图形化界面提供使用，打算分两部分，
一部分是运维需要的接口，例如搭建快速集群，还有一部分封装kubernetes的监控接口，如 etcd 的备份等。


## 接口：
### 初始化集群：
endponint: POST /api/downloadResource <br>
描述： 预先下载事先准备好的docker image, image 约1G，包含kubernetes 集群所需的image, 具体包版本可以查看 config/resource/spray-v2.21.0c_k8s-v1.26.4_v4.4-amd64/package.yaml <br>
技术逻辑： POST 接口接受到image 信息类容，如果没有则读取默认配置config/resource/spray-v2.21.0c_k8s-v1.26.4_v4.4-amd64/package.yaml，然后运行pull-resource-package.sh 下载。<br>
待改进： 这个接口使用了对文件加锁的形式来防止同时运行多个task。我觉得这样不太优雅，可以参考 kubernetes 源码使用queen + 生产者，消费者的模式提高效率，这种写法像Java的写法，失去Golang的灵活性

