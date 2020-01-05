package logging

import "os"

func ExampleLogger() {
	SetLevel(InfoLevel)
	SetWriter(os.Stdout)
	SetFlags(LevelFlag | IdentifierFlag)
	log1 := Get("test1")
	log2 := Get("test2")

	p1 := 123
	p2 := "abc"
	log1.Error("error messsage ", p1, " ", p2)
	log2.Trace("trace messsage")
	log1.Warning("warning messsage")
	log1.Info("info messsage")
	log2.Warningf("warning messsage %v %v", p1, p2)
	log1.Debug("debug messsage")
	log1.Trace("trace messsage")
	log2.Info("1st line\n2nd line")

	// Output:
	// ERROR  |test1          |error messsage 123 abc
	// WARNING|test1          |warning messsage
	// INFO   |test1          |info messsage
	// WARNING|test2          |warning messsage 123 abc
	// INFO   |test2          |1st line\n2nd line
}
