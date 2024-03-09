GORUN = go run
NODEMON = nodemon --exec
SWAGFMT = swag fmt
SWAGBUILD = swag init
TARGET = main.go

run:
	$(NODEMON) $(GORUN) $(TARGET)