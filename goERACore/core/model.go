package core

import (
    "time"
    "github.com/go-redis/redis"
    "fmt"
    "encoding/json"
)

var (
    FRAMEWORKMAP       = make(map[int32]DLFramework)
    FRAMEWORKKEYBYNAME = make(map[string]DLFramework)
)

func init() {
    // 从redis中读共享配置信息
    var client *redis.Client
    client = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "", // no password set
        DB:       0,  // use default DB
    })
    pong, err := client.Ping().Result()
    fmt.Println(pong, err)
    resp := client.ZRangeByScore(REDISFRAMEWORKSET, redis.ZRangeBy{Min: "-inf", Max: "inf", Count: -1})
    for _, item := range resp.Val() {
        // Val() ==> []string
        dlFramework := DLFramework{}
        if err := json.Unmarshal([]byte(item), &dlFramework); err != nil {
            ErrorLog("unmarshall from redis to get framework, reason: %s", err)
            continue
        }
        FRAMEWORKMAP[dlFramework.Id] = dlFramework
        FRAMEWORKKEYBYNAME[dlFramework.Name] = dlFramework
    }
    InfoLog("init from redis %s %s", FRAMEWORKMAP[1].Name, FRAMEWORKKEYBYNAME["tensorflow-1.5:py2"].Name)
}

const (
    // TODO
    CPUFLAG        = 0x0000
    GPUFLAG        = 0x1000
    MEMFLAG        = 0x2000
    FRWFLAG        = 0x3000
    GPUTYPEDEFAULT = 0x0000
    CPUTYPEDEFAULT = 0x0000
    MEMTYPEDEFAULT = 0x0000
    // 框架类型（即名称）由后三个字节表示
    FRWTYPEDEFAULT = 0x0000

    REDISFRAMEWORKSET    = "framework"
    REDISACCEPTEDCHANNEL = "era_accepted_queue"
    REDISACCEPTEDSET     = "era_accepted_set"

    // 计价单位：人民币分
    DEFAULTCPUPRICEPERHOUR = 100
    DEFAULTGPUPRICEPERHOUR = 800
    DEFAULTMEMPRICEPERGB   = 5
)

type ResourceList struct {
    Cpu int32 `json:"cpu"` //four bytes, first byte show cpu, second show cpu type, left show number of Cpu cores
    Gpu int32 `json:"gpu"` //four bytes, first byte show gpu, second show gpu type, left show number of Gpu
    Mem int32 `json:"mem"` //four bytes, first byte show mem, second show mem type, left show GB of memory
    Frw int32 `json:"frw"` //four bytes, first byte show deep learning framework, the left show framework type(Id)
}

func (self *ResourceList) GetCpuNum() int32 {
    return self.Cpu & 0x0011
}
func (self *ResourceList) GetGpuNum() int32 {
    return self.Gpu & 0x0011
}
func (self *ResourceList) GetMemNum() int32 {
    return self.Mem & 0x0011
}
func (self *ResourceList) GetMemType() int32 {
    return self.Mem & 0x0100
}
func (self *ResourceList) GetGpuType() int32 {
    return self.Gpu & 0x0100
}
func (self *ResourceList) GetCpuType() int32 {
    return self.Cpu & 0x0100
}
func (self *ResourceList) GetFrwType() int32 {
    return self.Frw & 0x0111
}
func NewResourceList(cpuType, cpuNum, gpuType, gpuNum, memType, memNum, frwType int32) *ResourceList {
    return &ResourceList{
        Cpu: CPUFLAG | cpuType | cpuNum,
        Mem: MEMFLAG | memType | memNum,
        Gpu: GPUFLAG | gpuType | gpuNum,
        Frw: FRWFLAG | frwType,
    }
}

type ID string
type JobRequest struct {
    Id ID `json:"id"`
    // 预估运行时长
    Duration time.Duration `json:"duration"`
    // 可接受的时间窗口
    TwStart time.Time `json:"tw_start"`
    TwEnd   time.Time `json:"tw_end"`
    // 与job实体绑定的id
    JobId ID `json:"job_id"`
    //竞价价格，单位分
    Value uint32 `json:"value"`
    // 预定资源
    Resources *ResourceList `json:"resources"`
}
type Allocation struct {
    JobId     ID            `json:"job_id"`
    Resources *ResourceList `json:"resources"`
    TStart    time.Time     `json:"t_start"`
    TEnd      time.Time     `json:"t_end"`
}
type Response2JobReq struct {
    Id            ID        `json:"id"`
    Accepted      bool      `json:"accepted"`
    ArrivalTime   time.Time `json:"arrival_time"`
    AcceptedPrice uint32    `json:"accepted_price"`
}
type DataIdsT []ID
type Job struct {
    Id ID `json:"id"`
    // 实际运行时长
    Duration time.Duration `json:"duration"`
    // 创建和修改时间
    TCreated  time.Time `json:"t_created"`
    TModified time.Time `json:"t_modified"`

    // 归属于项目
    ProjectId ID `json:"project_id"`

    // 代码文件
    CodeId ID `json:"code_id"`
    // 数据文件
    DataIds DataIdsT `json:"data_ids"`

    //// 运行环境
    Env *Environment `json:"env"`

    // 启动后的容器入口命令
    EntryCmd string `json:"entry_cmd"`
    // CLI启动命令
    StartCmd string `json:"start_cmd"`

    // 是否开启xx模式
    BTensorboard bool `json:"b_tensorboard"`
    BJupyter     bool `json:"b_jupyter"`

    // （是否）已启动/结束（时间）
    TStarted time.Time `json:"t_started"`
    TEnded   time.Time `json:"t_ended"`

    // 所属用户ID
    Uid ID `json:"uid"`
    // 所属团队ID
    Gid ID `json:"gid"`
    // 权限
    Perm int8 `json:"perm"`
    // 用户描述
    Doc string `json:"doc"`
}

type DLFramework struct {
    Name string `json:"name"` //相当于ID
    // 镜像命名空间
    Namespace string `json:"name_space"`
    // 镜像版本
    CpuImgTag string `json:"cpu_img_tag"`
    GpuImgTag string `json:"gpu_img_tag"`
    // Python版本
    PyVer rune `json:"py_ver"`
    // 唯一int值（由后两个字节表明name，相当于ID，）(
    Id int32 `json:"id"`
    // 基础配置费用(单位：分）
    ConfigPrice uint32 `json:"config_price"`
}
type Environment struct {
    //DlFr    *DLFramework `json:"dl_fr"`
    DlFrName string `json:"dl_fr_name"`
    Os       string `json:"os"` // 操作系统
    WithGpu  bool   `json:"with_gpu"`
}
