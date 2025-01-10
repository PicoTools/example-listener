BUILD_FLAGS = -trimpath
LD_FLAGS = -s -w 

listener:
	@mkdir -p bin/
	@echo "Building listener..."
	@CGO_ENABLED=0 go build $(BUILD_FLAGS) -ldflags="$(LD_FLAGS)" -o bin/listener ./cmd/listener

run: listener
	@bin/listener run 
