package cloud

import (
	"encoding/json"
	"github.com/go-redis/redis"
	. "goERACore/core"
	"testing"
)

// Redis lib usage:
//func ExampleClient() {
//    err := redisClient.Set("key", "value", 0).Err()
//    if err != nil {
//        panic(err)
//    }
//
//    val, err := redisClient.Get("key").Result()
//    if err != nil {
//        panic(err)
//    }
//    fmt.Println("key", val)
//
//    val2, err := redisClient.Get("key2").Result()
//    if err == redis.Nil {
//        fmt.Println("key2 does not exist")
//    } else if err != nil {
//        panic(err)
//    } else {
//        fmt.Println("key2", val2)
//    }
//    // Output: key value
//    // key2 does not exist
//}

func TestRedisPubSub(t *testing.T) {
	//pubSubConn := redisClient.Subscribe("rchat")
	//ch := pubSubConn.Channel()
	//t.Log(ch)
	//for msg := range ch {
	//    t.Log(msg.Payload, msg.Pattern)
	//}
}

func TestInit_Framework_Map(t *testing.T) {
	Init_Framework_Map()
	resp := redisClient.ZRangeByScore(REDISFRAMEWORKSET, redis.ZRangeBy{"-inf", "inf", 0, -1})
	for _, item := range resp.Val() {
		// Val() ==> []string
		dlFramework := DLFramework{}
		err := json.Unmarshal([]byte(item), &dlFramework)
		if err != nil {
			t.Error(err)
		}
		t.Logf("item.Name is: %s, cputag is: %s", dlFramework.Name, dlFramework.CpuImgTag)
	}

}
