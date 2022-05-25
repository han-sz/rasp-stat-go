debug:
	go build -o build/rasp-stat_debug ./rasp-stat

prod-arm:
	GOARCH=arm
	export GOARCH=arm GOOS=linux; go build -ldflags="-s -w" -o build/rasp-stat_$${GOARCH}-$${GOOS} ./rasp-stat
	export GOARCH=arm64 GOOS=linux; go build -ldflags="-s -w" -o build/rasp-stat_$${GOARCH}-$${GOOS} ./rasp-stat

prod-x86:
	# export GOARCH=386 GOOS=linux; go build -ldflags="-s -w" -o build/rasp-stat_$${GOARCH}-$${GOOS} ./rasp-stat
	export GOARCH=amd64 GOOS=darwin; go build -ldflags="-s -w" -o build/rasp-stat_$${GOARCH}-$${GOOS} ./rasp-stat

prod: prod-arm prod-x86