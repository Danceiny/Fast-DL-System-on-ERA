package loadgen

import (
    "time"
    "fmt"
    "errors"
    . "goSimulator/loadgen/lib"
    "math"
    "bytes"
)

type implGenerator struct {
    timeoutNs   time.Duration //响应超时时间，单位：纳秒
    lps         uint32        // 每秒载荷发送量
    durationNs  time.Duration //负载持续时间，单位：纳秒
    resultCh    chan *CallResult
    concurrency uint32      //并发量
    tickets     GoTickets   //goroutine票池
    stopSign    chan byte   //停止信号的传递通道
    cancelSign  byte        // 取消发送后续结果的信号
    endSign     chan uint64 //
    status      GenStatus   //载荷发生器的状态
    caller      Caller      //调用器
    callCount   uint64
}

func (gen *implGenerator) Start() {
    logger.Infoln("Starting load generator...")

    // 设定节流阀
    var throttle <-chan time.Time
    if gen.lps > 0 {
        interval := time.Duration(1e9 / gen.lps)
        logger.Infof("Setting throttle (%v)...", interval)
        throttle = time.Tick(interval)
    }

    // 初始化停止信号
    go func() {
        time.AfterFunc(gen.durationNs, func() {
            logger.Infof("Stopping load generator...")
            gen.stopSign <- 0
        })
    }()

    // 初始化完结信号通道
    gen.endSign = make(chan uint64, 2)

    // 初始化调用执行计数
    gen.callCount = 0

    // 设置已启动状态
    gen.status = STATUS_STARTED

    go func() {
        // 生成载荷
        logger.Infoln("Generating loads...")
        gen.genLoad(throttle)

        // 接收调用执行计数
        callCount := <-gen.endSign
        gen.status = STATUS_STOPPED
        logger.Infof("Stopped. (callCount=%d)\n", callCount)
    }()
}
func (gen *implGenerator) Stop() (uint64, bool) {
    if gen.stopSign == nil {
        return 0, false
    }
    if gen.status != STATUS_STARTED {
        return 0, false
    }
    gen.status = STATUS_STOPPED
    gen.stopSign <- 1
    callCount := <-gen.endSign
    return callCount, true
}
func (self *implGenerator) Status() GenStatus {
    return self.status
}

func NewGenerator(
    caller Caller,
    timeoutNs time.Duration,
    lps uint32,
    durationNs time.Duration,
    resultCh chan *CallResult) (Generator, error) {
    var errMsg string
    if caller == nil {
        errMsg = fmt.Sprintf("Invalid caller!")
    }
    if timeoutNs == 0 {
        errMsg = fmt.Sprintf("Invalid timeoutNs!")
    }
    if lps == 0 {
        errMsg = fmt.Sprintf("Invalid lps!")
    }
    if durationNs == 0 {
        errMsg = fmt.Sprintf("Invalid durationNs!")

    }
    if resultCh == nil {
        errMsg = fmt.Sprintf("Invalid !")

    }
    if errMsg != "" {
        return nil, errors.New(errMsg)
    }

    gen := &implGenerator{
        caller:     caller,
        timeoutNs:  timeoutNs,
        lps:        lps,
        durationNs: durationNs,
        stopSign:   make(chan byte, 1),
        cancelSign: 0,
        status:     STATUS_ORIGINAL,
        resultCh:   resultCh,
    }
    logger.Infof("Passed. (timeoutNs=%v, lps=%d, durationNs=%v)",
        timeoutNs, lps, durationNs)
    if err := gen.init(); err != nil {
        return nil, err
    }

    return gen, nil
}
func (gen *implGenerator) init() error {
    //并发量 约等于 单个载荷的响应超时时间/载荷的发送间隔时间
    //即 concurrency = timeoutNs / (1e9 / lps ) + 1
    //计算并发量的最大意义在于：为约束被并发运行的goroutine的数量提供支撑，即依据此数值确定载荷发生器的tickets字段所代表的goroutine票池的容量

    var total64 int64 = int64(gen.timeoutNs)
    if total64 > math.MaxInt32 {
        total64 = math.MaxInt32
    }
    gen.concurrency = uint32(total64)
    tickets, err := NewGoTickets(gen.concurrency)
    if err != nil {
        return err
    }
    gen.tickets = tickets
    logger.Infof("Initialized. (concurrency: %d)", gen.concurrency)
    return nil
}

