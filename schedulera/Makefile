build:
	CGO_ENABLED=0 GOOS=linux go build  -tags netgo 
	mv schedulera karmada-scheduler
	docker build -t ghcr.io/lank8s/karmada-schedulera:latest .
	docker build -t ghcr.io/lank8s/karmada-schedulera:test .
	rm -f karmada-scheduler
	docker push ghcr.io/lank8s/karmada-schedulera:latest
	docker push ghcr.io/lank8s/karmada-schedulera:test
