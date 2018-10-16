# TasQ 
TasQ is a background task worker with high performance.

## Cases for using:
You can use **TasQ** if your application requires some background operations. If a task requires some logic which not needs for client You can move this logic to background.

For example, there are HTTP server which requires doing some logic (send mail, exec DB queries etc.) after handling of HTTP request. In this case you may use **TasQ**.

## Benchmarking

**pkg: github.com/begmaroman/tasq**

|Test name|Iteration count|Time|
|---|---|---|
|BenchmarkTasQ_Enqueue_SuccessDone-4|3000000|565 ns/op|
|BenchmarkTasQ_Enqueue_CustomErrorDone-4|2000000|609 ns/op|
|BenchmarkTasQ_Enqueue_RetryError3-4|3000000|569 ns/op|