package cloud

const (
    // TODO
    CPU_FLAG         = 0x0000
    GPU_FLAG         = 0x1000
    MEM_FLAG         = 0x2000
    GPU_TYPE_DEFAULT = 0x0000
    CPU_TYPE_DEFAULT = 0x0000
    MEM_TYPE_DEFAULT = 0x0000
)

type Resource_List struct {
    Cpu int32 //two bytes, first byte show cpu, second show cpu type, left show number of Cpu cores
    Gpu int32 //two bytes, first byte show gpu, second show gpu type, left show number of Gpu
    Mem int32 //four bytes, first byte show mem, second show mem type, left show GB of memory
}

func init() {
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
