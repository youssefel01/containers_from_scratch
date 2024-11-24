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
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("Help!")
	}
}

// run() going to:
// going to create the new namespaces UTS
// and it's going create a new process in which going to run the same program again
// but substittuing the parameter "child" instead of "run"
func run() {
	fmt.Printf("Running %v as PID %d\n ", os.Args[2:], os.Getpid())
	// cmd is the structure that describes the command i want run
	// "/proc/self/exe" = run the same program
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 1-we going start adding a namespaceS here using syscall
	//  SysProcAttr = system called attributes that pass in when i want run the command
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

	// create a Cgroup
	cg()

	// 2-set the hostname
	syscall.Sethostname([]byte("myContainer1"))
	// 3-make our own root file
	must(syscall.Chroot("/home/abstract/Documents/UBUNTU"))
	// when you change the root directory it leaves you in some undefined directory
	// so we have to spicifiyted explistly
	must(syscall.Chdir("/"))
	// we need to mount
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())

	syscall.Unmount("proc", 0)

}

// Cgroup
func cg() {
	pids := "/sys/fs/cgroup/"
	os.Mkdir(filepath.Join(pids, "youssef"), 0755)                                                           // create a Cgroup
	must(os.WriteFile(filepath.Join(pids, "youssef/pids.max"), []byte("20"), 0700))                          // set the maximum processes in the cgroup
	must(os.WriteFile(filepath.Join(pids, "youssef/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700)) // assign the current running process to the Cgroup
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
