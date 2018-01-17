package query

import (
	"fmt"
	"runtime/debug"
	"strings"
	"time"
)

var Debug = false

var Logger func(timestamp int64, cmd string, args []interface{})

func init() {
	Logger = func(timestamp int64, cmd string, args []interface{}) {
		spent := time.Duration((time.Now().UnixNano() - timestamp)) * time.Nanosecond
		fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " SQL query debug:")
		fmt.Println(cmd)
		fmt.Println("Args:")
		argsString := make([]string, len(args))
		for k := range args {
			argsString[k] = fmt.Sprint(args[k])
		}
		fmt.Println("[" + strings.Join(argsString, " , ") + "]")
		stacks := string(debug.Stack())
		stacklines := strings.Split(stacks, "\n")
		if len(stacklines) > 9 {
			fmt.Println("Stack:")
			fmt.Println(stacklines[7])
			fmt.Println(stacklines[8])
		}
		fmt.Println("Time spent:")
		if spent > time.Hour {
			fmt.Printf("%f %s \n", spent.Hours(), "hours")
		} else if spent > time.Minute {
			fmt.Printf("%f %s \n", spent.Minutes(), "minutes")
		} else if spent > time.Second {
			fmt.Printf("%f %s \n", spent.Seconds(), "seconds")
		} else if spent > time.Millisecond {
			fmt.Printf("%f %s \n", spent.Seconds()*1000, "milliseconds")
		} else if spent > time.Microsecond {
			fmt.Printf("%f %s \n", spent.Seconds()*1000*1000, "microseconds")
		} else {
			fmt.Printf("%d %s \n", spent.Nanoseconds(), "nanoseconds")
		}
		fmt.Println()
	}
}
