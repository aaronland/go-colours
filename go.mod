module github.com/aaronland/go-colours

go 1.24.2

// https://github.com/RobCherry/vibrant/pull/3	     
replace github.com/RobCherry/vibrant => github.com/sfomuseum/vibrant v0.0.0-20250430212339-abb21560aa26

require (
	github.com/RobCherry/vibrant v0.0.0-20160904011657-0680b8cf1c89
	github.com/aaronland/go-roster v1.0.0
	github.com/lucasb-eyer/go-colorful v1.2.0
	github.com/marekm4/color-extractor v1.2.1
	github.com/sfomuseum/go-flags v0.10.0
	github.com/sfomuseum/go-www-show v1.0.0
	golang.org/x/image v0.26.0
)

require (
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	golang.org/x/sys v0.1.0 // indirect
)
