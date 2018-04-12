package core

import (
    "time"
    "reflect"
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
func pricingFrw(t *time.Time, f int32) uint32 {
    freqAtT := getFrwFreqAtTime(t, f) // 在t时刻，该框架被使用的频率
    default_price := FRAMEWORKMAP[f].ConfigPrice
    //TODO
    // 给该框架加上一分热度
    strV := FRAMEWORKMAP[f].Name
    //score := redisClient.ZScore(REDISFRAMEWORKSET_WITHSCORE, strV).Val()
    redisClient.ZIncrBy(REDISFRAMEWORKSET_WITHSCORE, 1, strV)
    rank := redisClient.ZRank(REDISFRAMEWORKSET_WITHSCORE, strV).Val()
    count := redisClient.ZCard(REDISFRAMEWORKSET_WITHSCORE).Val()
    // score is float64, should care type-cast
    default_price = uint32(float32(default_price) * (1 + freqAtT*(1+float32(rank/count)))) // 配置费用计算公式
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
        price = pricingFrw(t, res&0x0111)
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
