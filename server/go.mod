module github.com/mizanmahi/aiusage/server

go 1.25.4

require (
	github.com/go-chi/chi/v5 v5.3.0
	github.com/mizanmahi/aiusage/types v0.0.0
)

require github.com/lib/pq v1.12.3

replace github.com/mizanmahi/aiusage/types => ../types
