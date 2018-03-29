package core

import "time"

const (
    // TODO
    CPU_FLAG         = 0x0000
    GPU_FLAG         = 0x1000
    MEM_FLAG         = 0x2000
    GPU_TYPE_DEFAULT = 0x0000
    CPU_TYPE_DEFAULT = 0x0000
    MEM_TYPE_DEFAULT = 0x0000

    REDIS_ACCEPTED_CHANNEL = "queue"
    REDIS_ACCEPTED_SET = "queue"

    // 计价单位：人民币分
    DEFAULT_CPU_PRICE_PER_HOUR = 100
    DEFAULT_GPU_PRICE_PER_HOUR = 800
    DEFAULT_MEM_PRICE_PER_GB   = 5
)

type Resource_List struct {
    Cpu int32 //two bytes, first byte show cpu, second show cpu type, left show number of Cpu cores
    Gpu int32 //two bytes, first byte show gpu, second show gpu type, left show number of Gpu
    Mem int32 //four bytes, first byte show mem, second show mem type, left show GB of memory
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

func NewResourceList(cpuType, cpuNum, gpuType, gpuNum, memType, memNum int32) *Resource_List {
    return &Resource_List{
        Cpu: CPU_FLAG | cpuType | cpuNum,
        Mem: MEM_FLAG | memType | memNum,
        Gpu: GPU_FLAG | gpuType | gpuNum,
    }
}

type ID int64
type JobRequest struct {
    Id ID
    // 预估运行时长
    Duration time.Duration
    // 可接受的时间窗口
    TwStart time.Time
    TwEnd   time.Time
    // 与job实体绑定的id
    JobId ID
    //竞价价格，单位分
    Value uint32
    // 预定资源
    Resources *Resource_List
}
type Allocation struct {
    JobId     ID
    Resources *Resource_List
    TStart    time.Time
    TEnd      time.Time
}
type Response2JobReq struct {
    Id            ID
    Accepted      bool
    ArrivalTime   time.Time
    AcceptedPrice uint32
}
