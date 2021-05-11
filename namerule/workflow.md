# workflow

工作流(workflow)指按顺序执行一系列工序(process)

每个工序不应该返回任何值。

工序的执行包含两个部分，交付前和交付后

交付指将工作的当前状态交给接受者(receiver)，由接受者继续处理

如不交付意味着工作流中止

工序交付前前可以直接panic

工序交付后可以执行相应代码，一般为资源释放等处理，需要自行处理所有错误



工序的标准形式为

type Receiver func(*Context)

type Process func(ctx *Context,next Receiver)


ctx *Context可以替换为其他的需要在每次执行时展开的上下文

多个工序可以组合(compose)成一个工序

```

ComposeProcess(series []Process) Process{
    return func(ctx *Context,next Receiver){{
        if len(series)==0{
            next(ctx)
            return
        }
        ComposeProcess(series[1])(ctx,next)
    }
}   

```

