package core

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "time"
)

func estimateDemand(t *time.Time, p int64) {
    // t时刻
    // p价格，为负时视为省略

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
func getFrwFreqAtTime(t *time.Time, f int32, tw time.Duration, sources []RussellExperiment) {
    /*
    冷启动问题：由RussellCloud那边的数据给出初始评分（ordered set），后面就直接取redis评分+本地作业记录的时间来计算频率了
     */

}
func getFrwCreatedHourSlot(sources []RussellExperiment) (result map[int32]*[24]float32) {
    // 在每个小时用户预定某框架的频率（以用户请求创建时间为准）
    // 返回 {框架A: [p0,...,p23],框架B: [p0,...,p23]}
    //result := make(map[int32]*[24]float32)
    if sources == nil {
        sources = loadAllSourceData(originDatas)
    }
    slots := make([]int, 24)
    // 计数
    for _, source := range sources {
        createdHour := source.DateCreated.Hour()
        if frw, ok := FRAMEWORKKEYBYNAME[source.Environment]; ok {
            // 合法的框架名
            slots[createdHour] += 1 //该时段的总数
            result[frw.Id][createdHour] += 1
        } else {
            WarningLog("getFrwCreatedHourSlot error")
        }
    }
    // 频率
    for _, v := range result {
        for hour := range v {
            v[hour] /= float32(slots[hour])
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
