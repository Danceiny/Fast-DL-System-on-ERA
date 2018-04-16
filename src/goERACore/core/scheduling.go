package core

import (
    "encoding/json"
    "fmt"
    "github.com/go-redis/redis"
    "time"
    "bytes"
)

//Input: a new job request {W*T in [A, D), V}
//Output: accept or reject, and price if accepted
//procedure MAKE RESERVATION
//    for each t ∈ [A,D) do
//        demand(t) ← the demand estimate function at t
//        for each i ∈ [1, W] do
//            price(t)(i) ← the highest price p s.t.demand(t)(p) + promised[t] + i > Capacity
//        cost[t] ← price(t)(1) + price(t)(2) + ... + price(t)(W)
//    for each t ∈ [A, D-T] do
//        totalCost[t] ← cost[t] + ... + cost[t+T-1]
//    t* ← arg min(t ∈ [A, D-T])totalCost[t]
//    if V >= totalCost[t*] then
//        schedule the job to start at t*
//        return accept at cost totalCost[t*]
//    else
//        return reject
//end procedure
func BasicEconScheduling(jobRequest *JobRequest) *Response2JobReq {
    defer func() {
        // 处理异常
        if p := recover(); p != nil {
            err, ok := interface{}(p).(error)
            var buff bytes.Buffer
            buff.WriteString("BasicEconScheduling algorithm panic! ")
            if ok {
                buff.WriteString(fmt.Sprintf("error: %s", err.Error()))
            } else {
                buff.WriteString(fmt.Sprintf("clue: %v", p))
            }
            buff.WriteString(")")
            errMsg := buff.String()
            ErrorLog(errMsg)
        }
    }()
    // 检查请求的时间窗口是否合法
    if jobRequest.TwStart.After(jobRequest.TwEnd) {
        jobRequest.TwStart, jobRequest.TwEnd = jobRequest.TwEnd, jobRequest.TwStart
    }
    timeWindowDuration := jobRequest.TwEnd.Sub(jobRequest.TwStart)
    now := time.Now()
    if jobRequest.TwEnd.Before(now) {
        jobRequest.TwEnd = now.Add(timeWindowDuration)
    }
    if jobRequest.TwStart.Before(now) {
        DebugLog("请求的时间窗口已过期")
        jobRequest.TwStart = now
    }

    // 时间窗口的粒度为秒，可能过细（算法复杂度，耗时），可以考虑作调整（时间分片的粒度动态调整？）

    totalCost := make([]uint32, timeWindowDuration/ESTIMATE_INTERVAL+1)
    DebugLog("timeWindowDuration %d second, interval %d second, array with length: %d", timeWindowDuration.Seconds(), ESTIMATE_INTERVAL.Seconds(), len(totalCost))
    //for t := uint64(0); t < timeWindowDuration; t += uint64(ESTIMATE_INTERVAL) {
    //    current_time := jobRequest.TwStart.Add(time.Second * time.Duration(t))
    //    estimateDemand(&current_time, jobRequest.Resources)
    //    totalCost[t] = pricingResourceList(&current_time, jobRequest.Resources)
    //}
    index := 0
    for t := jobRequest.TwStart; t.Before(jobRequest.TwEnd); t = t.Add(time.Duration(ESTIMATE_INTERVAL)) {
        estimateDemand(&now, &t, jobRequest.Resources)
        //DebugLog("index: %d", index)
        totalCost[index] = pricingResourceList(&t, jobRequest.Resources)
        index++
    }
    t := findMinT(totalCost)
    arrivalTime := jobRequest.TwStart.Add(time.Duration(t+1) * ESTIMATE_INTERVAL)
    DebugLog("findMinT got index: %d, so arrivalTime is: %v", t, arrivalTime)
    minTotalPrice := totalCost[t] // 可接受的最低价
    if jobRequest.Value >= minTotalPrice {
        InfoLog("接受")
        return scheduleJob(jobRequest, &arrivalTime, minTotalPrice)
    } else {
        InfoLog("拒绝。最低价:%d, 出价:%d", minTotalPrice, jobRequest.Value)
        return rejectJobRequest(jobRequest, minTotalPrice)
    }
}
func rejectJobRequest(request *JobRequest, v uint32) *Response2JobReq {
    // 响应拒绝
    return &Response2JobReq{
        Id:            request.Id,
        Accepted:      false,
        ArrivalTime:   time.Time{},
        AcceptedPrice: v,
    }
}
func scheduleJob(request *JobRequest, t *time.Time, v uint32) *Response2JobReq {
    // 将作业发布到队列中，等待真正的调度器执行调度
    // TODO:
    alloc := Allocation{
        JobId:     request.JobId,
        Resources: request.Resources,
        TStart:    *t,
        TEnd:      t.Add(request.Duration),
        Value:     v,
    }
    allcName := GetJobAllocRedisKey(alloc.JobId) //重要规则！！！
    msg, err := json.Marshal(alloc)
    if err != nil {
        ErrorLog("marshall allocation failed, reason: %s", err)
    }
    InfoLog("发布任务调度消息 %s", msg)
    // 发布消息（不支持历史查看） ==> cloud
    if err := G_RedisBrokerClient.Publish(REDISACCEPTEDCHANNEL, allcName).Err(); err != nil {
        //panic(err)
        ErrorLog("publish task and resource allocation info failed, reason: %s", err)
    }
    // 添加到队列（有序集合，按启动时间+价值排序，其中启动时间优先排序）
    // TODO：目前该有序set并没有被使用
    if err := G_RedisBackendClient.ZAdd(REDISACCEPTEDSET,
        redis.Z{Score: float64(t.Second()) + convert2Float64LessThanOne(v),
            Member: allcName}).Err(); err != nil {
        //panic(err)
        WarningLog("saved to redis backend failed, reason: %s", err)
    }
    // 分配的详情，由全局可见的persis数据库保存，存储在一个单独的键，键名为allocName
    if err := G_RedisPersisClient.Set(allcName, msg, alloc.TEnd.Sub(time.Now())).Err(); err != nil {
        //panic(err)
        ErrorLog("failed to save allocation detail, reason: %s", err)
    } //启动任务的deadline时刻过期
    // 向发起请求者返回响应 ==> user
    return &Response2JobReq{
        Id:            request.Id,
        Accepted:      true,
        ArrivalTime:   *t,
        AcceptedPrice: v,
    }
}
