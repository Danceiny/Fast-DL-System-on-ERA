package core

import (
    "encoding/json"
    "github.com/go-redis/redis"
    "time"
)

const (
    // broker: worker需要对接的
    RedisBrokerHost = "test.dl.russellcloud.com"
    RedisBrokerPort = 6380
    RedisBrokerPass = ""
    RedisBrokerDB   = 0
    // backend: era自己的独立后端
    //RedisBackendHost = "localhost"
    RedisBackendHost = "api.cannot.cc"
    RedisBackendPort = 6380
    RedisBackendPass = ""
    RedisBackendDB   = 0
    // config: 整个系统的公用配置
    RedisConfigHost = "api.cannot.cc"
    RedisConfigPort = 6379
    RedisConfigPass = ""
    RedisConfigDB   = 0
    // persis: 整个系统的共用缓存（该版本实现中充当持久化数据库mysql~）
    RedisPersisHost = "api.cannot.cc"
    RedisPersisPort = 6380
    RedisPersisPass = ""
    RedisPersisDB   = 0
)

var (
    G_RedisBrokerClient  *redis.Client
    G_RedisBackendClient *redis.Client
    G_RedisConfigClient  *redis.Client
    G_RedisPersisClient  *redis.Client
    ESTIMATE_INTERVAL    = 60 * time.Second
    FRAMEWORKMAP         = make(map[int32]*DLFramework)
    FRAMEWORKKEYBYNAME   = make(map[string]*DLFramework)
)

func init() {
    G_RedisBrokerClient = redis.NewClient(&redis.Options{
        Addr:     GetEndpointByHostAndName(RedisBrokerHost, RedisBrokerPort),
        Password: RedisBrokerPass,
        DB:       RedisBrokerDB,
    })
    G_RedisBackendClient = redis.NewClient(&redis.Options{
        Addr:     GetEndpointByHostAndName(RedisBackendHost, RedisBackendPort),
        Password: RedisBackendPass,
        DB:       RedisBackendDB,
    })
    G_RedisConfigClient = redis.NewClient(&redis.Options{
        Addr:     GetEndpointByHostAndName(RedisConfigHost, RedisConfigPort),
        Password: RedisConfigPass,
        DB:       RedisConfigDB,
    })
    G_RedisPersisClient = redis.NewClient(&redis.Options{
        Addr:     GetEndpointByHostAndName(RedisPersisHost, RedisPersisPort),
        Password: RedisPersisPass,
        DB:       RedisPersisDB,
    })
}

