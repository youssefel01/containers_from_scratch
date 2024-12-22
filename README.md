# Custom Container Engine in Go

This project is a minimal implementation of a container engine written in Go. It demonstrates key concepts behind containerization, such as namespaces, chroot, and cgroups, to isolate and manage processes.

## Features
- **Process Isolation**: Uses Linux namespaces to isolate the containerized process.
- **Filesystem Isolation**: Implements a custom root filesystem using `chroot` and mounts special filesystems like `/proc`, `/sys`, and `/dev`.
- **Resource Control**: Limits the number of processes using cgroups.

## How It Works
The project creates a lightweight container by:
1. Spawning a new process in isolated namespaces.
2. Changing the root filesystem to an isolated directory.
3. Configuring cgroups to limit resource usage.
4. Running the specified command within the isolated environment.

## Usage
To run a command inside the containerized environment:

1. **Prepare a Root Filesystem**:
   - Create a directory (e.g., `/path/to/root`) containing necessary files like `/bin/bash` and libraries.

2. **Run the Program**:
   ```bash
   go run main.go run /bin/bash
