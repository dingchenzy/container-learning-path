# helm介绍

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

在 helmv2 中使用，之后直接使用 helm client 实现。

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

首先使用 `helm show value` 查看 chart 安装的配置，然后可以创建 yaml 文件来修改配置。

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

# Chart介绍

本文以 `wordpress` chart 包为例，介绍 chart 的构成及使用。

```sh
helm pull bitnami/wordpress
```

`pull` 到本地之后，可以通过修改 `values.yaml` 文件修改 chart 内部的值。

```sh
wordpress/
	Chart.yaml		#  内含 chart 信息的 yaml 文件，包括 chart 的版本等
	values.yaml		# 配置值文件
	values.schema.json	# 使用 json 格式的 values.yaml 文件
	charts			# chart 依赖的其他 chart
	crds		# 自定义资源的定义
	templates	# 模板目录，当和 values 结合时，可以生成有效的 kubernetes manifest 文件（资源清单）
	templates/NOTES.txt	# 简要使用说明的文件
```

在一个 Chart 中，`Chart.yaml` 文件是必须的，其中包含了 chart 的版本等。

```yaml
annotations:  # 自定义元数据内容
  category: CMS
# type: chart 类型，分为 application 和 library 两种类型，默认类型是应用类型，库类型提供针对 chart 构建的实用程序和功能
#       库类型不包括任何资源对象，应用类型可以作为库类型 chart 使用可以通过设置为 library 实现
apiVersion: v2    # Chart APi 版本，helm 3 选择的应该是 v2，helm 2 选择的版本应该是 v1
# v1 和 v2 的区别如下：
# dependencies字段定义了chart的依赖，针对于v1 版本的chart被放置在分隔开的 requirements.yaml 文件中 （查看 Chart 依赖).
# type字段, 用于识别应用和库类型的chart（查看 Chart 类型).
appVersion: "5.8.0" # 应用版本，建议使用引号，也就是 wordpress 的版本，appVersion 对之后的 chart 计算并无影响，并且 chart 看的是通过 version 字段指定的版本看
dependencies: # Chart 必要条件的列表，依赖的 chart 包
- condition: mariadb.enabled  # 解析的 yaml 存放路径，以及是否要启用 （enabled/disabled）
  name: mariadb   # chart 名称
  repository: https://charts.bitnami.com/bitnami    # 仓库 URL
  version: 9.x.x  # chart 版本
- condition: memcached.enabled
  name: memcached
  repository: https://charts.bitnami.com/bitnami
  version: 5.x.x
- name: common
  repository: https://charts.bitnami.com/bitnami
  tags:   # 用于一次启用禁用一组 chart 的 tag
  - bitnami-common
  version: 1.x.x
# kubeVersion: 该字段指定的语义化版本，可以使用这种方式 >= 1.13.0 < 1.15.0，检测不成功将会出现问题
# 也可以使用操作符 || 连接 >= 1.13.0 < 1.14.0 || >= 1.14.1 < 1.15.0
description: Web publishing platform for building blogs and websites.   # 添加的注释
home: https://github.com/bitnami/charts/tree/master/bitnami/wordpress # 项目 home 页面的 url
icon: https://bitnami.com/assets/stacks/wordpress/img/wordpress-stack-220x234.png # 用作 icon 的图片
# deprecated: 用来标识已经废弃的 chart
keywords: # 项目的一组关键字
- application
- blog
- cms
- http
- php
- web
- wordpress
maintainers:  # 声明
- email: containers@bitnami.com # 维护者邮箱
  name: Bitnami   # 维护者名字
  # url: 维护者 url
name: wordpress
sources:  # 项目源码的 url
- https://github.com/bitnami/bitnami-docker-wordpress
- https://wordpress.org/
version: 12.1.7 # 语义化版本
# 具体内容 https://semver.org/lang/zh-CN/
```

## chart dependency

helm 中可能某些 chart 会依赖到其他的某个 chart，这些可以使用 chart.yaml 文件中的 dependencies 字段动态链接，并进行手动配置

```yaml
# helm repo add bitnami https://charts.bitnami.com/bitnami
dependencies:
- name: mysql
  version: 8.8.6		# [chart 版本]
  # repository: "https://charts.bitnami.com/bitnami"
  repository: "@bitnami"	# 可以使用 repo 名字代替 repo url
```

确定好了依赖之后，就可以使用 `helm dependency update` 命令将你的依赖文件下载到 `charts/` 目录。

## chart alias

`alias` 字段可以实现一个 chart 在依赖多个相同的 chart 时使用，需要手动移动到 charts 目录下并修改目录名称。

## helm condition tags

condition：设置该 chart 是否启用

tags：与 chart 关联的 yaml 格式标签列表，在顶层 vlaue 可以通过指定 tag 和 bool 值，启用和禁用所有带 tag 的 chart

