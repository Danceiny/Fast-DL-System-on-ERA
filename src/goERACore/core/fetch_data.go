package core

import "time"

/*
获取数据（新数据）
*/

func getJob(id ID) *Job {
	// TODO: 数据库查询
	return &Job{
		Id:        "",
		Duration:  0,
		TCreated:  time.Time{},
		TModified: time.Time{},
		ProjectId: "",
		CodeId:    "",
		DataIds:   nil,
		Env: &Environment{
			DlFrName: "caffe2",
			Os:       "",
			WithGpu:  false,
		},
		EntryCmd:     "",
		StartCmd:     "",
		BTensorboard: false,
		BJupyter:     false,
		TStarted:     time.Time{},
		TEnded:       time.Time{},
		Uid:          "",
		Gid:          "",
		Perm:         0,
		Doc:          "",
	}
}
func getAllJobs(r interface{}) []Job {
	/*
	   根据范围r来查询，r相当于过滤器，可以为nil
	*/
	return []Job{{}}
}
