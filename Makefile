runDev:
	@CompileDaemon -build="go build -o findyourpet-backend.exe" -command="./cmd/findyourpet-backend.exe" -directory=cmd -directory=./ -exclude=Makefile -exclude=.exe -exclude=.exe~ -exclude=.git exclude-dir=".trunk"
# mocks:
# 	@mockgen -source=models/postModel.go -destination=mocks/postModel.go -imports=models/userModel.go -package=mocks

build-go:
	CGO_ENABLED=0 && ENV="PROD" && go build -a -installsuffix cgo -o findyourpet-backend cmd/main.go cmd/config.go

mocks:
	mockgen -source=internal/store/users/store.go Store >test/users/mock_store.go

docker-build:
	docker build --progress=plain --no-cache -t findyourpet-backend .

cloud-build:
	docker buildx build --builder cloud-kach022-findyourpet --no-cache -t findyourpet-backend .

generate-mock:
	mockgen -destination=mocks/mock_s3api.go -package=mocks github.com/anastasiagoncharova/findyourpet/findyourpet-backend/pkg/awsS3 S3API

run-tests:
