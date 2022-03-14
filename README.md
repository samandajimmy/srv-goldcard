# srv-goldcard

Goldcard Service for Pegadaian Digital Service

## Local Development

### Prerequisites

1. Goland IDE or Visual Studio Code
2. Go 1.17
3. UNIX Shell
   > Use `wsl2` in Windows 10
4. Git
5. Make
6. Docker and Docker Compose CE

### Quick Start

```shell
# Set-up development environment
make configure

# Run database upgrade migration scripts
make db-up

# Start development server
make serve
```