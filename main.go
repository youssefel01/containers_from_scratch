package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: [run <command>]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("Unknown command!")
	}
}

// run() going to:
// create the new namespaces
// and create a new process in which going to run the same program again
// but substittuing the parameter "child" instead of "run"
func run() {
	fmt.Printf("Running %v as PID %d\n ", os.Args[2:], os.Getpid())

	// cmd is the structure that describes the command i want run
	// setup command to rerun the program with "child" argument
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// configure namespace isolation
	// SysProcAttr = system called attributes that pass in when i want run the command
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags:   syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS, // the system call that creating this new process that we're going to run our executable in
		Unshareflags: syscall.CLONE_NEWNS,                                               // to make the mount not visible to the system
	}

	// to run/create the process cmd
	must(cmd.Run())

}

// the child will:
// set the new hostname
// then run the actual command that we tried to run
func child() {
	fmt.Printf("Running %v as PID %d\n ", os.Args[2:], os.Getpid())

	// 1-create and configure the Cgroup
	cg()

	// 2-set the hostname
	syscall.Sethostname([]byte("myContainer1"))

	// 3-make our own root file
	must(syscall.Chroot("/home/abstract/Documents/UBUNTU"))

	// when you change the root directory it leaves you in some undefined directory
	// so we have to spicifiyted explistly
	must(syscall.Chdir("/"))

	// we need to mount
	mountSpecialFS()

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())

	// unmount
	cleanupMounts()

}

// Cgroup set up a c cgroup to restrict resources for the containerzed process
func cg() {
	pids := "/sys/fs/cgroup/"
	group := filepath.Join(pids, "youssef")

	// Create a new cgroup directory.
	os.Mkdir(group, 0755)

	// Set maximum allowed processes.
	must(os.WriteFile(filepath.Join(group, "pids.max"), []byte("20"), 0700))

	// Assign the current process to the cgroup.
	must(os.WriteFile(filepath.Join(group, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))

	// Add cleanup logic: Remove the cgroup after use.
	defer os.RemoveAll(group)
}

// mountSpecialFS
func mountSpecialFS() {
	fmt.Println("Mounting special filesystems... ")

	// Mount `./proc`
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	// Mount `./sys` as read-only for safety
	must(syscall.Mount("sysfs", "/sys", "sysfs", syscall.MS_RDONLY, ""))

	// Mount `/dev` as a tmpfs for device isolation
	must(syscall.Mount("tmpfs", "/dev", "tmpfs", 0, ""))

	fmt.Println("Filesystems mounted successfully")
}

// clean up Mounts
func cleanupMounts() {
	fmt.Println("Cleaning up mounts...")

	// Unmount `/proc`.
	if err := syscall.Unmount("/proc", 0); err != nil {
		fmt.Printf("Warning: failed to unmount /proc: %v\n", err)
	}

	// Unmount `/sys`.
	if err := syscall.Unmount("/sys", 0); err != nil {
		fmt.Printf("Warning: failed to unmount /sys: %v\n", err)
	}

	// Unmount `/dev`.
	if err := syscall.Unmount("/dev", 0); err != nil {
		fmt.Printf("Warning: failed to unmount /dev: %v\n", err)
	}

	fmt.Println("Mounts cleaned up.")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
