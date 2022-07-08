## nsenter

The `nsenter` package registers a special init constructor that is called before 
the Go runtime has a chance to boot.  This provides us the ability to `setns` on 
existing namespaces and avoid the issues that the Go runtime has with multiple 
threads.  This constructor will be called if this package is registered, 
imported, in your go application.

> `nsenter`包注册了一个特殊的init函数，在Go运行时启动之前被调用。
> 这为我们提供了在现有命名空间上`setns`的能力，避免了Go运行时在多线程方面的问题。
> init构造函数这个包在被注册、导入时被调用。

The `nsenter` package will `import "C"` and it uses [cgo](https://golang.org/cmd/cgo/)
package. In cgo, if the import of "C" is immediately preceded by a comment, that comment, 
called the preamble, is used as a header when compiling the C parts of the package.
So every time we  import package `nsenter`, the C code function `nsexec()` would be 
called. And package `nsenter` is only imported in `init.go`, so every time the runc
`init` command is invoked, that C code is run.

> `nsenter`包会`inport "C"`并使用cgo包。
> 在cgo中，`import "C"`前面紧挨着的一段注释称为前言，在编译包中C语言部分时，这个前言作为header。
> 所以每当`nsenter`包被导入时，C代码中的`nsexec()`就会被调用。
> 而`nsenter`包只在`init.go`中被导入，所以每次调用runc的`init`命令时，都会运行那段C代码。

Because `nsexec()` must be run before the Go runtime in order to use the
Linux kernel namespace, you must `import` this library into a package if
you plan to use `libcontainer` directly. Otherwise Go will not execute
the `nsexec()` constructor, which means that the re-exec will not cause
the namespaces to be joined. You can import it like this:

> 由于`nsexec()`必须运行在在Go运行时之前，才能使用 Linux 内核命名空间。
> 因此如果需要直接使用`libcontainer`，就必须导入此库。
> 否则Go将不会执行`nsexec()`构造函数，这意味着再次执行将不会join到命名空间。
> 你可以这样导入它：

```go
import _ "github.com/opencontainers/runc/libcontainer/nsenter"
```

`nsexec()` will first get the file descriptor number for the init pipe
from the environment variable `_LIBCONTAINER_INITPIPE` (which was opened
by the parent and kept open across the fork-exec of the `nsexec()` init
process). The init pipe is used to read bootstrap data (namespace paths,
clone flags, uid and gid mappings, and the console path) from the parent
process. `nsexec()` will then call `setns(2)` to join the namespaces
provided in the bootstrap data (if available), `clone(2)` a child process
with the provided clone flags, update the user and group ID mappings, do
some further miscellaneous setup steps, and then send the PID of the
child process to the parent of the `nsexec()` "caller". Finally,
the parent `nsexec()` will exit and the child `nsexec()` process will
return to allow the Go runtime take over.

> `nsexec()`首先从环境变量`_LIBCONTAINER_INITPIPE`中获取初始管道的文件描述符。
> （该变量由父进程打开，并在`nsexec()`初始进程的fork-exec中保持打开）
> init管道用于从父进程中读取bootstrap数据（namespace paths、clone flags、uid和gid映射、console path）。
> `nsexec()`调用`setns(2)`来加入bootstrap数据中提供的命名空间（如果有），
> 调用`clone(2)`用克隆flag克隆一个子进程，更新uid和gid映射，以及其他设置步骤，
> 然后将子进程的PID发送给`nsexec()`的"调用者"。
> 最后，父`nexec()`退出，子`nexec()`进程返回，Go运行时可以接管。

NOTE: We do both `setns(2)` and `clone(2)` even if we don't have any
`CLONE_NEW*` clone flags because we must fork a new process in order to
enter the PID namespace.

> 注意：即使没有任何`CLONE_NEW*`克隆标志，也会调用setns(2)和clone(2)，
> 因为我们必须fork一个新进程才能进入PID命名空间。


