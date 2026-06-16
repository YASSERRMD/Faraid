package api

import "expvar"

// Request counters exported via expvar at /debug/vars.
var (
	metricSolveTotal    = expvar.NewInt("faraid_solve_total")
	metricSolveErrors   = expvar.NewInt("faraid_solve_errors")
	metricCompareTotal  = expvar.NewInt("faraid_compare_total")
	metricExplainTotal  = expvar.NewInt("faraid_explain_total")
	metricParseTotal    = expvar.NewInt("faraid_parse_total")
)
