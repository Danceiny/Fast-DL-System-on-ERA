package core

import "fmt"

func PrintLog(level, msg string, a ...interface{}) {
    format := fmt.Sprintf("[%s]: %s\n", level, msg)
    fmt.Printf(format, a)
}

func DebugLog(msg string, a ...interface{}) {
    PrintLog("debug", msg, a)
}
func InfoLog(msg string, a ...interface{}) {
    PrintLog("info", msg, a)
}
func ErrorLog(msg string, a ...interface{}) {
    PrintLog("error", msg, a)
}
func CriticalLog(msg string, a ...interface{}) {
    PrintLog("critical", msg, a)
}
func WarningLog(msg string, a ...interface{}) {
    PrintLog("warning", msg, a)
}
