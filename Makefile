GO=go
FILES=api.go models.go
FILES_ONE=one.go models.go
RUN=run
EXPORTS=GOARCH=amd64 GOOS=linux
BUILD=build
OUT=api
PRE_CMD=rm -rf encode.files.txt

run:
	$(PRE_CMD) | $(GO) $(RUN) $(FILES)
one:
	$(PRE_CMD) | $(GO) $(RUN) $(FILES_ONE)
release:
	$(PRE_CMD) | $(EXPORTS) $(GO) $(BUILD) -o $(OUT) $(FILES)
