# 笔记

[toc]

此文件记录`nsenter`包下源码学习笔记。

这个包下的go文件的作用就是`import "C"`并将C的前言(preamble)放到go代码之前执行，
而C代码的入口只是调用了`nsexec()`方法。有些困扰的地方是：
C代码一上来就会找`_LIBCONTAINER_INITPIPE`这个环境变量
（以及用于log的`_LIBCONTAINER_LOGPIPE`和`_LIBCONTAINER_LOGLEVEL`环境变量），
如果找不到就退出。这个原因是这样的，并不是每一次执行runc都需要执行`nsexec()`，
比如只在命令行执行`runc --help`是完全没必要走一遍`nsexec()`的。
所以只有在需要的情况下（可能只有`runc init`），才会在go代码执行前走一遍特定的初始化流程。

todo [runc源码分析](https://toutiao.io/posts/9t5ta44/preview)

## nsexec.c

### 主体逻辑

建立用于log的pipe
获取跟parent通信的pipe
确保克隆文件（避免CVE-2019-5736漏洞，这部分逻辑没看懂）
通知parent已完成initial（通过通信pipe）
从通信pipe中读取netlink配置
把config中的`oom_score_adj`数据写入`/proc/self/oom_score_adj`
*未完成*

### pipe

`nsexec.c`中的`pipe`就是通过文件来相互传递数据。
从`setup_logpipe()`和`write_log()`两个函数的定义可以看出来，
`setup_logpipe()`实际上就是从环境变量中读取了log级别和log文件描述符，
（linux的文件描述符就是个非负整数）
然后在`write_log()`中调用`write()`函数向文件描述符写入字符串。

linux中特殊的文件描述符：
`0`表示标准输入`stdin`，
`1`表示标准输出`stdout`，
`2`表示标准错误输出`stderr`。

在命令行使用`0 1 2`文件描述符：
`ls 2>&1`表示将`stderr`输出到`stdout`。
`ls 2>/dev/null`表示将`stderr`的输出“扔掉”。

### /proc/self

proc文件系统是一个伪文件系统，它只存在内存当中，而不占用外存空间。它以文件系统的方式为访问系统内核数据的操作提供接口。
通过`/proc/$pid/`来获取指定进程的信息。
`/proc/self/`等价于`/proc/本进程pid/`。
进程可以通过访问/proc/self/目录来获取自己的系统信息，而不用每次都获取pid。
`cmdline` 进程的完整命令。
`cwd` 进程启动时所在路径。
`exe` 指向进程对应的二进制文件。
`environ` 进程运行时的所有环境变量。
`fd` 进程打开的每个文件的文件描述符。

## 参考引用

- [linux特殊文件描述符及在命令行中使用](https://blog.csdn.net/weiqifa0/article/details/106270760)
- [/proc/self](https://blog.csdn.net/Zero_Adam/article/details/114853022)
