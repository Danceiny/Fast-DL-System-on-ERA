package user

import (
    "net"
    "bytes"
    "time"
    "io"
    "encoding/json"
    "encoding/binary"
    . "goERACore/core"
    //"goERACore/core"
)

const (
    SERVER_NETWORK = "tcp"
    SERVER_ADDRESS = "127.0.0.1:5555"
    DELIMITER      = '\t'
    // 连续出现 `{REQ_EOF_HEADER}{REQ_EOF_CONFIRM}`这两个字符，视为字节流的结束标志
    REQ_EOF_HEADER  = '\t'
    REQ_EOF_CONFIRM = '\n'
)

//func PrintLog(format string, a ...interface{}) (n int, err error) {
//    return fmt.Printf(format, a...)
//}
func TCP_Interface() {
    var listener net.Listener
    listener, err := net.Listen(SERVER_NETWORK, SERVER_ADDRESS)
    if err != nil {
        PrintLog("debug", "Listen Error: %s\n", err)
        return
    }
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
func handleConn(conn net.Conn) {
    //defer wg.Done()
    defer conn.Close()
    for {
        conn.SetReadDeadline(time.Now().Add(10 * time.Second))
        bytesReq, err := readBytesUntilSpecifiedEOF(conn, REQ_EOF_HEADER, REQ_EOF_CONFIRM)
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
        json.Unmarshal(bytesReq, &jobRequest)
        //PrintLog("debug", "resources: %v", jobRequest.Resources)
        PrintLog("debug", "resources.Gpu: %d", jobRequest.Resources.Gpu)
        resp := BasicEconScheduling(&jobRequest)
        PrintLog("debug", "resp struct is: %v", resp)
        respBytes, err := json.Marshal(resp)
        if err != nil {
            PrintLog("debug", "marshal struct to json string, error: %s", err)
        } else {
            PrintLog("debug", "respBytes is %v", respBytes)
        }
        var lenBytes = make([]byte, 4)
        PrintLog("debug", "respBytes length: %d", uint32(len(respBytes)))
        binary.LittleEndian.PutUint32(lenBytes, uint32(len(respBytes)))
        sendBytes := append(lenBytes, respBytes...)
        PrintLog("debug", "the full bytes to send is %v, %s\n, length:%d, msg:%s", sendBytes, string(sendBytes),
            sendBytes[:4], string(sendBytes[4:]))
        n, err := writeBytesWithSpecifiedEOF(conn, sendBytes, REQ_EOF_HEADER, REQ_EOF_CONFIRM)
        if err != nil {
            PrintLog("debug", "Write error: %s (Server)\n", err)
        }
        PrintLog("debug", "Sent response (written %d bytes): error: %s (Server)\n", n, err)
    }
}
func readBytesUntilSpecifiedEOF(conn net.Conn, eof_h byte, eof_c byte) ([]byte, error) {
    readBytes := make([]byte, 1)
    var buffer bytes.Buffer
    b_eof_started := false
    var eofBuf byte
    for {
        _, err := conn.Read(readBytes)
        if err != nil {
            return nil, err
        }
        readByte := readBytes[0]
        if b_eof_started {
            if readByte == eof_c {
                // 确认终结符
                break
            } else {
                // 并不是真正的终结符
                b_eof_started = false
                buffer.WriteByte(eofBuf)
            }
        } else if readByte == eof_h {
            b_eof_started = true
            eofBuf = readByte // 缓存一个字节，等待终结符确认
            continue
        }

        buffer.WriteByte(readByte)
    }
    return buffer.Bytes(), nil
}
func writeBytesWithSpecifiedEOF(conn net.Conn, content []byte, eof_h byte, eof_c byte) (int, error) {
    var buffer bytes.Buffer
    content = append(content, eof_h, eof_c)
    for _, b := range content {
        buffer.WriteByte(b)
    }
    return conn.Write(buffer.Bytes())
}
