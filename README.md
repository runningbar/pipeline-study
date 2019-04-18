# `pipeline`学习代码

> 每个文件都是一个独立单元，互不关联，对任意单个文件`go run`检查效果

* `basic-squaring-numbers.go`：一个简单的3阶段`pipeline`，一条线到底，没有分叉路径。

* `fan-in-and-fan-out.go`：示例扇入和扇出。

* `explicit-cancellation.go`：使用`channel`向上游明确退出信号。

* `son-or-sibling.go`：**证明了协程之间没有父子关系**，都是并列的，但是`main`进程退出时，所有的协程都会强制退出。

* `md5-pipeline/`是一个实际的`pipeline`用例，包含`fan-in/out`，`cancel`，固定数量的`goroutine`。
    * `md5-pipeline/serial.go`：最简单，无pipeline，串行处理
    * `md5-pipeline/parallel.go`：每个文件启一个协程计算MD5，fan-in，主进程遍历文件树
    * `md5-pipeline/parallel.go`：文件树平铺处理（stage1），对平铺的文件路径多协程消费（fan-out）并写入同一`channel`（fan-in），限制消费者数量（bound）