//go:build linux && !gccgo
// +build linux,!gccgo

package nsenter

/*
#cgo CFLAGS: -Wall
extern void nsexec();
void __attribute__((constructor)) init(void) {
	nsexec();
}
*/
import "C"

/*
`import "C"`之前如果有注释，这部分注释称为前言(preamble)，编译C代码的时候前言作为作为header。
前言中可以包含任何C语言代码，包括函数和变量的声明和定义。
这些代码可以在Go代码中被使用，就像它们被定义在名称为"C"的包中一样。
所有在序言中声明的名字都可以使用，即使它们以小写字母开头。
例外：序言中的静态变量不能在Go代码中使用。
参考引用`https://pkg.go.dev/cmd/cgo`
*/

/*
CFLAGS、CPPFLAGS、CXXFLAGS、FFLAGS和LDFLAGS可以在前言注释中用#cgo指令定义，
用来调整C、C++或Fortran编译器的行为。定义在多个指令中的值会被串联起来。
该指令可以包含一个构建约束的列表，将其作用限制在满足某一约束的系统上
（关于约束语法的细节，请参见https://golang.org/pkg/go/build/#hdr-Build_Constraints）。
参考引用`https://pkg.go.dev/cmd/cgo`
*/

/*
`__attribute__((constructor))`和` __attribute__((destructor))`是gcc特定的语法。
`__attribute__((constructor))`可以在main函数执行之前(或library被加载之前)被调用。
` __attribute__((destructor))`可以在main函数执行之后(或library被卸载之后)被调用。
参考引用`https://www.geeksforgeeks.org/__attribute__constructor-__attribute__destructor-syntaxes-c/`
*/
