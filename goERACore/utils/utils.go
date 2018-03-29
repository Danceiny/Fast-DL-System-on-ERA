package utils

func max(first int, args... int) int {
    for _ , v := range args{
        if first < v {
            first = v
        }
    }
    return first
}
