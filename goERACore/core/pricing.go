package core

import (
    "time"
    "reflect"
    "encoding/json"
)

// 价格单位分，数据类型uint32 Range: 0 through 4294967295.
// 因此最大价格为 4294,9672.95元人民币

func pricingCpu(t *time.Time) uint32 {
    return DEFAULTCPUPRICEPERHOUR
}
func pricingGpu(t *time.Time) uint32 {
    return DEFAULTGPUPRICEPERHOUR
}
func pricingMem(t *time.Time) uint32 {
    return DEFAULTMEMPRICEPERGB
}
func pricingFramework(t *time.Time, f int32) uint32 {
    freqAtT := 0.1 // 在t时刻，该框架被使用的频率
    default_price := FRAMEWORKMAP[f].ConfigPrice
    //TODO
    // 给该框架加上一分热度
    v, _ := json.Marshal(FRAMEWORKMAP[f])
    client.ZIncrBy(REDISFRAMEWORKSET, 1, string(v))
    score := client.ZScore(REDISFRAMEWORKSET, string(v)).Val()
    // score is float64, should care type-cast
    default_price += uint32(freqAtT * score / float64(client.ZCard(REDISFRAMEWORKSET).Val()))
    return default_price
}
func pricing(t *time.Time, res int32) uint32 {
    var price uint32
    switch res & 0x1000 {
    case CPUFLAG:
        price = pricingCpu(t)
    case GPUFLAG:
        price = pricingGpu(t)
    case MEMFLAG:
        price = pricingMem(t)
    case FRWFLAG:
        price = pricingFramework(t, res&0x0111)
    }
    return price

}
func pricingResourceList(t *time.Time, rl *ResourceList) uint32 {
    //PrintLog("debug", "resource list is: %v", *rl)
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
