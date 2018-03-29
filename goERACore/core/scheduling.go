package core

import "time"
import (
    "goERACore/cloud"
)

type ID int64
type JobRequest struct {
    Id ID
    // 预估运行时长
    Duration time.Duration
    // 可接受的时间窗口
    TwStart time.Time
    TwEnd   time.Time
    // 与job实体绑定的id
    JobId ID
    //竞价价格，单位分
    Value uint32
    // 预定资源
    Resources *cloud.Resource_List
}

type Response2JobReq struct {
    Id            ID
    Accepted      bool
    ArrivalTime   time.Time
    AcceptedPrice uint32
}

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

    totalCost := make([]uint32, timeWindowDuration)
    for t := uint64(0); t < timeWindowDuration; t++ {
        current_time := jobRequest.TwStart.Add(time.Second * time.Duration(t))
        estimateDemand(&current_time, -1)
        totalCost[t] = pricingResourceList(&current_time, jobRequest.Resources)
    }
    t := findMinT(totalCost)
    if jobRequest.Value >= totalCost[t] {
        start_time := jobRequest.TwStart.Add(time.Second * time.Duration(t))
        return scheduleJob(jobRequest, &start_time)
    } else {
        return rejectJobRequest(jobRequest)
    }
}
func rejectJobRequest(request *JobRequest) *Response2JobReq {
    // 响应拒绝
    return &Response2JobReq{
        Id:            request.Id,
        Accepted:      false,
        ArrivalTime:   time.Time{},
        AcceptedPrice: 0,
    }
}
func scheduleJob(request *JobRequest, t *time.Time) *Response2JobReq {
    // 实际的调度器，将作业加入等候队列
    arrivalTime := time.Time{}
    acceptedPrice := uint32(0)
    return &Response2JobReq{
        Id:            request.Id,
        Accepted:      true,
        ArrivalTime:   arrivalTime,
        AcceptedPrice: acceptedPrice,
    }
}
func findMaxT(arr []uint32) uint64 {
    max := uint32(0)
    maxI := uint64(-1)
    for i, v := range arr {
        if v > max {
            maxI = uint64(i)
        }
    }
    return maxI
}
func findMinT(arr []uint32) uint64 {
    // TODO: 优化
    min := uint32(0)
    minI := uint64(-1)
    for i, v := range arr {
        if v < min {
            minI = uint64(i)
        }
    }
    return minI
}
