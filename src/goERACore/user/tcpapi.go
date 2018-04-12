package user

import (
    "net"
    "bytes"
    "time"
    "io"
    "encoding/json"
    "encoding/binary"
    . "goERACore/core"
)

const (
    SERVERNETWORK = "tcp"
    SERVERADDRESS = "127.0.0.1:5555"
    //DELIMITER     = '\t'
    // 连续出现 `{REQ_EOF_HEADER}{REQ_EOF_CONFIRM}`这两个字符，视为字节流的结束标志
    REQEOFHEADER  = '\t'
    REQEOFCONFIRM = '\n'
)

var (
    // 处理请求数，用于测量性能
    REQCOUNTCHAN = make(chan uint64, 1)
)

func init() {
    REQCOUNTCHAN <- 1
}

//func PrintLog(format string, a ...interface{}) (n int, err error) {
//    return fmt.Printf(format, a...)
//}
func TCPInterface() {
    var listener net.Listener
    listener, err := net.Listen(SERVERNETWORK, SERVERADDRESS)
    if err != nil {
        PrintLog("debug", "Listen Error: %s\n", err)
    } else {
        defer listener.Close()
        PrintLog("debug", "Got listener for the server. (local address: %s)\n", listener.Addr())

        for {
            conn, err := listener.Accept() // blocked until new connection arrives
            if err != nil {
                PrintLog("debug", "Accept Error: %s\n", err)
            }
            PrintLog("debug", "Established a connection with a client application. (remove address: %s)\n",
                conn.RemoteAddr())
            go handleConn(conn)
        }
    }

}
func handleConn(conn net.Conn) {

    if cnt, ok := <-REQCOUNTCHAN; ok {
        //fetch from channel is [%!s(bool=true) %!s(uint64=4)], The %!d(MISSING)'st request
        PrintLog("info", "【请求处理计数】fetch from channel, will process the %d'st request", cnt)
        cnt++
        REQCOUNTCHAN <- cnt
    } else {
        PrintLog("info", "fetch from channel is failed: %b", ok)
    }

    //defer wg.Done()
    defer conn.Close()
    for {
        conn.SetReadDeadline(time.Now().Add(10 * time.Second))
        bytesReq, err := readBytesUntilSpecifiedEOF(conn, REQEOFHEADER, REQEOFCONFIRM)
        if err != nil {
            if err == io.EOF {
                PrintLog("debug", "The connection is closed by another side.(Server)\n")
            } else {
                PrintLog("debug", "Read Error: %s (Server)\n", err)
            }
            break
        }
        PrintLog("debug", "Received request: %s (Server)\n", bytesReq)

        // part 2
        var jobRequest JobRequest
        if err := json.Unmarshal(bytesReq, &jobRequest); err != nil {
            ErrorLog("unmarshal request failed, reason: %s", err)
            continue
        }
        resp := BasicEconScheduling(&jobRequest)
        respBytes, err := json.Marshal(resp)
        if err != nil {
            PrintLog("debug", "marshal struct to json string, error: %s", err)
        } else {
            //PrintLog("debug", "respBytes is %v", respBytes)
        }
        var lenBytes = make([]byte, 4)
        binary.LittleEndian.PutUint32(lenBytes, uint32(len(respBytes)))
        sendBytes := append(lenBytes, respBytes...)
        n, err := writeBytesWithSpecifiedEOF(conn, sendBytes, REQEOFHEADER, REQEOFCONFIRM)
        if err != nil {
            PrintLog("debug", "Write error: %s, written %d bytes however (Server)\n", err, n)
        } else {

            PrintLog("debug", "Sent response (written %d bytes) (Server)\n", n)
        }
    }
}
func readBytesUntilSpecifiedEOF(conn net.Conn, eofH byte, eofC byte) ([]byte, error) {
    readBytes := make([]byte, 1)
    var buffer bytes.Buffer
    bEofStarted := false
    var eofBuf byte
    for {
        _, err := conn.Read(readBytes)
        if err != nil {
            return nil, err
        }
        readByte := readBytes[0]
        if bEofStarted {
            if readByte == eofC {
                // 确认终结符
                break
            } else {
                // 并不是真正的终结符
                bEofStarted = false
                buffer.WriteByte(eofBuf)
            }
        } else if readByte == eofH {
            bEofStarted = true
            eofBuf = readByte // 缓存一个字节，等待终结符确认
            continue
        }

        buffer.WriteByte(readByte)
    }
    return buffer.Bytes(), nil
}
func writeBytesWithSpecifiedEOF(conn net.Conn, content []byte, eofH byte, eofC byte) (int, error) {
    var buffer bytes.Buffer
    content = append(content, eofH, eofC)
    for _, b := range content {
        buffer.WriteByte(b)
    }
    return conn.Write(buffer.Bytes())
}
