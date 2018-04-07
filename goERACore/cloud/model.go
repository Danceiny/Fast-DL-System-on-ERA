package cloud

import (
    . "goERACore/core"
    "github.com/go-redis/redis"
    "encoding/json"
)

//var FRAMEWORKMAP = make(map[int32]DLFramework)
var frameworkSet = map[int32]DLFramework{
    0x0000: {
        Name:        "tensorflow-1.5",
        Namespace:   "tensorflow",
        CpuImgTag:   "1.5.0-py3_aws.22",
        GpuImgTag:   "1.5.0-gpu.cuda9cudnn7-py3_aws.22",
        PyVer:       3,
        Id:          0x0000,
        ConfigPrice: 1,
    },
    0x0001: {
        Name:        "tensorflow-1.5:py2",
        Namespace:   "tensorflow",
        CpuImgTag:   "1.5.0-py2_aws.22",
        GpuImgTag:   "1.5.0-gpu.cuda9cudnn7-py2_aws.22",
        PyVer:       2,
        Id:          0x0001,
        ConfigPrice: 2,
    },
    0x00ff: {
        Name:        "caffe2:py2",
        Namespace:   "caffe",
        CpuImgTag:   "caffe:2.0-py2.7",
        GpuImgTag:   "caffe:2.0-gpu-py2.7",
        PyVer:       2,
        Id:          0x00ff,
        ConfigPrice: 3,
    },
}

func Init_Framework_Map() {
    for _, v := range frameworkSet {
        bytes, _ := json.Marshal(v)
        client.ZAdd(REDISFRAMEWORKSET, redis.Z{0, string(bytes)})
    }
}
