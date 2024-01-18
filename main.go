// ball.go
// /bin/true && exec /usr/bin/env go run "$0" "$@"
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	args := os.Args

	// Print other command-line arguments
	fmt.Println("Command-line arguments:")
	for i, arg := range args[0:] {
		fmt.Printf("%d: %s\n", i+1, arg)
	}
	ballerinaHome := "/path/to/ballerina/home"
	javaPath := filepath.Join(ballerinaHome, "../../dependencies/jdk-17.0.7+7-jre")
	if stat, err := os.Stat(javaPath); err == nil && stat.IsDir() {
		javaHome := javaPath
		// Set JAVA_HOME environment variable
		os.Setenv("JAVA_HOME", javaHome)
	}
	fmt.Println(runtime.GOOS)
	var cygwin, mingw, os400, darwin bool
	switch runtime.GOOS {
	case "darwin":
		darwin = true
		javaHome := os.Getenv("JAVA_HOME")
		javaVersion := os.Getenv("JAVA_VERSION")
		if javaHome == "" {
			if javaVersion == "" {
				// If both JAVA_HOME and JAVA_VERSION are not set, get the default Java home
				javaHomeCmd, _ := exec.Command("/usr/libexec/java_home").Output()
				javaHome = string(javaHomeCmd)
				os.Setenv("JAVA_HOME", javaHome)
			} else {
				// If JAVA_HOME is not set, but JAVA_VERSION is set, get the Java home for the specified version
				javaHomeCmd := exec.Command("/usr/libexec/java_home", "-v", javaVersion)
				javaHomeOutput, err := javaHomeCmd.Output()
				if err == nil {
					javaHome = string(javaHomeOutput)
					os.Setenv("JAVA_HOME", javaHome)
				} else {
					fmt.Printf("Error determining Java home for version %s: %v\n", javaVersion, err)
					os.Exit(1)
				}
			}
		}
	case "os400":
		os400 = true
	}

	if _, ok := os.LookupEnv("CYGWIN"); ok {
		cygwin = true // Running in Cygwin environment
	}

	if _, ok := os.LookupEnv("MINGW"); ok {
		mingw = true // Running in MinGW environment
	}

	fmt.Println("Cygwin:", cygwin)
	fmt.Println("MinGW:", mingw)
	fmt.Println("OS400:", os400)
	fmt.Println("Darwin:", darwin)

	prg, err := os.Executable()
	prg = "/usr/lib/ballerina/distributions/ballerina-2201.8.4/bin/bal"
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	// Resolve symbolic links
	for {
		link, err := filepath.EvalSymlinks(prg)
		if err != nil {
			fmt.Println("Error resolving symbolic link:", err)
			return
		}

		if link == prg {
			break
		}

		prg = link
	}

	// Get the directory of the script
	prgDir := filepath.Dir(prg)

	// Set BALLERINA_HOME
	ballerinaHome, err = filepath.Abs(filepath.Join(prgDir, ".."))
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	// Output the result
	fmt.Println("BALLERINA_HOME:", ballerinaHome)

	// Get the directory containing the executable

	//give realpath and execDir
	//execDir = "/usr/lib/ballerina/distributions/ballerina-2201.8.4/bin"
	//realPath = "/usr/lib/ballerina/distributions/ballerina-2201.8.4/bin/bal"

	// absParentDir, err := filepath.Abs(execDir)
	// if err != nil {
	// 	fmt.Println("Error getting absolute path:", err)
	// 	return
	// }

	//fmt.Println("BALLERINA_HOME:", absParentDir)

	//Setting Java CMD

	javaCmd := os.Getenv("JAVACMD")
	javaHome := os.Getenv("JAVA_HOME")
	//fmt.Println("JavaCmD= ", javaCmd)
	fmt.Println("-----------------------------")
	if javaCmd == "" {
		if javaHome != "" {
			// If JAVA_HOME is set, determine the appropriate JAVACMD
			var javacmdPath string

			// Check if the specified Java home is for IBM's JDK on AIX
			if _, err := os.Stat(filepath.Join(javaHome, "jre", "sh", "java")); err == nil {
				javacmdPath = filepath.Join(javaHome, "jre", "sh", "java")
			} else {
				// Use the standard Java executable path
				javacmdPath = filepath.Join(javaHome, "bin", "java")
			}

			// Check if the determined JAVACMD path is executable
			if _, err := os.Stat(javacmdPath); err == nil {
				javaCmd = javacmdPath
			} else {
				fmt.Printf("Error: Unable to find executable JAVACMD in specified JAVA_HOME: %s\n", javaHome)
				os.Exit(1)
			}
		} else {
			// If neither JAVACMD nor JAVA_HOME is set, use the default 'java' command
			javaCmd = "java"
		}
	}

	javaCmd = strings.TrimSpace(javaCmd)

	//fmt.Printf("JAVACMD: %s\n", javaCmd)
	//fmt.Printf("JAVAHOME: %s\n", javaHome)

	//assign value to ballerinaCliWidth
	var ballerinaCLIWidth string

	// Check if tput is available
	tputCmd := exec.Command("type", "tput")
	if tputCmd.Run() != nil {
		// tput is not available, set default width
		ballerinaCLIWidth = "80"
	} else {
		// tput is available, check if it can determine the number of columns
		colsCmd := exec.Command("tput", "cols")
		output, err := colsCmd.Output()

		if err != nil {
			// tput cols failed, set default width
			ballerinaCLIWidth = "80"
		} else {
			// tput cols succeeded, parse the output as an integer
			width, err := strconv.Atoi(string(output))
			if err != nil {
				// Failed to parse, set default width
				ballerinaCLIWidth = "80"
			} else {
				// Set the width to the parsed value
				ballerinaCLIWidth = strconv.Itoa(width)
			}
		}
	}
	fmt.Println("CLI Width :", ballerinaCLIWidth)

	//Setting Ballerina debug port

	balJavaDebug := os.Getenv("BAL_JAVA_DEBUG")
	if balJavaDebug != "" {
		// BAL_JAVA_DEBUG is set, proceed
		if os.Getenv("JAVA_OPTS") != "" {
			fmt.Println("Warning !!!. User specified JAVA_OPTS may interfere with BAL_JAVA_DEBUG")
		}

		// Check if BAL_DEBUG_OPTS is set
		balDebugOpts := os.Getenv("BAL_DEBUG_OPTS")
		if balDebugOpts != "" {
			// If BAL_DEBUG_OPTS is set, use it
			os.Setenv("JAVA_OPTS", os.Getenv("JAVA_OPTS")+" "+balDebugOpts)
		} else {
			// If BAL_DEBUG_OPTS is not set, set default debug options
			os.Setenv("JAVA_OPTS", os.Getenv("JAVA_OPTS")+
				" -Xdebug -Xnoagent -Djava.compiler=NONE -Xrunjdwp:transport=dt_socket,server=y,suspend=y,address="+balJavaDebug)
		}
	} else {
		fmt.Println("BAL_JAVA_DEBUG is not set.")
	}

	//Ballerina jarpath setting
	//ballerinaHome = "/usr/lib/ballerina/distributions/ballerina-2201.8.4"
	jarPath := filepath.Join(ballerinaHome, "bre", "lib")
	var ballerinaJarPath string

	// Get a list of all JAR files in the jarPath directory
	files, err := filepath.Glob(filepath.Join(jarPath, "*.jar"))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	pattern := filepath.Join(ballerinaHome, "bre", "lib", "*.jar")

	// Iterate over the JAR files
	for _, f := range files {
		matched, err := filepath.Match(pattern, f)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		if !matched {
			ballerinaJarPath += ":" + f
		}
	}

	//fmt.Printf("BALLERINA_JARPATH: %s\n", ballerinaJarPath)

	//Ballerina Classpath Setting
	var ballerinaClasspath string
	toolsJarPath := filepath.Join(ballerinaHome, "bre", "lib", "tools.jar")
	if _, err := os.Stat(toolsJarPath); err == nil {
		ballerinaClasspath = filepath.Join(javaHome, "lib", "tools.jar")
	}
	// Add JAR files from bre/lib
	//brePath := filepath.Join(ballerinaHome, "bre", "lib")
	breFiles, err := filepath.Glob(filepath.Join(jarPath, "*.jar"))
	if err == nil {
		//fmt.Println("Files in bre folder")
		for _, f := range breFiles {
			//fmt.Println(f)
			ballerinaClasspath += string(filepath.ListSeparator) + f
		}
	}

	// Add JAR files from lib/tools/lang-server/lib
	langServerPath := filepath.Join(ballerinaHome, "lib", "tools", "lang-server", "lib")
	langServerFiles, err := filepath.Glob(filepath.Join(langServerPath, "*.jar"))
	if err == nil {
		for _, f := range langServerFiles {
			ballerinaClasspath += string(filepath.ListSeparator) + f
		}
	}

	// Add JAR files from lib/tools/debug-adapter/lib
	debugAdapterPath := filepath.Join(ballerinaHome, "lib", "tools", "debug-adapter", "lib")
	debugAdapterFiles, err := filepath.Glob(filepath.Join(debugAdapterPath, "*.jar"))
	if err == nil {
		for _, f := range debugAdapterFiles {
			ballerinaClasspath += string(filepath.ListSeparator) + f
		}
	}

	//fmt.Println("BALLERINA_CLASSPATH:", ballerinaClasspath)
	//jarFilePath := "/usr/lib/ballerina/distributions/ballerina-2201.8.4/bre/lib/ballerina-cli-2201.8.4.jar"

	// if ballerinaClasspath != "" {
	// 	classpathEntries := strings.Split(ballerinaClasspath, string(filepath.ListSeparator))
	// 	for _, entry := range classpathEntries {
	// 		if entry == jarFilePath {
	// 			fmt.Printf("The JAR file '%s' is in BALLERINA_CLASSPATH.\n", jarFilePath)
	// 			return
	// 		}
	// 	}
	// }

	// fmt.Printf("The JAR file '%s' is not in BALLERINA_CLASSPATH.\n", jarFilePath)

	//Define Ballerina CLI Arguments
	cmdArgs := []string{
		"-Xms256m", "-Xmx1024m",
		"-classpath", ballerinaClasspath,
		"-Dballerina.home=" + ballerinaHome,
		"-Denable.nonblocking=false",
		"-Djava.security.egd=file:/dev/./urandom",
		"-Dfile.encoding=UTF8",
		"-Dballerina.target=jvm",
		"-Djava.command=" + javaCmd,
	}
	//fmt.Println(cmdArgs)
	//Implementing the running method for bal
	cmdArgs = append(cmdArgs, os.Args[1:]...)
	fmt.Println("length", len(os.Args))

	//Implementation of bal run --debug=<PORT> <*.jar>
	// handles "bal run --debug <PORT> <JAR_PATH>" command
	//Implementation of bal run <*.jar>
	if len(os.Args) == 3 && os.Args[1] == "run" && isJarFile(os.Args[2]) {
		jarFilePath := os.Args[2]
		cmd := exec.Command("java", "-jar", jarFilePath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		err := cmd.Run()
		if err != nil {
			fmt.Println("Error running bal run jar command:", err)
			os.Exit(1)
		}
	} else if len(os.Args) >= 4 && os.Args[1] == "run" && os.Args[2] == "--debug=*=" && isJarFile(os.Args[3]) {
		debugPort, err := extractDebugPort(os.Args[2])
		if err != nil {
			fmt.Println("Error: Invalid debug port number specified.")
			os.Exit(1)
		}

		if validateDebugPort(debugPort) {
			cmdArgs := []string{"-jar"}
			cmdArgs = append(cmdArgs, os.Args[3:]...)

			cmd := exec.Command("java", "-agentlib:jdwp=transport=dt_socket,server=y,suspend=y,address="+strconv.Itoa(debugPort))
			cmd.Args = append(cmd.Args, cmdArgs...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				fmt.Println("Error running Java command:", err)
				os.Exit(1)
			}
		}
	} else if len(os.Args) == 5 && os.Args[1] == "run" && os.Args[2] == "--debug" && isJarFile(os.Args[4]) {
		debugPort, err := extractDebugPort(os.Args[3])
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}

		if validateDebugPort(debugPort) {
			cmdArgs := []string{"-jar"}
			cmdArgs = append(cmdArgs, os.Args[4:]...)

			cmd := exec.Command("java", "-agentlib:jdwp=transport=dt_socket,server=y,suspend=y,address="+strconv.Itoa(debugPort))
			cmd.Args = append(cmd.Args, cmdArgs...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				fmt.Println("Error running Java command:", err)
				os.Exit(1)
			}
		} else {
			fmt.Println("Error: Invalid debug port number specified.")
			os.Exit(1)
		}
	} else {
		// Create the command
		MainClass := "io.ballerina.cli.launcher.Main"
		//Args := append([]string{"-cp", ballerinaClasspath, MainClass}, os.Args[1:]...)
		//cmd := exec.Command("java", Args...)
		cmd := exec.Command("java", "-cp", ballerinaClasspath, MainClass)

		// Append user-provided arguments to cmd.Args
		cmd.Args = append(cmd.Args, os.Args[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()

		// Run the command

		if err != nil {
			fmt.Println("Error running Java command:", err)
			os.Exit(1)
		}
	}
	///home/pabadhi/myGoProjects/sample2/app/build/libs/app.jar

}

func isJarFile(filePath string) bool {
	return len(filePath) > 4 && filePath[len(filePath)-4:] == ".jar"
}

func extractDebugPort(debugArg string) (int, error) {
	re := regexp.MustCompile(`--debug=([0-9]+)`)
	matches := re.FindStringSubmatch(debugArg)
	if len(matches) != 2 {
		return 0, fmt.Errorf("invalid --debug argument format")
	}

	return strconv.Atoi(matches[1])
}

func validateDebugPort(port int) bool {
	return port > 0
}
