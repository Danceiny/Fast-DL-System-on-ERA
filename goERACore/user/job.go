package user

import "time"

type ID int64
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
    DataId ID `json:"data_id"`

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
    UId ID `json:"u_id"`
    // 所属团队ID
    GId ID `json:"g_id"`
    // 权限
    Perm int8 `json:"perm"`
    // 用户描述
    Doc string `json:"doc"`
}

type DL_Framework struct {
    Name string `json:"name"` //相当于ID
    // 镜像命名空间
    NameSpace string `json:"name_space"`
    // 镜像版本
    ImgTag string `json:"img_tag"`
    // Python版本
    PyVer rune `json:"py_ver"`
}
type Environment struct {
    DlFr  *DL_Framework `json:"dl_fr"`
    OsVer string        `json:"os_ver"` // 操作系统
}

// 竞价描述语言
type Bid_Desc_Lang struct {
}
