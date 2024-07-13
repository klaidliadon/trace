module github.com/klaidliadon/trace

go 1.22.5

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-chi/httplog/v2 v2.0.11
	github.com/webrpc/webrpc v0.19.3
)

replace github.com/go-chi/httplog/v2 => github.com/klaidliadon/httplog/v2 v2.0.0-20240714061042-420abc5265d5
