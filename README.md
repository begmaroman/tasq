# TasQ 
TasQ is a background task worker with high performance.

## How it works
0. Task (go type) must implement *tasq.Task* interface which has one method *Do*. Signature of *Do* func: `func Do() error`
1. Configure **TasQ** for your application using the following settings:

    - `WorkersPollSize` - count of workers (go-routines) which will be run after start **TasQ**
    - `TaskMaxRetry` - number of function calls *Do* when that function returns `tasq.ErrRetryTask`. That's mean If function `Do` returns `tasq.ErrRetryTask` background worker try to call this function again until the number of trying less than `tasq.TaskMaxRetry` option or `Do` function returns something else.
    
2. Create new instance of **TasQ**: `tq := tasq.New()`;
3. Set up `TaskDone` and `TaskFailed` functions if need:
    
    - `TaskDone` - the function which will be called when the task is done. Parameters: `id int64` - task ID;
    - `TaskFailed` - the function which will be called when the task is failed. Parameters: `id int64` - task ID, `err error` - error;

4. Run **TasQ**: `tq.Start()`;
5. Use `tq.Enqueue(task)` for adding a task to background queue;
6. Use `tq.Close()` for stopping workers and for closing channel;

## Cases for using:
You can use **TasQ** if your application requires some background operations. If a task requires some logic which not needs for client You can move this logic to background.

For example, there are HTTP server which requires doing some logic (send mail, exec DB queries etc.) after handling of HTTP request. In this case you may use **TasQ**.

## Benchmarks:

**pkg: github.com/begmaroman/tasq**

|Test name|Iteration count|Time|
|---|---|---|
|BenchmarkTasQ_Enqueue_SuccessDone-4|3000000|565 ns/op|
|BenchmarkTasQ_Enqueue_CustomErrorDone-4|2000000|609 ns/op|
|BenchmarkTasQ_Enqueue_RetryError3-4|3000000|569 ns/op|