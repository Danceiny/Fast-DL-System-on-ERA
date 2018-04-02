package core

import (
    "time"
    "reflect"
)

// 价格单位分，数据类型uint32 Range: 0 through 4294967295.
// 因此最大价格为 4294,9672.95元人民币

func pricingCpu(t *time.Time) uint32 {
    return DEFAULT_CPU_PRICE_PER_HOUR
}
func pricingGpu(t *time.Time) uint32 {
    return DEFAULT_GPU_PRICE_PER_HOUR
}
func pricingMem(t *time.Time) uint32 {
    return DEFAULT_MEM_PRICE_PER_GB
}
func pricingFramework(t *time.Time, f int32) uint32 {
    default_price := FRAMEWORK_MAP[f].ConfigPrice
    //TODO
    return default_price
}
func pricing(t *time.Time, res int32) uint32 {
    var price uint32
    switch res & 0x1000 {
    case CPU_FLAG:
        price = pricingCpu(t)
    case GPU_FLAG:
        price = pricingGpu(t)
    case MEM_FLAG:
        price = pricingMem(t)
    case FRW_FLAG:
        price = pricingFramework(t, res&0x0111)
    }
    return price

}
func pricingResourceList(t *time.Time, rl *Resource_List) uint32 {
    PrintLog("debug", "resource list is: %v", *rl)
    v := reflect.ValueOf(*rl)
    count := v.NumField()
    costs := make([]uint32, count)
    var cost uint32 = 0
    for i := 0; i < count; i++ {
        costs[i] = pricing(t, int32(v.Field(i).Int()))
        cost += costs[i]
    }
    return cost
}
