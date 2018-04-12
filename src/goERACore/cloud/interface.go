package cloud

import (
    "github.com/go-redis/redis"
    //"encoding/json"
    //. "goERACore/core"
)

var client *redis.Client

func init() {
    client = redis.NewClient(&redis.Options{
        Addr:     "localhost:6381",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
}

// This is the main interface with the actual cloud scheduler. The cloud should repeatedly
// call this method (quite often, say, every few seconds) and ask ERA for the current
// allocations to be made.（For performance, it is also possible to replace this query with
// an event-driven scheme in which ERA pushes an event to the cloud scheduler when the
// allocations change.）The method returns an allocation, which is the list of jobs that
// should be instantaneously allocated resources and the resources that should be allocated
// to them. In the simple case of a single resource, it is a list of “job J should now be
// getting W resources.” The actual cloud infrastructure should update the resources that
// it currently allocates to all jobs to fit the results of the current allocation returned
// by this query. This new allocation remains in effect until a future query returns a
// different allocation. It is the responsibility of the underlying cloud scheduling system
// to query ERA often enough, and to put these new allocations into effect ASAP, so that any
// changes are effected with reasonably small delay. The main responsibility of the ERA
// system is to ensure that the sequence of answers to this query reflects a plan that can
// accommodate all accepted reservation requests.
//
// The main architectural aspect of this query is to make the interface between ERA and
// the cloud system narrow, such that it completely hides the plan ERA has for future
// allocation. It is assumed that the cloud has no information on the total requirements
// of the jobs, and follows ERA as accurately as possible.

// 实际的云调度器需要反复调用该接口。
// 由于RussellCloud基于Python语言（RC的云调度器实际上是Celery），在实现Go版本的阿里云封装之前，可以考虑将本函数
// 封装为一个TCP接口，由Celery进行请求调用
// 另：Redis内建了发布订阅机制，该函数是否可以弃用：在接受请求后将分配发布到redis，celery订阅。

// 具体的allocation，应该是一个简洁的数据结构，并且可以存放于Redis
//func getCurrentAllocation() {
//    //该方法返回一个分配，该分配是应当立即分配资源的作业的列表以及应分配给它们的资源。
//    //在单一资源的简单情况下，它是“作业J现在应获取W资源”的列表。
//    //实际云基础架构应更新当前分配给所有作业的资源，以适应当前分配结果这个查询。
//    //这个新分配保持有效，直到将来的查询返回不同的分配。底层云调度系统有责任经常查询ERA，
//    //并尽快将这些新的分配生效，以便任何更改都以合理的小延迟进行。
//    //ERA系统的主要职责是确保该查询的答案顺序反映了可以适应所有接受的预订请求的计划。
//    pubSubConn := client.Subscribe(REDISACCEPTEDCHANNEL)
//    ch := pubSubConn.Channel()
//    for msg := range ch {
//        var allocation Allocation
//        json.Unmarshal([]byte(msg.Payload), &allocation)
//    }
//
//}
