GO=go
FILES=api.go models.go
D_FILES=downloader.go models.go
C_FILES=control.go models.go
BW_FILES=bigWigToWig.go
BAM_FILES=bamToWig.go
FILES_ONE=one.go models.go
STEREO_FILES=stereogene.go models.go
RUN=run
EXPORTS=GOARCH=amd64 GOOS=linux
BUILD=build
OUT=api
D_OUT=downloader
C_OUT=control
BW_OUT=bigWigToWigGo
BAM_OUT=bamToWigGo
STEREO_OUT=stereogeneGo
PRE_CMD=rm -rf encode.files.txt

downloader:
	$(PRE_CMD) | $(GO) $(RUN) $(D_FILES)
downloader_release:
	$(PRE_CMD) | $(EXPORTS) $(GO) $(BUILD) -o $(D_OUT) $(D_FILES)
control:
	$(PRE_CMD) | $(GO) $(RUN) $(C_FILES)
control_release:
	$(PRE_CMD) | $(EXPORTS) $(GO) $(BUILD) -o $(C_OUT) $(C_FILES)
bigwig:
	$(PRE_CMD) | $(GO) $(RUN) $(BW_FILES)
bigwig_release:
	$(PRE_CMD) | $(EXPORTS) $(GO) $(BUILD) -o $(BW_OUT) $(BW_FILES)
bam:
	$(PRE_CMD) | $(GO) $(RUN) $(BAM_FILES)
bam_release:
	$(PRE_CMD) | $(EXPORTS) $(GO) $(BUILD) -o $(BAM_OUT) $(BAM_FILES)
run:
	$(PRE_CMD) | $(GO) $(RUN) $(FILES)
one:
	$(PRE_CMD) | $(GO) $(RUN) $(FILES_ONE)
release:
	$(PRE_CMD) | $(EXPORTS) $(GO) $(BUILD) -o $(OUT) $(FILES)
stereogene_release:
	$(PRE_CMD) | $(EXPORTS) $(GO) $(BUILD) -o $(STEREO_OUT) $(STEREO_FILES)