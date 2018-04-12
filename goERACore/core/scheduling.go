package core

import (
    "time"
    "encoding/json"
    "github.com/go-redis/redis"
    "fmt"
)

var ESTIMATE_INTERVAL = 60
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
    // 时间窗口的粒度为秒，可能过细（算法复杂度，耗时），可以考虑作调整（时间分片的粒度动态调整？）
    timeWindowDuration := uint64((jobRequest.TwEnd.Sub(jobRequest.TwStart)).Seconds())
    PrintLog("debug", "timeWindowDuration %d second", timeWindowDuration)
    totalCost := make([]uint32, timeWindowDuration)
    //for t := uint64(0); t < timeWindowDuration; t += uint64(ESTIMATE_INTERVAL) {
    //    current_time := jobRequest.TwStart.Add(time.Second * time.Duration(t))
    //    estimateDemand(&current_time, jobRequest.Resources)
    //    totalCost[t] = pricingResourceList(&current_time, jobRequest.Resources)
    //}
    for t := jobRequest.TwStart; t.Before(jobRequest.TwEnd); t = t.Add(time.Duration(ESTIMATE_INTERVAL)) {
        estimateDemand(&t, jobRequest.Resources)
        totalCost[t.Sub(jobRequest.TwStart)] = pricingResourceList(&t, jobRequest.Resources)
    }
    t := findMinT(totalCost)
    PrintLog("debug", "findMinT is: %d", t)
    minTotalPrice := totalCost[t] // 可接受的最低价
    if jobRequest.Value >= minTotalPrice {
        InfoLog("接受")
        start_time := jobRequest.TwStart.Add(time.Second * time.Duration(t))
        return scheduleJob(jobRequest, &start_time, minTotalPrice)
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
    }
    allcName := fmt.Sprintf("accepted_%s", alloc.JobId)
    msg, _ := json.Marshal(alloc)
    InfoLog("发布任务调度消息")
    // 发布消息（不支持历史查看） ==> cloud
    client.Publish(REDISACCEPTEDCHANNEL, allcName)
    // 添加到队列（有序集合，按启动时间+价值排序，其中启动时间优先排序）
    client.ZAdd(REDISACCEPTEDSET,
        redis.Z{Score: float64(t.Second()) + convert2Float64LessThanOne(v),
            Member: allcName})
    // 分配的详情存储在一个单独的键，键名为allocName
    client.Set(allcName, msg, alloc.TEnd.Sub(time.Now())) //启动任务的deadline时刻过期
    // 向发起请求者返回响应 ==> user
    return &Response2JobReq{
        Id:            request.Id,
        Accepted:      true,
        ArrivalTime:   *t,
        AcceptedPrice: v,
    }
}
func findMaxT(arr []uint32) uint64 {
    max := uint32(0)
    maxI := uint64(0)
    for i, v := range arr {
        if v > max {
            maxI = uint64(i)
        }
    }
    return maxI
}
func findMinT(arr []uint32) uint64 {
    // TODO: 优化
    min := arr[0]
    minI := 0
    for i, v := range arr {
        if v < min {
            minI = i
        }
    }
    return uint64(minI)
}

func convert2Float64LessThanOne(n uint32) float64 {
    // 1331 ==> 0.1331
    // Algorithm Reference:
    // https://stackoverflow.com/questions/701322/how-can-you-get-the-first-digit-in-an-int-c/701621#701621
    i := float64(n)
    if i >= 100000000 {
        i /= 100000000
    }
    if i >= 10000 {
        i /= 10000
    }
    if i >= 100 {
        i /= 100
    }
    if i >= 10 {
        i /= 10
    }
    return i / 10
}
