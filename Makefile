GO=go
FILES=api.go models.go
D_FILES=downloader.go models.go
FILES_ONE=one.go models.go
RUN=run
EXPORTS=GOARCH=amd64 GOOS=linux
BUILD=build
OUT=api
D_OUT=downloader
PRE_CMD=rm -rf encode.files.txt

downloader:
	$(PRE_CMD) | $(GO) $(RUN) $(D_FILES)
downloader_release:
	$(PRE_CMD) | $(EXPORTS) $(GO) $(BUILD) -o $(D_OUT) $(D_FILES)
run:
	$(PRE_CMD) | $(GO) $(RUN) $(FILES)
one:
	$(PRE_CMD) | $(GO) $(RUN) $(FILES_ONE)
release:
	$(PRE_CMD) | $(EXPORTS) $(GO) $(BUILD) -o $(OUT) $(FILES)
