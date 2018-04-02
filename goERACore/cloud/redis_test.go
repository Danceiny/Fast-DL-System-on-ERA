package cloud

import (
    "testing"
    . "goERACore/core"
    "github.com/go-redis/redis"
    "encoding/json"
)

func TestRedisPubSub(t *testing.T) {
    //pubSubConn := client.Subscribe("rchat")
    //ch := pubSubConn.Channel()
    //t.Log(ch)
    //for msg := range ch {
    //    t.Log(msg.Payload, msg.Pattern)
    //}
}

func TestInit_Framework_Map(t *testing.T) {
    Init_Framework_Map()
    resp := client.ZRangeByScore(REDIS_FRAMEWORK_SET, redis.ZRangeBy{"-inf", "inf", 0, -1})
    for _, item := range resp.Val() {
        // Val() ==> []string
        dlFramework := DL_Framework{}
        err := json.Unmarshal([]byte(item), &dlFramework)
        if err != nil{
            t.Error(err)
        }
        t.Logf("item.Name is: %s, cputag is: %s",dlFramework.Name, dlFramework.CpuImgTag)
    }

}
