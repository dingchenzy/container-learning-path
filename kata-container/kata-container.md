# 简介

## 安全容器的诞生

安全容器是一种运行时技术，并且为容器应用提供一个完整的操作系统执行环境。

> 与普通容器相比，最重要的区别是每个容器，都运行在一个单独的微型虚拟机中，并且拥有独立的操作系统内核，以及虚拟化层的安全隔离。因为云容器实例采用的是共享多租集群，因此容器的安全隔离比用户独立拥有私有 kubernetes 集群有更严格的要求。通过安全容器，不同租户之间的容器之间，内核、计算资源、存储和网络都是隔离开的。保护了用户的资源和数据不被其他用户抢占和窃取。

但是应用的执行与宿主机操作系统隔离开，避免应用直接访问主机资源，从而可以在容器主机之间或容器之间提供额外的保护。

## 间接层

安全容器的基础，是 Torvalds 提出的 

“安全问题的唯一正解在于允许那些 Bug 发生，但通过额外的隔离层来挡住。”

因为 linux 系统本身这样的规模已经非常庞大，无法通过语句分析或者是理论上证明程序是没有 Bug 的。所以需要增加额外的隔离层，来减少漏洞或者是这些漏洞造成的被彻底攻破的风险。

## kata-container

安全容器，项目的前身是 runV 以及 intel 的 clear Container 项目。

制作安全容器，只需要一个隔离层，虚拟机本身（是一个现有的隔离层，比如像阿里云和 AWS）只要虚拟机有个内核就可以满足 OCI 的定义，也就是提供了 Linux ABI 的运行环境，在这个环境中跑应用程序并不难。

> Linux ABI 是应用程序的二进制接口，这个接口内提供了可以提供给应用程序调用的各种库文件，提供调用，Windows 与 Linux 是不同的库，所以说应用程序分为 exe 以及执行文件。

**但是唯一的缺陷是虚拟机不够快，阻碍了容器在环境中的应用，如果够快的话，那么就可以使用虚拟机做隔离的安全容器技术，这也是 kata-container 的一个思路。就是用虚拟机做 kubernetes 的 PodSandbox（Pod沙盒）**

![img](https://static001.infoq.cn/resource/image/48/d5/488cea8b7dacf38c28d7e38ca24f9ed5.png)

通过 kubelet 的 CRI 找到 containerd，containerd 找到 containerd-shim 执行容器命令，容器 cmd/spec 经过 kata-runtime 或者是 kata-shim 发送请求到 kata-proxy 实现对 Pod Sandbox 中的容器执行命令，执行命令时使用的是 Sandbox 中单独的内核空间，并不是使用的宿主机的内核。

![kata2](kata-container/kata2.png)

现在云中的容器使用的是通过 vm 虚拟化之后的虚拟机 kernel，而 kata-container 实现的是直接在容器中创建一个 kernel，让容器中的应用程序之后使用都是会调用自身的 kernel。

![kata3](kata-container/kata3-1617004559949.png)

- Runtime：符合 OCI 规范的容器运行时工具。主要用来创建轻量级虚拟机并通过 Agent 控制虚拟机内容器的生命周期。
- Agent：运行在虚拟机中的运行时代理组件，主要用来执行 Runtime 传给它的指令并在虚拟机内管理容器的生命周期。
- Shim：Shim 相当于 Containerd-Shim 的适配，用来处理容器进程的 stdio 和 signals。Shim 可以将 Containerd-Shim 发来的数据流传给 Proxy，Proxy 再将数据流传输给微型虚拟机中的 Agent，Agent 传输给容器并执行相应的动作，同时Shim 也支持将内部 Agent 产生的信号传输给 Proxy，Proxy 再传输给 Shim。
- Proxy：为 Runtime 和 Shim 分配访问 Agent 的控制隧道，路由 Shim 和 Agent 发出的数据流。
- Kernel：是一个单独的内核，提供轻量虚拟机的内核，最小的 4M，根据不同的需要提供几个内核。

> 一个 Pod 的多个容器是被放到同一个微型虚拟机中，也可以根据某些需求，实现共享某些 namespace

## Pod Sandbox

![img](https://static001.infoq.cn/resource/image/a4/dd/a45309ee7528e02ee8e3f44140b595dd.png)

DAX：共享内存，让多个微型虚拟机共享内存中只读的内存空间，节省内存空间。

Virtio-fs：存储共享，共享宿主机中 docker 的 rootfs 目录。

Guest Kernel：仅提供容器的运行时环境，并不提供其他虚拟机给用户的感受，例如登录系统等，纯粹的只是一个内核。

使用 Pod Sandbox 这种方式，使用起来与容器无异，但是实际上在底层看来是一个虚拟机。

## gVisor

![img](https://static001.infoq.cn/resource/image/74/79/74ea3425d3e3136995061f27b38cb079.png)

进程级别的虚拟化，google 公司为 kata-container 的贡献就是自己内部开发五年并一直在应用的 gVisor 安全容器的解决方案。

不同于 kata 的运行逻辑，在 gVisor 中，是采用 Go 语言重新对内核中的用户空间进行了二次开发，生成 sentry 内核。

sentry 内核，该内核的特点就是可以将攻击者的攻击面减少，也就是会将内核空间中的一些 syscall 调用进行筛选和删减在 sentry 内核中的使用，筛选出一些常用并且较为安全的 syscall 继续提供给用户使用，而一些容易受到攻击且用户并不常用的系统调用将会被删减。

syscall open()，原始 linux kernel 中的 open 系统调用极易受到系统的攻击，因为 linux 中的功能都是依赖于文件实现的，所以通过 open() 调用就会有很多操作。而 sentry 内核将 open() 调用封装到了 Gofer 中，当在容器中调用 open() 调用时会交给 Gofer 进程来执行，执行过程不会允许不良操作。

## 安全容器未来

安全容器创造出的微型虚拟机，在未来的云原生基础架构中会广泛应用，因为其提供的隔离性以及运行的速度相对于传统虚拟化是有过之而无不及的。

在内部集成的 DAX virtio-fs 这种虚拟化组件是非常的节省资源并适用于云原生环境，并且其内部的隔离机制，可以将用户的敏感信息封装到容器的应用中，这样即使运维人员在维护时也不会触及内部的敏感数据，使用的仅仅是外部的 Pod Sandbox 中的内核。

gViosr 虽然目前没有任何优势，但是这种通过在用户空间运行一个 linux 内核并运行应用程序的思路，在之后可能会广泛应用。kata 目前的技术已经成熟并且相对比较好理解。

