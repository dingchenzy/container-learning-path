## rke

rancher 公司开发的一款快速安装 kubernetes 集群的工具，平时安装 kubernetes 需要配置电脑支持网络以及 kubectl 等插件，但是使用 rke 时，仅仅需要主机安装上 rke 支持的 docker 版本即可，简化了 kubernetes 的安装。

## 安装要求

1. ssh 用户，需要在其他主机创建用户，然后将其添加到 docker 组中，添加完成后会自动获得 root 权限。

```bash
usermod -aG docker sshuser
```

2. 禁用所有 worker 节点上的交换功能 swap

```bash
1. 可以修改 fstab 文件禁用
2. 可以 swapoff 
```