```yaml
# parentchart/Chart.yaml

dependencies:
  - name: subchart1
    repository: http://localhost:10191
    version: 0.1.0
    condition: subchart1.enabled, global.subchart1.enabled
    tags:
      - front-end
      - subchart1
  - name: subchart2
    repository: http://localhost:10191
    version: 0.1.0
    condition: subchart2.enabled,global.subchart2.enabled
    tags:
      - back-end
      - subchart2
```

```yaml
# parentchart/values.yaml

subchart1:
  enabled: true
tags:
  front-end: false
  back-end: true
```

在上面的例子中，含有 front-end 标签的 chart 都会被禁用，但是含有该标签的 chart subchart1 的 condition 中包含 `subchart1.enabled=true` 所以会覆盖 `front-end ` 标签。

执行命令 `

```sh
# subchart2 并不会生效，尽管 back-end 设置为 true，但是因为上级字段也就是 sunchart2.enabled 设置为 false，所以该标签并不会生效
helm install --set tags.front-end=true --set subchart2.enabled=false
```

## import-values字段

在 charts 文件中声明子 chart 中的 value 内容，可以从子 chart 的 values.yaml 文件中导入键名

```yaml
# wordpress/Chart.yaml
dependencies:
- name: subchart
  repository: http://localhost:9000
  version: 0.1.0
  import-values:
  - data
```

```yaml
# wordpress/charts/mysql/values.yaml
exports:
  data:
    myint: 90
```

```yaml
# wordpress/values.yaml
myint: 99
```

另一种方式，使用子 chart 中未声明 `exports` 的值，需要指定导入的源键名和目标路径。

```yaml
# wordpress/Chart.yaml
dependencies:
- name: subchart1
  repository: http://localhost:9000
  version: 0.1.0
  import-values:
  - child: default.data
    parent: myimports
```

这个例子中会直接在 `sunchart1` chart 中的 `values.yaml` 文件中的 `default.data` 目录下查找

```yaml
# wordpress/subchart1/values.yaml

default:
  data:
    myint: 999
    mybool: true
```

```yaml
# wordpress/values.yaml

myimports:
  myint: 999
  mybool: true
  mystring: "helm rocks!"
```

## values.yaml

该文件直接通过 go template 进行渲染得到，helm 并不提供任何配置项。

在 template 目录下配置的 template，会通过 helm 解析后通过 values.yaml 文件中的内容进行渲染。

### 内置 values

- `Release.Name`: 版本名称(非chart的)
- `Release.Namespace`: 发布的chart版本的命名空间
- `Release.Service`: 组织版本的服务
- `Release.IsUpgrade`: 如果当前操作是升级或回滚，设置为true
- `Release.IsInstall`: 如果当前操作是安装，设置为true
- `Chart`: `Chart.yaml`的内容。因此，chart的版本可以从 `Chart.Version` 获得， 并且维护者在`Chart.Maintainers`里。
- `Files`: chart中的包含了非特殊文件的类图对象。这将不允许您访问模板， 但是可以访问现有的其他文件（除非被`.helmignore`排除在外）。 使用`{{ index .Files "file.name" }}`可以访问文件或者使用`{{.Files.Get name }}`功能。 您也可以使用`{{ .Files.GetBytes }}`作为`[]byte`访问文件内容。
- `Capabilities`: 包含了Kubernetes版本信息的类图对象。(`{{ .Capabilities.KubeVersion }}`) 和支持的Kubernetes API 版本(`{{ .Capabilities.APIVersions.Has "batch/v1" }}`)

### 注意

chart 包含的默认 values 文件，必须被命名为 `values.yaml` 但是可以在命令行指定的文件可以是其他名称。

如果 `helm install` 或者是 `helm upgrade` 使用 `--set` 参数，这些值会在客户端被简单的转换为 yaml

如果 values 文件存在任何必须的条目，会在 chart 模板中使用 required 函数声明为必须的

## crd 资源的使用 

如果 chart 中有依赖 crd 资源的资源，那么可以将 crd 资源的定义文件放置在 `crds/` 目录下，这样当 helm 在安装 chart 时会首先执行 `crds/` 目录下的 crd 定义清单，防止 chart 内有使用该资源的 kubernetes 资源。

## starter包

`helm create` 命令可以附带 `--starter` 选项，可以指定一个 `starter chart`。

Starter就只是普通chart，但是被放置在`$XDG_DATA_HOME/helm/starters`。作为一个chart开发者， 您可以编写被特别设计用来作为启动的chart。设计此类chart应注意以下考虑因素：

- `Chart.yaml`会被生成器覆盖。
- 用户将希望修改此类chart的内容，所以文档应该说明用户如果做到这一点。
- 所有出现的`<CHARTNAME>`都会被替换为指定为chart名称，以便chart可以作为模板使用。

当前增加一个chart的唯一方式就是拷贝chart到`$XDG_DATA_HOME/helm/starters`。在您的chart文档中，您可能需要解释这个过程。