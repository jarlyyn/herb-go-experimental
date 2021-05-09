# 模块,系统，组件

一个服务可以有多个模块

模块有顺序

系统可以Start以及Stop

系统有StopHook，函数列表

系统t推荐有HandleRuntimeError 函数，用于接受运行时错误

模块在Start时可以依次挂载，挂载时可添加StopHook进行资源释放

模块需要有模块ID

在stop时或者Start失败时需要依次逆序调用StopHook

原则上Stophook需要强制关闭，不应该返回任何信息。如需要记录错误需要通过HandleRuntimeError方式调用。如有需要fatal的情况，可以直接exit。

组件是系统运行时 模块处理数据的地方，string/interface{}对，只有在start前可以注册，start后可以获取，需自行处理并发问题

组件的释放应该通过模块挂载时添加的StopHook处理