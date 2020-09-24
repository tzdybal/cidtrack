all: ~/.ipfs/plugins/cidtrack.so

clean:
	rm -rf ~/.ipfs/plugins/cidtrack.so ../cidtrack.so
	go clean

test: *.go
	gotestsum ./...

~/.ipfs/plugins/cidtrack.so: ../cidtrack.so
	cp ../cidtrack.so ~/.ipfs/plugins/cidtrack.so

../cidtrack.so: *.go
	$(MAKE) -C ../../.. plugin/plugins/cidtrack.so
