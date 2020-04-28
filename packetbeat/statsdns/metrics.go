package statsdns


import "expvar"

var (
	counterReq = expvar.NewInt("counterReq")
	counterResp = expvar.NewInt("counterResp")
)