# myzinx

## 一、简介

myzinx是根据[刘丹冰老师视频课程](https://www.bilibili.com/video/av71067087/)在学习[zinx](https://github.com/aceld/zinx)，在学习zinx框架的过程中的产物，使用`go1.22`进行开发，项目使用`go modules`的方式进行开发。

## 二、建议

1. 因为刘丹冰老师录制视频课程的时间比较早，当时还是采用的`go path`的开发模式，建议现阶段学习的小伙伴一开始就采用`go module`的开发模式。
2. 在测试模块是否可用时，可以采用单元测试的方式，而非另外开包去`main`函数中开启。
3. 可以优化一些方法，提高开发效率和开发时的容错，比如`server`端监听时采用`net.Listen()`,回显发送的内容使用`io.copy()`等。

## 参考资料

1. net包的一些常用方法的解读：[银角代王——go源码解析之TCP连接](https://www.jianshu.com/p/8e41a7aa5f07)
2. 通过impl命令快速生成给对应接口实现，提升代码开发效率：[go install github.com/josharian/impl@latest](https://github.com/josharian/impl)

## 加入我们

公众号搜索：程序员vastea
