package user

import (
    "net"
    "time"
    "fmt"
    "math/rand"
    "io"
    "sync"
    "testing"
    "encoding/json"
)

var wg sync.WaitGroup

func userClient(id ID) {
    defer wg.Done()
    conn, err := net.DialTimeout(SERVER_NETWORK, SERVER_ADDRESS, 2*time.Second)
    if err != nil {
        printLog("Dial Error: %s (Client[%d])\n", err, id)
        return
    }
    defer conn.Close()
    printLog("Connected to server. (remote address： %s, local address: %s)(Client【%d】)\n", conn.RemoteAddr(), conn.LocalAddr(), id)
    time.Sleep(200 * time.Millisecond)
    requestNumber := 5
    conn.SetDeadline(time.Now().Add(5 * time.Millisecond))
    for i := 0; i < requestNumber; i++ {
        i32Req := rand.Int31()
        n, err := write(conn, fmt.Sprintf("%d", i32Req))
        if err != nil {
            printLog("Write Error: %s(Client[%d])\n", err, id)
            continue
        }
        printLog("Send request (written %d bytes): %d (Client[%d])\n", n, i32Req, id)
    }
    for j := 0; j < requestNumber; j++ {
        strResp, err := read(conn)
        if err != nil {
            if err == io.EOF {
                printLog("The connection is closed by another side.(Client[%d])\n", id)

            } else {
                printLog("Read error: %s (client[%d])\n", err, id)

            }
            break
        }
        printLog("Received response: %s (client[%d])\n", strResp, id)
    }
}

func TestJsonStruct(t *testing.T) {
    job := Job{
        Id:        0,
        Duration:  0,
        TCreated:  time.Time{},
        TModified: time.Time{},
        ProjectId: 0,
        CodeId:    0,
        DataId:    1,
        Env: &Environment{
            DlFr: &DL_Framework{
                Name:      "name",
                NameSpace: "",
                ImgTag:    "",
                PyVer:     0,
            },
            OsVer: "",
        },
        EntryCmd:     "",
        StartCmd:     "",
        BTensorboard: false,
        BJupyter:     false,
        TStarted:     time.Time{},
        TEnded:       time.Time{},
        UId:          0,
        GId:          0,
        Perm:         0,
        Doc:          "",
    }
    s, err := json.Marshal(job)
    //println(string(s))
    t.Log(string(s), err)
    var newjob Job
    err = json.Unmarshal(s, &newjob)
    //_, err = simplejson.NewJson(s)
    t.Log(newjob.Env.DlFr.Name, err)
}

//func TestClient(t *testing.T) {
//    wg.Add(2)
//    go serverGo()
//    time.Sleep(500 * time.Millisecond)
//    go userClient(1)
//    wg.Wait()
//    t.Log("Success")
//}
