package cloud

import (
    //"github.com/go-redis/redis"
    //"fmt"
    "testing"
)

//var client *redis.Client

//func init() {
//    client = redis.NewClient(&redis.Options{
//        Addr:     "localhost:6379",
//        Password: "", // no password set
//        DB:       0,  // use default DB
//    })
//    pong, err := client.Ping().Result()
//    fmt.Println(pong, err)
//}
func TestRedisPubSub(t *testing.T) {
    pubSubConn := client.Subscribe("rchat")
    //msg, err := pubSubConn.Receive()
    //var ch = make([]chan *redis.Message, 0)
    ch := pubSubConn.Channel()
    t.Log(ch)
    for msg := range ch {
        t.Log(msg.Payload, msg.Pattern)
    }
}

//func main(){
//}