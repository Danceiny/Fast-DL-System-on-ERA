package core

/*
    Util functions followed
*/
func findMaxT(arr []uint32) uint64 {
    max := uint32(0)
    maxI := uint64(0)
    for i, v := range arr {
        if v > max {
            maxI = uint64(i)
        }
    }
    return maxI
}
func findMinT(arr []uint32) uint64 {
    // TODO: 优化
    min := arr[0]
    minI := 0
    for i, v := range arr {
        if v < min {
            minI = i
        }
    }
    return uint64(minI)
}

func convert2Float64LessThanOne(n uint32) float64 {
    // 1331 ==> 0.1331
    // Algorithm Reference:
    // https://stackoverflow.com/questions/701322/how-can-you-get-the-first-digit-in-an-int-c/701621#701621
    i := float64(n)
    if i >= 100000000 {
        i /= 100000000
    }
    if i >= 10000 {
        i /= 10000
    }
    if i >= 100 {
        i /= 100
    }
    if i >= 10 {
        i /= 10
    }
    return i / 10
}
