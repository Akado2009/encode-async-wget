GO=go
FILES=api.go models.go
FILES_ONE=one.go models.go
CMD=run
PRE_CMD=rm -rf encode.files.txt

run:
	$(PRE_CMD) | $(GO) $(CMD) $(FILES)
one:
	$(PRE_CMD) | $(GO) $(CMD) $(FILES_ONE)