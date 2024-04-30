all: compile upload clean

compile: check-go build zip
upload:compile to-s3 to-lambda update-ses-template
call: invoke-lambda
clean: clean-all

# name of s3 bucket
BUCKET_NAME ?= shoestring-lambda-bucket
# 
FUNC_NAME ?= shoestringSendEmail


# defn for default path and args for build command
GOOS ?= linux
GOARCH ?= amd64
CGO_ENABLED ?= 0
BUILD_DIR ?= ./build
MAIN_DIR ?= ./function

GO_FILES ?= $(shell find . -name '*.go')
BUILD_FILE ?= main-email

check-go:
	@which go > /dev/null || (echo "Go not found. Please install Go." && exit 1)

.PHONY: build

build:
	@go env -w GOOS=$(GOOS)
	@go env -w GOARCH=$(GOARCH)
	@go env -w CGO_ENABLED=$(CGO_ENABLED)
	go build -o $(BUILD_DIR)/$(BUILD_FILE) $(GO_FILES)

zip: 
	@which zip > /dev/null || (echo "zip not found. Please install zip." && exit 1)
	zip -FS -j $(BUILD_DIR)/${BUILD_FILE}.zip $(BUILD_DIR)/${BUILD_FILE}

to-s3:
	aws s3 sync $(BUILD_DIR)/ s3://$(BUCKET_NAME) --exclude "*" --include "*.zip"

to-lambda:
	aws lambda update-function-code --function-name ${FUNC_NAME} --s3-bucket ${BUCKET_NAME} --s3-key ${BUILD_FILE}.zip --no-cli-pager

invoke-lambda:
	aws lambda invoke --function-name ${FUNC_NAME} out --log-type Tail --query 'LogResult' --output text |  base64 -d
	rm out

clean-bin:
	cd $(BUILD_DIR) && find . ! -name '*.zip' -type f -exec rm -f {} +

clean-all:
	rm -f $(BUILD_DIR)/*

update-ses-template:
	aws ses update-template --cli-input-json file://Templates/emailTemplate.json

