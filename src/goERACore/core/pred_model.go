package core

import "time"

//import "time"

// 与RussellCloud的experiment表结构保持一致
//+--------------------+--------------+------+-----+---------+-------+
//| Field              | Type         | Null | Key | Default | Extra |
//+--------------------+--------------+------+-----+---------+-------+
//| id                 | varchar(32)  | NO   | PRI | NULL    |       |
//| date_created       | datetime     | YES  |     | NULL    |       |
//| date_modified      | datetime     | YES  |     | NULL    |       |
//| is_deleted         | int(11)      | YES  |     | NULL    |       |
//| name               | varchar(128) | YES  | MUL | NULL    |       |
//| description        | varchar(64)  | YES  |     | NULL    |       |
//| module_id          | varchar(32)  | YES  | MUL | NULL    |       |
//| owner_id           | varchar(32)  | YES  | MUL | NULL    |       |
//| permission         | int(11)      | YES  |     | NULL    |       |
//| log_id             | varchar(32)  | YES  | UNI | NULL    |       |
//| state              | varchar(32)  | YES  |     | NULL    |       |
//| data_id            | varchar(32)  | YES  |     | NULL    |       |
//| family_id          | varchar(64)  | YES  |     | NULL    |       |
//| version            | double       | YES  |     | NULL    |       |
//| predecessor        | varchar(64)  | YES  |     | NULL    |       |
//| instance_type      | varchar(64)  | YES  |     | NULL    |       |
//| started            | datetime     | YES  |     | NULL    |       |
//| ended              | datetime     | YES  |     | NULL    |       |
//| duration           | int(11)      | YES  |     | NULL    |       |
//| canvas             | varchar(64)  | YES  |     | NULL    |       |
//| task_instance_ids  | varchar(256) | YES  |     | NULL    |       |
//| project_id         | varchar(32)  | YES  |     | NULL    |       |
//| mode               | varchar(32)  | YES  |     | NULL    |       |
//| inputs             | varchar(256) | YES  |     | NULL    |       |
//| outputs            | varchar(256) | YES  |     | NULL    |       |
//| Command            | varchar(256) | YES  |     | NULL    |       |
//| environment        | varchar(32)  | YES  |     | NULL    |       |
//| data_ids           | varchar(512) | YES  |     | NULL    |       |
//| full_command       | varchar(512) | YES  |     | NULL    |       |
//| output_size        | int(11)      | YES  |     | NULL    |       |
//| enable_tensorboard | tinyint(1)   | YES  |     | NULL    |       |
//| system_stopped     | varchar(32)  | YES  |     | NULL    |       |
//+--------------------+--------------+------+-----+---------+-------+
type RussellExperiment struct {
	Id           string    // Y
	DateCreated  time.Time //Y
	DateModified time.Time //Y
	//DateCreated  []uint8 //Y
	//DateModified []uint8 //Y
	IsDeleted    bool      //N
	Name         string    //Y
	Description  string    //Y
	ModuleId     string    //Y
	OwnerId      string    //Y
	Permission   int       //N
	LogId        string    //N
	State        string    //Y
	DataId       string    //N
	FamilyId     string    //N
	Version      int       //Y
	Predecessor  string    //N
	InstanceType string    //Y
	Started      time.Time //Y
	Ended        time.Time //Y
	//Started  []uint8 //Y
	//Ended    []uint8 //Y
	Duration int //Y
	//Duration time.Duration
	Canvas            string //N
	TaskInstanceIds   string //N
	ProjectId         string //Y
	Mode              string //Y
	Inputs            string //N
	Outputs           string //N
	Command           string //Y
	Environment       string //Y
	DataIds           string //Y
	FullCommand       string //N
	OutputSize        string //Y
	EnableTensorboard bool   //Y
	SystemStopped     string //N
}

const SELECTEXPUSEFULINFOSQL = "SELECT id,date_created,date_modified,name,description,module_id," +
	"owner_id,state,version,instance_type,started,ended,duration,project_id,mode,Command,environment," +
	"data_ids,output_size,enable_tensorboard FROM experiment "

	//SELECT * INTO OUTFILE '/var/lib/mysql-files/table.csv'
	// FIELDS TERMINATED BY ',' OPTIONALLY ENCLOSED BY '"'     LINES TERMINATED BY '\n' FROM experiment;
