# myzinx

## 一、简介

myzinx是根据[刘丹冰老师视频课程](https://www.bilibili.com/video/av71067087/)在学习[zinx](https://github.com/aceld/zinx)
，在学习zinx框架的过程中的产物，使用`go1.22`进行开发，项目使用`go modules`的方式进行开发。

## 二、Base版本纵览

1. v0.1: 实现服务器`server`，监听一个 ip 和端口，收到连接之后直接处理数据(将发送的信息写回)。
2. v0.2: 新增`connection`层，将监听到的连接交给对应的`connection`去处理，这样也方便对本次连接进行处理，决定处理方式和何时关闭连接。
3. v0.3: 新增`request`和`router`，即定义客户端请求，以及将消息处理方式提供给框架使用者 diy，增加框架的灵活性。
4. v0.4: 新增`zconf`包，将`server`和`myzinx`框架的属性通过配置文件确定，比如 ip、端口、版本等。
5. v0.5: 新增`message`层，用于定义一次请求中的数据格式，使用该格式 (TLV) 解决 tcp 粘包问题。
6. v0.6: 新增多路由模式，即根据 v0.5 中的`messageID`，定义不同 ID 的 message 的处理方式。
7. v0.7: `connection`层包含了对本次连接的处理，使用读写分离的 goroutine 处理，方便在读/写前做处理，各司其职。
8. v0.8: 新增消息队列及任务池机制，减少业务处理的 goroutine 的最大个数，以减小过多的 goroutine 之间切换的成本。
9. v0.9: 新增链接管理模块`connectionManager`，用于管理链接的生命周期，并增加链接创建和销毁时的钩子方法。

## 三、建议

1. 因为刘丹冰老师录制视频课程的时间比较早，当时还是采用的`go path`
   的开发模式，建议现阶段学习的小伙伴一开始就采用`go module`的开发模式。
2. 在测试模块是否可用时，可以采用单元测试的方式，而非另外开包去`main`函数中开启。
3. 开发时可以参考[刘老师zinx框架](https://github.com/aceld/zinx)最新的处理，经过各位大神的共同建造，相信大家都能学到新的东西。
4. 遇到新的接口，可以跟进看一下注释，有助于提高理解，减少犯错。

## 参考资料

1. `net`包的一些常用方法的解读：[银角代王——go源码解析之TCP连接](https://www.jianshu.com/p/8e41a7aa5f07)
2. 安装`impl`命令快速生成给对应接口实现，提升代码开发效率：
   [go install github.com/josharian/impl@latest](https://github.com/josharian/impl)

## 加入我们

公众号搜索：程序员vastea
