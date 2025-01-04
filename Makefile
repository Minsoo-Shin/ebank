.PHONY: all clean proto

# 변수 정의
PROTO_DIR := api/v1
FILES = data/account.json data/transaction.json data/user.json

all: proto

# Protocol Buffers 파일 컴파일
proto:
	@protoc --proto_path=. \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	    --grpc-gateway_out . --grpc-gateway_opt paths=source_relative \
	    --openapiv2_out . --openapiv2_opt use_go_templates=true \
		$(PROTO_DIR)/*.proto

# 생성된 파일 삭제
clean:
	@find $(PROTO_DIR) \( -name "*.pb.go" -o -name "*_grpc.pb.go" \) -type f -delete

# mock 파일 생성
mock:
	@mockery --dir=internal/service --all --outpkg=mocks --with-expecter=true --recursive=true
	@mockery --dir=internal/repository --all --outpkg=mocks --with-expecter=true --recursive=true
	@mockery --dir=pkg --all --outpkg=mocks --with-expecter=true --recursive=true


# .PHONY 타겟 지정
.PHONY: file check file-clean

# 기본 타겟
file: $(FILES)

# 파일 생성 규칙
%.json:
	@echo "Creating $@..."
	@echo "" > $@
	@echo "File $@ created successfully."

# 모든 파일 생성 확인
check:
	@echo "Checking files..."
	@for file in $(FILES); do \
		if [ -f $$file ]; then \
			echo "$$file exists."; \
		else \
			echo "$$file does not exist."; \
		fi \
	done

# 파일 삭제 규칙
file-clean:
	@echo "Cleaning up..."
	@rm -f $(FILES)
	@echo "Cleanup complete."


# 애플리케이션 실행
run:
	@go run cmd/app/main.go

# 테스트 실행
test:
	@go test -v -cover ./...
