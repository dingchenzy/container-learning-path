## helm简介

是第三方的软件包管理器，但是应用较为广泛。

可以简单的部署复杂的应用中间件，如果手动在 kubernetes 实现部署，那么需要创建指定的 Pod 指定的 service 以及 configmap 等资源，包括存储卷等。但是这些复杂的步骤逐步就会被 helm 取代。

所有的配置以及操作都通过 chart 的方式进行声明，当使用 helm 进行部署时会自动的调用。

## helm概念

chart：是一个 helm 包，其中包含了应用需要的工具以及 kubernetes 集群中的服务定义，类似于 rpm 文件。

release：在 kubernetes 运行的一个 chart 实例，一个 chart 可以被安装多次。但是如果想要在服务器上运行两个应用程序实例，也就是会基于同一个 chart 创建两个实例，每次创建都会在 kubernetes 上生成 release，并且会拥有不同的 release 名称。

repository：存放和共享 chart 仓库，其中 helm 官方提供了 chart 仓库。

helm 的主要任务就是从 repository 上查找指定的 chart 然后以 release 的方式安装到集群之上。

## helm组成

### HelmClient

在命令行的管理工具，可以通过命令来管理 release chart repository 对象的管理能力，可以直接通过二进制文件或者是脚本进行安装。

### TillerServer

是客户端工具与 kubernetes 集群交互的中间人，基于 chart 定义，生成和管理各种的 kubernetes 资源对象。

可以使用 helm init 命令创建出 tillerServer，并且会通过当前的 context 指定的集群名称在 kube-system 名称空间下创建 deployment 和 service。但是在初始化 tillerServer 时创建出的镜像源时 gcr.io 的，所以可以通过 helm init 命令中的 tiller-image 来指定下载 tillerServer 的镜像。

安装完成后可能会出现错误，socat not found，安装 socat 工具实现 tillerServer 端口转发即可。

## 安装

### helm初始化安装

```bash
可以使用 helm init --tiller-image [镜像] 设置指定的 TillerServer 的镜像。
```

### 本地安装

```bash
--host				选项指定 tillerserver 监听地址
HELM_HOST			环境变量指定的监听地址
```

也可以在服务器本地运行 tillerServer，但是需要指定 --host 监听的地址或者是设置 HELM_HOST 环境变量。但是即使是在本地安装 kubernetes，tillerServer 依然会使用 kubectl 配置中的 context 配置去连接 kubernetes 集群。

## 用法

```bash
helm search			#查找指定的 chart

helm inspect		#查看 chart 的详细信息

helm install		#安装 chart

helm status			#该命令查看安装完成后的 release 运行状态
```

## 自定义chart

首先使用 `helm inspect` 查看 chart 安装的配置，然后可以创建 yaml 文件来修改配置。

```bash
helm install -f/--values [file.yaml]		#修改 chart 配置文件
helm install --set key=value				#直接通过 set 的方式赋值
```

## 应用的更新回滚

```bash
helm upgrade			#该命令对 chart 安装的 release 进行升级操作
helm rollback			#回滚操作，会对 chart 已经升级的操作执行回滚操作，实现对 release 的版本回滚操作
helm history			#查看该 release 的升级操作的版本号，使用 rollback 命令可以进行回退操作
```

例如更新一个资源的部分内容，就可以使用 upgrade 操作进行升级配置。

## 删除release

```bash
helm delete				#可以删除一个 release，使用 helm delete release 从 k8s 集群中删除
```

## repo管理

```bash
helm repo list				#查看已有的仓库
helm repo add				#添加仓库
helm repo update			#更新仓库
helm fetch [chart]			#下载 chart
```

