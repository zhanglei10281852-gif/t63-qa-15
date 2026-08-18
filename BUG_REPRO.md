# Bug Reproduction

## 包的性质

当前 test_model_fix 保存的是被测模型修复后的结果源码，不是初始含 Bug 源码。要复现原始缺陷，必须检出下面固定的 parent SHA；不要在当前修复结果源码上期待重新出现修复前失败。生成系统使用的可信验证补丁和完整验证日志仅在本地留存，不提交到结果分支。

## 问题现象

先不要修改代码。服务启动时请求已取消，但初始化流程仍继续写入并占用启动协程。请定位取消信号在初始化调用链中的丢失位置，并提供可重复的取消后写入证据。

## 含 Bug 版本

- 仓库：zhanglei10281852-gif/t63-qa-15
- 仓库地址：https://github.com/zhanglei10281852-gif/t63-qa-15.git
- parent SHA：9ab0b23052f2ca3de49366aa3ef8b8b58fe2e186

## 复现步骤

```bash
git clone -- https://github.com/zhanglei10281852-gif/t63-qa-15.git bug-repro
cd bug-repro
git checkout --detach 9ab0b23052f2ca3de49366aa3ef8b8b58fe2e186
go test ./internal/seed -run TestSeedStopsWhenBootstrapContextIsCancelled -count=1
```

## 双架构完整错误信息

### linux/amd64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/seed -run TestSeedStopsWhenBootstrapContextIsCancelled -count=1
--- FAIL: TestSeedStopsWhenBootstrapContextIsCancelled (0.50s)
    cancellation_test.go:48: seed did not stop after cancellation
FAIL
FAIL	sanitation-operations/internal/seed	0.504s
FAIL

```

stderr：

```text
warning: internal/seed/cancellation_test.go has type 100755, expected 100644
warning: internal/seed/cancellation_test.go has type 100755, expected 100644

```

### linux/arm64

- 容器内复现预期退出码：1
- 容器内复现实际退出码：1

stdout：

```text
$ go test ./internal/seed -run TestSeedStopsWhenBootstrapContextIsCancelled -count=1
--- FAIL: TestSeedStopsWhenBootstrapContextIsCancelled (0.51s)
    cancellation_test.go:48: seed did not stop after cancellation
FAIL
FAIL	sanitation-operations/internal/seed	0.647s
FAIL

```

stderr：

```text
warning: internal/seed/cancellation_test.go has type 100755, expected 100644
warning: internal/seed/cancellation_test.go has type 100755, expected 100644

```

## 通过条件

在触发条件下，定向测试 TestSeedStopsWhenBootstrapContextIsCancelled 应通过，相关包、全量测试、竞态测试和构建检查均通过；回退 gold 唯一修复后定向测试重新失败。
