package main

import (
    "github.com/go-redis/redis"
    "fmt"
)

var client *redis.Client

func init() {
    client = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    pong, err := client.Ping().Result()
    fmt.Println(pong, err)
}
func TestRedisPubSub() {
    pubSubConn := client.Subscribe("rchat")
    ch := pubSubConn.Channel()
    for msg := range ch {
        print(msg.Payload, msg.Pattern, '\n')
    }
}
func convert2Float64LessThanOne(n int) float64 {
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
func main() {
    //TestRedisPubSub()
    for i := 0; i < 10; i += 1 {
        r := convert2Float64LessThanOne(i * 2)
        println(i*2, r)
    }
    r := convert2Float64LessThanOne(120000000000)
    println("3个0", r)
}
