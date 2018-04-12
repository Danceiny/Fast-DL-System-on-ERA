package core

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/go-redis/redis"
    _ "github.com/go-sql-driver/mysql"
    "time"
)

func estimateDemand(t *time.Time, r *ResourceList) {
    // t时刻的软件环境需求
    estimateFrwDemand(t, r.Frw)

    // 获取当前时刻（不一定是t时刻，通常早于t时刻）的正在使用的硬件资源情况
    cur := time.Now()
    rRes := getRunningSourcesAtTimeT(&cur)

    //从已分配的历史队列中估计当前时刻t的运行任务列表，从而推算时刻t的资源需求
    //分配已执行的，将分数置为负值，因此范围的最小值为0
    if resp, err := redisClient.ZRangeByScore(REDISACCEPTEDSET, redis.ZRangeBy{"0", "inf", 0, -1}).Result(); err != nil {
        panic(err)
    } else {
        //DebugLog("future demand by %d allocation records", len(resp))
        for _, item := range resp {
            // Val() ==> []string
            alloc := Allocation{}
            if record, err := redisClient.Get(item).Result(); err != nil {
                //错误消息为：redis: nil，有可能是已过期
                ErrorLog("redisclient get %s failed, reason: %s", item, err)
                continue
            } else {
                DebugLog("accepted record key name: %s, content: %s", item, record)
                if err := json.Unmarshal([]byte(record), &alloc); err != nil {
                    ErrorLog("get allocation key: %s failed, reason: %s", item, err)
                    continue
                }
                if t.Before(alloc.TEnd) && t.After(alloc.TStart) {
                    rRes.CpuUsed += int(alloc.Resources.GetCpuNum())
                    rRes.MemUsed += int(alloc.Resources.GetMemNum())
                    rRes.GpuUsed += int(alloc.Resources.GetGpuNum())
                }
            }

        }
    }
}
func estimateFrwDemand(t *time.Time, f int32) {
    // 获取概率的时隙列表
    //getFrwFreqAtTime(t, f) // 软件环境的使用不需要预测，硬件资源更加需要
}
func estimateCpuDemand(t *time.Time) {

}
func estimateGpuDemand(t *time.Time) {

}
func estimateMemDemand(t *time.Time) {

}

