build:
	CGO_ENABLED=0 GOOS=linux go build  -tags netgo 
	mv schedulerb karmada-schedulerb
	docker build -t ghcr.io/lank8s/karmada-schedulerb:latest .
	docker build -t ghcr.io/lank8s/karmada-schedulerb:test .
	rm -f karmada-schedulerb
	docker push ghcr.io/lank8s/karmada-schedulerb:latest
	docker push ghcr.io/lank8s/karmada-schedulerb:test