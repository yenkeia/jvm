package main

import (
	"fmt"
	"github.com/yenkeia/jvm/classpath"
	"strings"
)

// jvm -version
// version 0.0.1

// jvm -cp foo/bar MyApp arg1 arg2
// classpath:foo/bar class:MyApp args:[arg1 arg2]

// jvm -Xjre "/Library/Java/JavaVirtualMachines/jdk1.8.0_111.jdk/Contents/Home/jre" java.lang.Object
// classpath:/Users/Mccree/gopath/src/github.com/yenkeia/jvm class:java.lang.Object args:[]
// class data:[202 254 186 190 0 0 0 52 0 78 7 0 49 10 0 1 0 ...

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

func startJVM(cmd *Cmd) {
	var (
		cp        *classpath.Classpath
		className string
		err       error
		classData []byte
	)
	cp = classpath.Parse(cmd.XjreOption, cmd.cpOption)
	fmt.Printf("classpath:%v class:%v args:%v\n", cp, cmd.class, cmd.args)
	className = strings.Replace(cmd.class, ".", "/", -1)
	classData, _, err = cp.ReadClass(className)
	if err != nil {
		fmt.Printf("Could not find or load main class %s\n", cmd.class)
		return
	}
	fmt.Printf("class data:%v\n", classData)
}
