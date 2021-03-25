# 链式调用的命名

* SetXXX 不使用，不应该被链式调用，容易混淆
* WithXXX 链式调用，过于常用，尽量不使用。
* VaryXXX 复制并修改，不影响原值
* MergeXXX 修改，影响原有值
* ConcatXXX 链式调用，复制并追加
* AppendXXX 链式调用，追加 