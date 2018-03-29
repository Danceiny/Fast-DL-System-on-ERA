package core

import (
    "time"
    "encoding/binary"
    "goERACore/cloud"
    "reflect"
)

// 价格单位分，数据类型uint32 Range: 0 through 4294967295.
// 因此最大价格为 4294,9672.95元人民币

const (
    // 计价单位：人民币分
    DEFAULT_CPU_PRICE_PER_HOUR = 100
    DEFAULT_GPU_PRICE_PER_HOUR = 800
    DEFAULT_MEM_PRICE_PER_GB   = 5
)

func pricingCpu(t *time.Time) uint32 {
    return DEFAULT_CPU_PRICE_PER_HOUR
}
func pricingGpu(t *time.Time) uint32 {
    return DEFAULT_GPU_PRICE_PER_HOUR
}
func pricingMem(t *time.Time) uint32 {
    return DEFAULT_MEM_PRICE_PER_GB
}
func pricing(t *time.Time, r []byte) uint32 {
    res := int32(binary.BigEndian.Uint32(r))
    var price uint32
    switch res & 0x1000 {
    case cloud.CPU_FLAG:
        price = pricingCpu(t)
    case cloud.GPU_FLAG:
        price = pricingGpu(t)
    case cloud.MEM_FLAG:
        price = pricingMem(t)
    }
    return price

}

//    for each t ∈ [A,D) do
//        demand(t) ← the demand estimate function at t
//        for each i ∈ [1, W] do
//            price(t)(i) ← the highest price p s.t.demand(t)(p) + promised[t] + i > Capacity
//        cost[t] ← price(t)(1) + price(t)(2) + ... + price(t)(W)
//    for each t ∈ [A, D-T] do
//        totalCost[t] ← cost[t] + ... + cost[t+T-1]
func pricingResourceList(t *time.Time, rl *cloud.Resource_List) uint32 {
    v := reflect.ValueOf(rl)
    count := v.NumField()
    costs := make([]uint32, count)
    var cost uint32 = 0
    for i := 0; i < count; i++ {
        costs[i] = pricing(t, v.Field(i).Bytes())
        //costs[i] = pricing(t, v.Field(i).Int()) => int64
        cost += costs[i]
    }
    return cost
}
