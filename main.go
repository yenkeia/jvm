package main

import "fmt"

func main() {
	var cmd *Cmd
	cmd = parseCmd()
	if cmd.versionFlag {
		fmt.Println("version 0.0.1")
	} else if cmd.helpFlag || cmd.class == "" {
		printUsage()
	} else {
		startJVM(cmd)
	}
}

// jvm -version
// version 0.0.1
// jvm -cp foo/bar MyApp arg1 arg2
// classpath:foo/bar class:MyApp args:[arg1 arg2]

func startJVM(cmd *Cmd) {
	fmt.Printf("classpath:%s class:%s args:%v\n", cmd.cpOption, cmd.class, cmd.args)
}