func getRunningSourcesAtTimeT(t *time.Time) *MachineUsageDetail {
    /*
     */
    cpuTotal, cpuUsed, gpuTotal, gpuUsed, memTotal, memUsed := 0, 0, 0, 0, 0, 0
    ids := getRunningClusters(t)
    for _, id := range ids {
        clusterUsage := getClusterUsageAtT(id, t)
        for _, detail := range clusterUsage.HUsage.HostDetail {
            cpuTotal += detail.CpuTotal
            cpuUsed += detail.CpuUsed
            gpuTotal += detail.GpuTotal
            gpuUsed += detail.GpuUsed
            memTotal += detail.MemTotal
            memUsed += detail.MemUsed
        }
    }
    return &MachineUsageDetail{
        MemTotal: memTotal,
        MemUsed:  memUsed,
        GpuTotal: gpuTotal,
        GpuUsed:  gpuUsed,
        CpuTotal: cpuTotal,
        CpuUsed:  cpuUsed,
    }

}
func getClusterUsageAtT(id ID, t *time.Time) *ClusterUsage {
    // TODO: 实时请求云，并缓存到本地的redis中
    return &ClusterUsage{
        ClusterId: id,
        HUsage: &HardwareResourceUsage{
            HostNumEmployeed: 0,
            HostDetail: []MachineUsageDetail{
                MachineUsageDetail{
                    MemTotal: 0,
                    MemUsed:  0,
                    GpuTotal: 0,
                    GpuUsed:  0,
                    CpuTotal: 0,
                    CpuUsed:  0,
                },
                MachineUsageDetail{
                    MemTotal: 0,
                    MemUsed:  0,
                    GpuTotal: 0,
                    GpuUsed:  0,
                    CpuTotal: 0,
                    CpuUsed:  0,
                },
            },
        },
    }
}
func getRunningClusters(t *time.Time) []ID {
    return []ID{"c4e9436b8530e4d739970e94943b18d9f", "c13208a667ad14e048ec10b3818a60470"}
}
func getAveCreatedHourSlot(sources []RussellExperiment) (int, []float32) {
    /*
       2018-04-02 数据
       2518
       [0.008737093 0.057585385 0.058776807 0.057188246 0.045274027 0.042891182 0.07625099 0.0802224
       0.064733915 0.035345513 0.037331216 0.06870532 0.08776807 0.07069103 0.072676726 0.059571087
       0.048848294 0.02104845 0.0051628277 0 0 0 0.0003971406 0.0007942812]
    */
    // 返回0~24的uint8，其中0表示00:00:00-00:59:59这个时间段内创建任务
    if sources == nil {
        sources = loadAllSourceData(originDatas)
    }
    slots := make([]int, 24)
    total := len(sources)
    for _, source := range sources {
        hours := source.DateCreated.Hour()
        slots[hours] += 1
        //InfoLog("%d'st source, id: %s created in %s", index, source.Id, source.DateCreated)
    }
    slotsFreq := make([]float32, 24)
    for i, v := range slots {
        slotsFreq[i] = float32(v) / float32(total)
    }
    return total, slotsFreq
}
func getAveStartedHourSlot(sources []RussellExperiment) (int, []float32) {
    /*
       2018-04-02 数据
       2518
       [0.008737093 0.056791104 0.059571087 0.056393962 0.04686259 0.042891182 0.07545671 0.07783955
       0.06513105 0.035345513 0.03931692 0.06830818 0.08935663 0.071088165 0.072676726 0.057982527
       0.04845115 0.021445591 0.0051628277 0 0 0 0.0003971406 0.0007942812]
    */
    if sources == nil {
        sources = loadAllSourceData(originDatas)
    }
    slots := make([]int, 24)
    total := len(sources)
    for _, source := range sources {
        hours := source.Started.Hour()
        slots[hours] += 1
        //InfoLog("%d'st source, id: %s created in %s", index, source.Id, source.DateCreated)
    }
    slotsFreq := make([]float32, 24)
    for i, v := range slots {
        slotsFreq[i] = float32(v) / float32(total)
    }
    return total, slotsFreq
}
func getFrwFreqAtTime(t *time.Time, f int32) (f32Freq float32) {
    /*
       某时段的所有框架的使用频率（不按启动时刻计算，而是按框架在线频率）
       冷启动问题：由RussellCloud那边的数据给出初始评分（ordered set），后面就直接取redis评分+本地作业记录的时间来计算频率了
    */
    // 时隙一小时
    slot_num := t.Hour() // 时隙序列号
    r := getFrwAtTimeTSlot(STARTED_TIMESTAGE, time.Hour, nil)
    return r[f][slot_num]
}
func getFrwAtTimeTSlot(ts ST_TIMESTAGETYPE, slot_width time.Duration, sources []RussellExperiment) (result map[int32][]float32) {
    // 在每个时隙内用户预定某框架的频率（以用户请求创建时间为准）
    // 返回 {框架A: [p0,...,p23],框架B: [p0,...,p23]}
    //result := make(map[int32]*[24]float32)
    if sources == nil {
        sources = loadAllSourceData(originDatas)
    }
    slots := []int{}
    // 计数
    for _, source := range sources {
        var createdT int
        var t time.Time
        switch slot_width {
        case time.Hour:
            if ts == RESERVED_TIMESTAGE {
                t = source.DateCreated
            } else if ts == STARTED_TIMESTAGE {
                t = source.Started
            }
            createdT = t.Hour()
        case time.Hour * 24 * 7:
            if ts == RESERVED_TIMESTAGE {
                t = source.DateCreated
            } else if ts == STARTED_TIMESTAGE {
                t = source.Started // sunday == 0, ...
            }
            createdT = int(t.Weekday()) // sunday == 0, ...
        }
        if frw, ok := FRAMEWORKKEYBYNAME[source.Environment]; ok {
            // 合法的框架名
            slots[createdT] += 1 //该时段的总数
            result[frw.Id][createdT] += 1
        } else {
            WarningLog("getFrwCreatedAtSlot error")
        }
    }
    // 频率
    for _, v := range result {
        for i := range v {
            v[i] /= float32(slots[i])
        }
    }

    return result
}
func statAllSourceData(sources []RussellExperiment) {
    /*
       从RussellCloud中挖掘所有信息，并存储到 内存 并缓存到 redis中。
    */
}
func loadAllSourceData(sources []RussellExperiment) []RussellExperiment {
    lines := 0
    errLines := 0
    if rows, err := database.Query(SELECTEXPUSEFULINFOSQL); err != nil {
        ErrorLog("query experiment failed, reason: %s", err)
        return nil
    } else {
        defer rows.Close()
        for rows.Next() {
            rusExp := RussellExperiment{}
            //const SELECT_EXP_USEFUL_INFO_SQL = "SELECT id,date_created,date_modified,name,description,module_id," +
            //    "owner_id,state,version,instance_type,started,ended,duration,project_id,mode,Command,environment," +
            //    "data_ids,output_size,enable_tensorboard FROM experiment;"
            err := rows.Scan(&rusExp.Id, &rusExp.DateCreated, &rusExp.DateModified,
                &rusExp.Name, &rusExp.Description, &rusExp.ModuleId,
                &rusExp.OwnerId, &rusExp.State, &rusExp.Version,
                &rusExp.InstanceType, &rusExp.Started, &rusExp.Ended,
                &rusExp.Duration, &rusExp.ProjectId, &rusExp.Mode,
                &rusExp.Command, &rusExp.Environment, &rusExp.DataIds,
                &rusExp.OutputSize, &rusExp.EnableTensorboard)
            if err != nil {
                errLines++
                //ErrorLog("scan row failed, reason: %s", err)
            } else {
                sources = append(sources, rusExp)
                //InfoLog("get row %v", rusExp)
            }
            lines++
        }
        InfoLog("扫描数据库： 共 %d 行, 其中 %d 行出错", lines, errLines)
        return sources
    }

}

var (
    // OjbK ^_^
    dataSourceName = "zuiwan:zuiwan2018@tcp(api.cannot.cc:3306)/Russell?charset=utf8&parseTime=true"
    database       *sql.DB
    originDatas    = make([]RussellExperiment, 0)
)

func init() {
    err := error(nil)
    database, err = sql.Open("mysql", dataSourceName)
    if err != nil {
        fmt.Println(err)
    }
    //originDatas = loadAllSourceData(originDatas)
}
