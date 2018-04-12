package core

import (
	_ "github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

func TestExpTable(t *testing.T) {
	t.Skip()
	return
	//defer database.Close()
	//getTaskSQL := "show tables;"
	//rows, err := database.Query(getTaskSQL)
	//if (err != nil) {
	//    fmt.Println(err)
	//} else {
	//    //fmt.Println(rows)
	//}
	//defer rows.Close()
	//for rows.Next() {
	//    tablename := ""
	//    err := rows.Scan(&tablename)
	//    if err != nil {
	//        fmt.Println(err)
	//    } else {
	//        fmt.Println(tablename)
	//    }
	//}

	//var tablename string;
	//stmt, err := database.Prepare("SELECT * from experiment where username=?")
	//err = stmt.QueryRow("Danceiny").Scan(&tablename)
	//if err != nil {
	//    fmt.Println(err)
	//} else {
	//    fmt.Println(tablename)
	//}

	if rows, err := database.Query(SELECTEXPUSEFULINFOSQL + " LIMIT 10"); err != nil {
		t.Error("query experiment failed, reason: %s", err)
		return
	} else {
		defer rows.Close()
		for rows.Next() {
			rusExp := RussellExperiment{}
			//const SELECT_EXP_USEFUL_INFO_SQL = "SELECT id,date_created,date_modified,name,description,module_id," +
			//    "owner_id,state,version,instance_type,started,ended,duration,project_id,mode,Command,environment," +
			//    "data_ids,output_size,enable_tensorboard FROM experiment;"
			var tm = time.Time{}
			err := rows.Scan(&rusExp.Id, &tm, &rusExp.DateModified,
				&rusExp.Name, &rusExp.Description, &rusExp.ModuleId,
				&rusExp.OwnerId, &rusExp.State, &rusExp.Version,
				&rusExp.InstanceType, &rusExp.Started, &rusExp.Ended,
				&rusExp.Duration, &rusExp.ProjectId, &rusExp.Mode,
				&rusExp.Command, &rusExp.Environment, &rusExp.DataIds,
				&rusExp.OutputSize, &rusExp.EnableTensorboard)
			//t.Logf("%s %v", tm, tm.Unix())
			if err != nil {
				ErrorLog("scan row failed, reason: %s", err)
			} else {
				InfoLog("get row %v", rusExp)
			}
		}
	}
}