var frameworkSet = map[string]DLFramework{
    "tensorflow-1.5": {
        Name:        "tensorflow-1.5",
        Namespace:   "tensorflow",
        CpuImgTag:   "1.5.0-py3_aws.22",
        GpuImgTag:   "1.5.0-gpu.cuda9cudnn7-py3_aws.22",
        PyVer:       3,
        Id:          0x0000,
        ConfigPrice: 0,
    },
    "tensorflow-1.5:py2": {
        Name:        "tensorflow-1.5:py2",
        Namespace:   "tensorflow",
        CpuImgTag:   "1.5.0-py2_aws.22",
        GpuImgTag:   "1.5.0-gpu.cuda9cudnn7-py2_aws.22",
        PyVer:       2,
        Id:          0x0001,
        ConfigPrice: 0,
    },
    "caffe": {
        Name:        "caffe",
        Namespace:   "caffe",
        CpuImgTag:   "caffe:1.0-py3.6",
        GpuImgTag:   "caffe:1.0-gpu-py3.6",
        PyVer:       3,
        Id:          0x0002,
        ConfigPrice: 0,
    },
    "caffe2:py2": {
        Name:        "caffe2:py2",
        Namespace:   "caffe",
        CpuImgTag:   "caffe:2.0-py2.7",
        GpuImgTag:   "caffe:2.0-gpu-py2.7",
        PyVer:       2,
        Id:          0x0003,
        ConfigPrice: 0,
    },
    "caffe:py2": {
        Name:        "caffe:py2",
        Namespace:   "caffe",
        CpuImgTag:   "caffe:1.0-py2.6",
        GpuImgTag:   "caffe:1.0-gpu-py2.6",
        PyVer:       2,
        Id:          0x0004,
        ConfigPrice: 0,
    },
    "chainer": {
        Name:        "chainer",
        Namespace:   "chainer",
        CpuImgTag:   "chainer:1.23.0-py3.6",
        GpuImgTag:   "chainer:1.23.0-gpu-py3.9",
        PyVer:       3,
        Id:          0x0005,
        ConfigPrice: 0,
    },
    "chainer-1.23": {
        Name:        "chainer-1.23",
        Namespace:   "chainer",
        CpuImgTag:   "chainer:1.23.0-py3.6",
        GpuImgTag:   "chainer:1.23.0-gpu-py3.9",
        PyVer:       3,
        Id:          0x0006,
        ConfigPrice: 0,
    },
    "chainer-1.23:py2": {
        Name:        "chainer-1.23:py2",
        Namespace:   "chainer",
        CpuImgTag:   "chainer:1.23.0-py2.6",
        GpuImgTag:   "chainer:1.23.0-gpu-py2.9",
        PyVer:       2,
        Id:          7,
        ConfigPrice: 0},
    "chainer-2.0": {
        Name:        "chainer-2.0",
        Namespace:   "chainer",
        CpuImgTag:   "chainer:2.0.0-py3.6",
        GpuImgTag:   "chainer:2.0.0-gpu-py3.9",
        PyVer:       3,
        Id:          8,
        ConfigPrice: 0,
    },
    "chainer-2.0:py2": {
        Name:        "chainer-2.0:py2",
        Namespace:   "chainer",
        CpuImgTag:   "chainer:2.0.0-py2.6",
        GpuImgTag:   "chainer:2.0.0-gpu-py2.9",
        PyVer:       2,
        Id:          9,
        ConfigPrice: 0,
    },
    "chainer:py2": {
        Name:        "",
        Namespace:   "",
        CpuImgTag:   "chainer:1.23.0-py2.6",
        GpuImgTag:   "chainer:1.23.0-gpu-py2.9",
        PyVer:       0,
        Id:          0x000a,
        ConfigPrice: 0,
    },
    "keras": {
        Name:        "keras",
        Namespace:   "keras",
        CpuImgTag:   "tensorflow:1.1.0-py3_aws.7",
        GpuImgTag:   "tensorflow:1.1.0-gpu-py3_aws.7",
        PyVer:       3,
        Id:          0,
        ConfigPrice: 0,
    },
    "keras:py2": {
        Name:        "keras:py2",
        Namespace:   "keras",
        CpuImgTag:   "tensorflow:1.1.0-py2_aws.7",
        GpuImgTag:   "tensorflow:1.1.0-gpu-py2_aws.7",
        PyVer:       2,
        Id:          0x000b,
        ConfigPrice: 0,
    },

    "kur": {
        Name:        "",
        Namespace:   "",
        CpuImgTag:   "kur:0.6.0-py3.6",
        GpuImgTag:   "kur:0.6.0-gpu-py3.6",
        PyVer:       0,
        Id:          0x000c,
        ConfigPrice: 0,
    },
    "mxnet-1.0:py2":       {CpuImgTag: "mxnet:1.0.0-py2", GpuImgTag: "mxnet:1.0.0-gpu-py2", Id: 0x000d},
    "mxnet:py2":           {CpuImgTag: "mxnet:0.10.0-py2.6", GpuImgTag: "mxnet:0.10.0-gpu-py2.6", Id: 0x000e},
    "paddle:py2":          {CpuImgTag: "paddle:latest", GpuImgTag: "paddle:latest-gpu", Id: 0x000f},
    "pytorch":             {CpuImgTag: "pytorch:0.2.0-py3.15", GpuImgTag: "pytorch:0.2.0-gpu-py3.15", Id: 0x0010},
    "pytorch-0.1":         {CpuImgTag: "pytorch:0.1-py3.8", GpuImgTag: "pytorch:0.1-gpu-py3.8", Id: 0x0011},
    "pytorch-0.1:py2":     {CpuImgTag: "pytorch:0.1-py2.8", GpuImgTag: "pytorch:0.1-gpu-py2.8", Id: 0x0012},
    "pytorch-0.2":         {CpuImgTag: "pytorch:0.2.0-py3.15", GpuImgTag: "pytorch:0.2.0-gpu-py3.15", Id: 0x0013},
    "pytorch-0.2:py2":     {CpuImgTag: "pytorch:0.2.0-py2.15", GpuImgTag: "pytorch:0.2.0-gpu-py2.15", Id: 0x0014},
    "pytorch-0.3":         {CpuImgTag: "pytorch:0.3.0-py3.17", GpuImgTag: "pytorch:0.3.0-gpu.cuda8cudnn6-py3.17", Id: 0x0015},
    "pytorch-0.3:py2":     {CpuImgTag: "pytorch:0.3.0-py2.17", GpuImgTag: "pytorch:0.3.0-gpu.cuda8cudnn6-py2.17", Id: 0x0016},
    "pytorch:py2":         {CpuImgTag: "pytorch:0.2.0-py2.15", GpuImgTag: "pytorch:0.2.0-gpu-py2.15", Id: 0x0017},
    "tensorflow":          {CpuImgTag: "tensorflow:1.1.0-py3_aws.7", GpuImgTag: "tensorflow:1.1.0-gpu-py3_aws.7", Id: 0x0018},
    "tensorflow-0.12":     {CpuImgTag: "tensorflow:0.12.1-py3.6", GpuImgTag: "tensorflow:0.12.1-gpu-py3.6", Id: 0x0019},
    "tensorflow-0.12:py2": {CpuImgTag: "tensorflow:0.12.1-py2.6", GpuImgTag: "tensorflow:0.12.1-gpu-py2.6", Id: 0x001a},
    "tensorflow-1.0":      {CpuImgTag: "tensorflow:1.0.1-py3_aws.7", GpuImgTag: "tensorflow:1.0.1-gpu-py3.7", Id: 0x001b},
    "tensorflow-1.0:py2":  {CpuImgTag: "tensorflow:1.0.1-py2_aws.7", GpuImgTag: "tensorflow:1.0.1-gpu-py2.7", Id: 0x001c},
    "tensorflow-1.1":      {CpuImgTag: "tensorflow:1.1.0-py3_aws.7", GpuImgTag: "tensorflow:1.1.0-gpu-py3_aws.7", Id: 0x001d},
    "tensorflow-1.1:py2":  {CpuImgTag: "tensorflow:1.1.0-py2_aws.7", GpuImgTag: "tensorflow:1.1.0-gpu-py2_aws.7", Id: 0x001e},
    "tensorflow-1.2":      {CpuImgTag: "tensorflow:1.2.1-py3_aws.7", GpuImgTag: "tensorflow:1.2.1-gpu-py3_aws.7", Id: 0x001f},
    "tensorflow-1.2:py2":  {CpuImgTag: "tensorflow:1.2.1-py2_aws.7", GpuImgTag: "tensorflow:1.2.1-gpu-py2_aws.7", Id: 0x0020},
    "tensorflow-1.3":      {CpuImgTag: "tensorflow:1.3.1-py3_aws.13", GpuImgTag: "tensorflow:1.3.1-gpu-py3_aws.13", Id: 0x0021},
    "tensorflow-1.3:py2":  {CpuImgTag: "tensorflow:1.3.1-py2_aws.13", GpuImgTag: "tensorflow:1.3.1-gpu-py2_aws.13", Id: 0x0022},
    "tensorflow-1.4":      {CpuImgTag: "tensorflow:1.4.0-py3_aws.14", GpuImgTag: "tensorflow:1.4.0-gpu-py3_aws.14", Id: 0x0023},
    "tensorflow-1.4:py2":  {CpuImgTag: "tensorflow:1.4.0-py2_aws.14", GpuImgTag: "tensorflow:1.4.0-gpu-py2_aws.14", Id: 0x0024},
    "tensorflow-1.6":      {CpuImgTag: "dockfiles-danceiny:1.6.0-py3", Id: 0x0025},
    "tensorflow-1.6:py2":  {CpuImgTag: "dockfiles-danceiny:1.6.0-py2", Id: 0x0026},
    "tensorflow:py2":      {CpuImgTag: "tensorflow:1.1.0-py2_aws.7", GpuImgTag: "tensorflow:1.1.0-gpu-py2_aws.7", Id: 0x0027},
    "theano":              {CpuImgTag: "theano:0.9.0-py3.6", GpuImgTag: "theano:0.9.0-gpu-py3.6", Id: 0x0028},
    "theano-0.8":          {CpuImgTag: "theano:0.8.2-py3.6", GpuImgTag: "theano:0.8.2-gpu-py3.6", Id: 0x0029},
    "theano-0.8:py2":      {CpuImgTag: "theano:0.8.2-py2.6", GpuImgTag: "theano:0.8.2-gpu-py2.6", Id: 0x002a},
    "theano-0.9":          {CpuImgTag: "theano:0.9.0-py3.6", GpuImgTag: "theano:0.9.0-gpu-py3.6", Id: 0x002b},
    "theano-0.9:py2":      {CpuImgTag: "theano:0.9.0-py2.6", GpuImgTag: "theano:0.9.0-gpu-py2.6", Id: 0x002c},
    "theano:py2":          {CpuImgTag: "theano:0.9.0-py2.6", GpuImgTag: "theano:0.9.0-gpu-py2.6", Id: 0x002d},
    "torch":               {CpuImgTag: "torch:7-py3.6", GpuImgTag: "torch:7-gpu-py3.6", Id: 0x002e},
    "torch:py2":           {CpuImgTag: "torch:7-py2.6", GpuImgTag: "torch:7-gpu-py2.6", Id: 0x002f},
}

func InitFrameworkMap() {
    framework := make(map[string]interface{})
    for k, v := range frameworkSet {
        G_RedisBrokerClient.ZAdd(REDISFRAMEWORKSET_WITHSCORE, redis.Z{0, k})
        framework[k], _ = json.Marshal(v)
    }
    if val, err := G_RedisPersisClient.HMSet(REDISFRAMEWORKSET, framework).Result(); err != nil {
        ErrorLog("hmset redisframeworkset failed, reason: %s; val: %s", err, val)
    }
}
