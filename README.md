# Portify

Portify is a simple Go-based port scanner that enables you to scan a range of ports on a target host. This tool provides information about open and closed ports and, if desired, can retrieve banners from open ports.

## Installation

To use Portify, follow these steps:

1. Clone this repository to your local machine:

   ```bash
   git clone https://github.com/ni5arga/portify.git
   ```
2. Navigate to the project direcory
   ```bash
   cd portify
   ```
3. Build the project with Go
   ```bash
   go build portify.go
   ``` 
You can now use the `portify` binary to scan ports. 

## Command Line Usage 
```bash
Usage: ./portify [options] <host> <start-port> <end-port>

Options:
  -timeout string
        Connection timeout duration (default "3s")
  -parallel int
        Number of parallel scans (default 100)
  -show-closed
        Show closed ports in the output
  -show-banners
        Show banners from open ports
  -show-open
        Show open ports in the output
```

- `<host>` : The target host you want to scan.
- `<start-port>`: The first port in the range you want to scan.
- `<end-port>` : The last port in the range you want to scan.
