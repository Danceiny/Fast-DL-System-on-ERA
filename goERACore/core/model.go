package core

import (
    "time"
    "github.com/go-redis/redis"
    "fmt"
    "encoding/json"
)

var (
    FRAMEWORK_MAP DL_Framework_Map = make(DL_Framework_Map)
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
    resp := client.ZRangeByScore(REDIS_FRAMEWORK_SET, redis.ZRangeBy{"-inf", "inf", 0, -1})
    for _, item := range resp.Val() {
        // Val() ==> []string
        dlFramework := DL_Framework{}
        json.Unmarshal([]byte(item), dlFramework)
        FRAMEWORK_MAP[dlFramework.Id] = dlFramework
    }
}

const (
    // TODO
    CPU_FLAG         = 0x0000
    GPU_FLAG         = 0x1000
    MEM_FLAG         = 0x2000
    FRW_FLAG         = 0x3000
    GPU_TYPE_DEFAULT = 0x0000
    CPU_TYPE_DEFAULT = 0x0000
    MEM_TYPE_DEFAULT = 0x0000
    // 框架类型（即名称）由后三个字节表示
    FRW_TYPE_DEFAULT = 0x0000

    REDIS_FRAMEWORK_SET    = "framework"
    REDIS_ACCEPTED_CHANNEL = "era_accepted_queue"
    REDIS_ACCEPTED_SET     = "era_accepted_set"

    // 计价单位：人民币分
    DEFAULT_CPU_PRICE_PER_HOUR = 100
    DEFAULT_GPU_PRICE_PER_HOUR = 800
    DEFAULT_MEM_PRICE_PER_GB   = 5
)

var ()

type Resource_List struct {
    Cpu int32 `json:"cpu"` //four bytes, first byte show cpu, second show cpu type, left show number of Cpu cores
    Gpu int32 `json:"gpu"` //four bytes, first byte show gpu, second show gpu type, left show number of Gpu
    Mem int32 `json:"mem"` //four bytes, first byte show mem, second show mem type, left show GB of memory
    Frw int32 `json:"frw"` //four bytes, first byte show deep learning framework, the left show framework type(Id)
}

func (self *Resource_List) GetCpuNum() int32 {
    return self.Cpu & 0x0011
}
func (self *Resource_List) GetGpuNum() int32 {
    return self.Gpu & 0x0011
}
func (self *Resource_List) GetMemNum() int32 {
    return self.Mem & 0x0011
}
func (self *Resource_List) GetMemType() int32 {
    return self.Mem & 0x0100
}
func (self *Resource_List) GetGpuType() int32 {
    return self.Gpu & 0x0100
}
func (self *Resource_List) GetCpuType() int32 {
    return self.Cpu & 0x0100
}
func (self *Resource_List) GetFrwType() int32 {
    return self.Frw & 0x0111
}
func NewResourceList(cpuType, cpuNum, gpuType, gpuNum, memType, memNum, frwType int32) *Resource_List {
    return &Resource_List{
        Cpu: CPU_FLAG | cpuType | cpuNum,
        Mem: MEM_FLAG | memType | memNum,
        Gpu: GPU_FLAG | gpuType | gpuNum,
        Frw: FRW_FLAG | frwType,
    }
}

type ID int64
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
    Resources *Resource_List `json:"resources"`
}
type Allocation struct {
    JobId     ID             `json:"job_id"`
    Resources *Resource_List `json:"resources"`
    TStart    time.Time      `json:"t_start"`
    TEnd      time.Time      `json:"t_end"`
}
type Response2JobReq struct {
    Id            ID        `json:"id"`
    Accepted      bool      `json:"accepted"`
    ArrivalTime   time.Time `json:"arrival_time"`
    AcceptedPrice uint32    `json:"accepted_price"`
}
type Data_Ids_T []ID
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
    DataIds Data_Ids_T `json:"data_ids"`

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

type DL_Framework struct {
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
    DlFr    *DL_Framework `json:"dl_fr"`
    Os      string        `json:"os"` // 操作系统
    WithGpu bool          `json:"with_gpu"`
}
type DL_Framework_Map map[int32]DL_Framework

// 竞价描述语言
type Bid_Desc_Lang struct {
}
