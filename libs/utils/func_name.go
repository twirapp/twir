package utils

import "runtime"

func GetFuncName() string {
    pc, _, _, ok := runtime.Caller(1)
    if !ok {
        return "unknown"
    }
    fn := runtime.FuncForPC(pc)
    if fn == nil {
        return "unknown"
    }
    return fn.Name()
}
