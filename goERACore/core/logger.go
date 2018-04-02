package core

import "fmt"

//import "fmt"

func PrintLog(level, msg string, a ...interface{}) {
    //fmt.Printf("%s: %s ",level, msg, fmt...)
    format := fmt.Sprintf("[%s]: %s\n", level, msg)
    fmt.Printf(format, a)
}
