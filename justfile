default:
  @just --choose

test:
  grc gotestsum --format-icons=hivis --format=testdox

# End-2-End Test
# HOME | fdchute | TV (fresh build)
e2e:
  go build -o main
  fd -H --ignore-file ~/.config/fd/ignore_fdchute --base-directory ~/ . | ./main --debug | tv

benchmark:
  @just timehome

time PATH:
  go build main.go
  time (fd -H --ignore-file ~/.config/fd/ignore_fdchute . {{PATH}} | ./main | wc -l)

timewb:
  @just time ~/wb

timehome:
  @just time ~/



# fd baseline minimum required time
timefd:
  time (fd -H --ignore-file ~/.config/fd/ignore_fdchute . ~/ | wc -l)

# write full output to STDOUT time
timeoutput:
  go build main.go
  time (fd -H --ignore-file ~/.config/fd/ignore_fdchute . ~/ | ./main )

# initial baseline, ~30 seconds e2e. This is what I wanted to get rid of.
timebaseline:
  time (fd -H . ~/ | awk -F'/' '{print NF-1 ($(NF)=="" ? "0" : "1"), $0}' | sort -n | cut -d' ' -f2- | wc -l )

prof PATH:
  go build -o main
  time (./main -cpuprofile=cpu.prof -memprofile=mem.prof {{PATH}} | wc -l)
  # ./main -cpuprofile=cpu.prof -memprofile=mem.prof {{PATH}}

profwb:
  @just prof ~/wb

profhome:
  @just prof ~/

build:
    @echo "Building application..."
    @go build -o main
    @echo "Build complete."

install: build
    @echo "Installing application..."
    @sudo cp -f main /usr/local/bin/fdchute
    @echo "Installation complete. You can now run 'fdchute' from anywhere."

uninstall:
    @echo "Uninstalling application..."
    @sudo rm -f /usr/local/bin/fdchute
    @echo "Uninstallation complete."

# Use with emx:SPC d s: Go Dlv Remote Debug
dlv:
  fd -H --ignore-file ~/.config/fd/ignore_fdchute . ./ > input.txt
  dlv debug --headless --listen=:2345 --api-version=2 ./main.go -r stdin:input.txt

dlv_basic:
  fd -H --ignore-file ~/.config/fd/ignore_fdchute . ./ > input.txt
  dlv debug ./main.go -r stdin:input.txt
  # dlv exec ./main -r stdin:input.txt
  # fd -H --ignore-file ~/.config/fd/ignore_fdchute . ~/wb | dlv exec ./main -- -r STDIN


