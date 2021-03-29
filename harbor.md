## harbor简介

docker 的私有 registry，当企业内的镜像需求过多，如果一直从阿里云之上上传获取镜像，那么就会出现效率过低的情况。可以在本地搭建一个私有 registry，实现在本地上传和下载镜像并使用提高效率。

harbor 作为云计算基金会毕业的项目，优势不言而喻，其支持分布式的特性也是提高了高可用性。

### Harbor 目录结构

```bash
/data/
├── ca_download			#存放 ca 证书
├── database		#存放镜像的数据以及配置信息
├── job_logs		#存放工作日志
├── psc		#
├── redis		#存放缓存内容，如果 Harbor 被停止，那么 redis 容器会自动的保存 rdb 文件
├── registry		#存储整个镜像库的结构
└── secret		#存储私有证书文件以及验证信息
```

## harbor安装

是通过 docker-compose 构建的方式进行安装，docker-compose 仅支持单机构建，所以这种方式仅仅支持在一个主机上安装。

1. 下载 harbor 的安装包

```http
https://github.com/goharbor/harbor/releases
可以选择各种版本下载
```

2. prepare install

```bash
yum install docker docker-compose -y				#构建 harbor 的必须软件
tar -xf /usr/local/src/harbor.tar.gz -C /usr/local/		#解压 harbor
cd /usr/local/harbor				#进入安装目录
```

3. 配置 harbor.yml

```http
https://goharbor.io/docs/2.2.0/install-config/configure-yml-file/
```

4. 安装 harbor

```bash
./install.sh
```

5. 如果想要停止 harbor

```bash
docker-compose down				#停止
docker-compose up				#启动
```

## 升级Harbor

**1.8 版本升级到 1.9 版本。**

需要迁移镜像将旧 harbor 中的镜像以及配置导入到新 harbor 中，其内部工作原理就是通过这种 api 管理的方式获取配置，然后再通过遍历的方式将配置写入到新 harbor 中。

```bash
curl -u admin:Harbor12345 -k -X GET "http://192.168.1.31/api/projects" | jq -r .[].name
```

1. 备份旧 harbor 的数据库

```bash
cd /usr/local/harbor
docker-compose down
cd ..
mv harbor harbor_backup1 
cp -rf /data/database/ . 
```

2. 升级

该命令的作用就是基于原有 harbor 的配置文件生成一个新版本的配置文件（可以自行查看配置文件对比）

```bash
docker run -it --rm -v /usr/local/harbor_old/harbor.yml:/harbor-migration/harbor-cfg/harbor.yml -v /usr/local/harbor/harbor.yml:/harbor-migration/harbor-cfg-out/harbor.yml goharbor/harbor-migrator:v1.9.0 --cfg up




docker run -it --rm -v [old_yaml]:/harbor-migration/harbor-cfg/harbor.yml -v [new_yaml]:/harbor-migration/harbor-cfg-out/harbor.yml goharbor/harbor-migrator:v1.9.0 --cfg up
```

3. 安装新版本

```bash
./install
```

访问查看，内容与旧版本一致。

## 大版本升级Harbor

**harbor v1.10.4 升级到 harbor v2.0.0**

因为大版本差异过大，所以需要进行数据迁移，数据迁移实际上应用的是 harbor 安装包中的 prepare 脚本。其内容就是将目录移动到新版本中。

```bash
docker run -it --rm -v /:/hostfs goharbor/prepare:v2.1.1 migrate -i /usr/local/harbor_10/harbor.yml -o /usr/local/harbor_20/harbor.yml



docker run -it --rm -v /:/hostfs goharbor/prepare:v2.1.1 migrate -i [old_yml] -o [new_yml]
```

安装

````bash
 ./install.sh  --with-chartmuseum --with-clair
````

## 主从同步

pull 方式，slave 主动向 master harbor pull 镜像，减轻了 master harbor 的负载。

push 方式，master 主动将镜像 push 到远端 slave 服务器，增大了 master harbor 的负载。

触发方式分为手动和定时触发。

可以基于 namespace 或者是标签指定同步的镜像，例如从服务器只同步 latest，那么就可以指定标签选择器为 latest，并且选择覆盖方式，就可以保证 slave harbor 中的镜像始终都是最新的。

如果镜像同步到 slave harbor 需要保存，那么就不建议选择覆盖方式。

### 配置步骤

1. 配置规则
2. 配置同步策略

