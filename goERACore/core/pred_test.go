package core

import "testing"

func TestGetAveCreatedHourSlot(t *testing.T) {
    t.Skip()
    processed, res := getAveCreatedHourSlot(nil)
    t.Log(processed, res)
}
func procSlice(s []int) []int{
    s = append(s, 10)
    return s
}

func TestSlice(t *testing.T) {
    t.Skip()
    sli := make([]int, 0)
    t.Log(sli, len(sli), cap(sli))
    sli = procSlice(sli)
    t.Log(sli, len(sli), cap(sli))
}
func TestGetAveStartedHourSlot(t *testing.T){
    processed, res := getAveStartedHourSlot(nil)
    t.Log(processed, res)
}
