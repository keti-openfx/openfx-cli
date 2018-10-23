package log

import (
	"fmt"
	"os"

	"github.com/morikuni/aec"
)

func Fatalln(msg interface{}) {
	os.Stderr.WriteString(aec.RedF.Apply(fmt.Sprintf("%v\n", msg)))
	os.Exit(1)
}

func Fatal(msg ...interface{}) {
	os.Stderr.WriteString(aec.RedF.Apply(fmt.Sprintf(msg[0].(string), msg[1:]...)))
	os.Exit(1)
}

func Error(msg ...interface{}) {
	os.Stderr.WriteString(aec.RedF.Apply(fmt.Sprintf(msg[0].(string), msg[1:]...)))
}

func Warn(msg ...interface{}) {
	os.Stderr.WriteString(aec.LightRedF.Apply(fmt.Sprintf(msg[0].(string), msg[1:]...)))
}

func Debug(msg ...interface{}) {
	fmt.Printf(aec.LightBlueF.Apply(fmt.Sprintf(msg[0].(string), msg[1:]...)))
}

func Info(format string, a ...interface{}) {
	fmt.Printf(aec.BlueF.Apply(fmt.Sprintf(format, a...)))
}

func Print(format string, a ...interface{}) {
	fmt.Printf(fmt.Sprintf(format, a...))
}
