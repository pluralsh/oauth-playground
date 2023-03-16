
run:
	go run server.go -zap-devel -zap-log-level 2

generate:
	go run github.com/99designs/gqlgen generate

generate-db:
	go generate ./ent

setup-tilt-cluster:
	ctlptl create cluster kind --registry=ctlptl-registry
