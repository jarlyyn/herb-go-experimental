## 支持的待解析解析类型

* bool
* int
* int64
* float32
* float64
* string
* struct
* slice
* map
* 通过tag指定的 LazyLoadFunc 懒读取函数
* 通过tag指定的 LazyLoader 懒读取接口
* interface{} 

interface{} 解析出的数据取决于传入的Part的Iter字段的Step类型
* 如果为nil,即不可遍历，解析为part的Value返回值
* 如果为TypeArray，即数组对象，解析为[]interface{}
* 如果为TypeString，解析为map[string]interface{}
* 如果为TypeEmptyInterface ,解析为map[interface{}]interface{}


## 匿名字段规则

匿名字段的规则如下:

* 如果字段的对象不是个结构体的话，当作非匿名字段处理，不继续向下判断。
* 如果通过tag设置过name,强制当作子结构按设置的name处理，不继续向下判断。
* 如果配置中的匿名标签设置不为空，同时通过tag指定为匿名字段，强制将字段的所有字段当成父结构的字段处理，不继续向下判断。
* 如果通过字段名可以找到对应的子元素，当成子结构处理，不继续向下判断。
* 其他情况将所有字段当成父结构的字段处理，不继续向下判断。

 如果匿名字段还包含匿名字段的话，递归处理