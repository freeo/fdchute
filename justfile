default:
  @just --choose

# benchmark
time PATH:
  go build main.go
  time (fd -H --ignore-file ~/.config/fd/ignore_dirchute . {{PATH}} | ./main | wc -l)

timebaseline:
  time (fd -H . ~/ | awk -F'/' '{print NF-1 ($(NF)=="" ? "0" : "1"), $0}' | sort -n | cut -d' ' -f2- | wc -l )

timewb:
  @just time ~/wb

timehome:
  @just time ~/

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

# Install the application, depending on a fresh build
install: build
    @echo "Installing application..."
    @sudo cp -f main /usr/local/bin/dirchute
    @echo "Installation complete. You can now run 'dirchute' from anywhere."

# Uninstall the application
uninstall:
    @echo "Uninstalling application..."
    @sudo rm -f /usr/local/bin/dirchute
    @echo "Uninstallation complete."