func (gen *implGenerator) genLoad(throttle <-chan time.Time) {
    /*
    throttle: 节流阀
    stopSign: 传递停止信号
     */
    callCount := uint64(0)
Loop:
    for ; ; callCount++ {
        select {
        case <-gen.stopSign:
            gen.handleStopSign(callCount)
            break Loop
        default:

        }
        gen.asyncCall()
        if gen.lps > 0 {
            select {
            case <-throttle:
            case <-gen.stopSign:
                gen.handleStopSign(callCount)
                break Loop
            }
        }
    }
}

func (gen *implGenerator) handleStopSign(callCount uint64) {
    gen.cancelSign = 1
    logger.Infof("Closing result channel...")
    close(gen.resultCh)
    // 两种载荷发生器停止方式的冲突，所以需要确认
    gen.endSign <- callCount
    gen.endSign <- callCount
}
func (gen *implGenerator) asyncCall() {
    gen.tickets.Take()
    go func() {
        defer func() {
            // 处理异常
            if p := recover(); p != nil {
                err, ok := interface{}(p).(error)
                var buff bytes.Buffer
                buff.WriteString("Async call panic!")
                if ok {
                    buff.WriteString(fmt.Sprintf("error: %s", err.Error()))
                } else {
                    buff.WriteString(fmt.Sprintf("clue: %v", p))
                }
                buff.WriteString(")")
                errMsg := buff.String()
                logger.Fatalln(errMsg)
                result := &CallResult{
                    Id:   -1,
                    Code: RESULT_CODE_ERROR_CALL,
                    Msg:  errMsg,
                }
                gen.sendResult(result)
            }
        }()

        rawReq := gen.caller.BuildReq()
        var bTimeout bool
        timer := time.AfterFunc(gen.timeoutNs, func() {
            bTimeout = true
            result := &CallResult{
                Id:   rawReq.Id,
                Req:  rawReq,
                Code: RESULT_CODE_WARNING_CALL_TIMEOUT,
                Msg:  fmt.Sprintf("Timeout!(expected: < %v)", gen.timeoutNs),
            }
            gen.sendResult(result)
        })
        rawResp := gen.interact(&rawReq)
        if !bTimeout {
            timer.Stop()
            var result *CallResult
            if rawResp.Err != nil {
                result = &CallResult{
                    Id:     rawResp.Id,
                    Req:    rawReq,
                    Code:   RESULT_CODE_ERROR_CALL,
                    Msg:    rawResp.Err.Error(),
                    Elapse: rawResp.Elapse,
                }
            } else {
                result = gen.caller.CheckResp(rawReq, *rawResp)
                result.Elapse = rawResp.Elapse
            }
            gen.sendResult(result)
        }
        gen.tickets.Return()
    }()
}
func (self *implGenerator) interact(rawReq *RawReq) *RawResp {
    if rawReq == nil {
        return &RawResp{
            Id:     -1,
            Resp:   nil,
            Err:    errors.New("Invalid raw request."),
            Elapse: 0,
        }
    }
    start := time.Now().Nanosecond()
    resp, err := self.caller.Call(rawReq.Req, self.timeoutNs)
    end := time.Now().Nanosecond()
    elapsedTime := time.Duration(end - start)
    var rawResp RawResp
    if err != nil {
        errMsg := fmt.Sprintf("Sync call error: %s.", err)
        rawResp = RawResp{
            Id:     rawReq.Id,
            Resp:   nil,
            Err:    errors.New(errMsg),
            Elapse: elapsedTime,
        }
    } else {
        rawResp = RawResp{
            Id:     rawReq.Id,
            Resp:   resp,
            Err:    nil,
            Elapse: elapsedTime,
        }
    }
    return &rawResp
}

func (self *implGenerator) sendResult(result *CallResult) bool {
    if self.status == STATUS_STARTED && self.cancelSign == 0 {
        self.resultCh <- result
        return true
    }
    logger.Warn("Ignore result: %s.",
        fmt.Sprintf("Id=%d, Code=%d, Msg=%s, Elapse=%v",
            result.Id, result.Code, result.Msg, result.Elapse))
    return false
}
