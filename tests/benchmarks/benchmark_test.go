package benchmarks

// this file is used  to help benchmark the impact of recreating the runspace every time, or if we should cache them.
// go test . -count 1 -v -run .*Power.*

import (
	"testing"

	"sync"

	"github.com/KnicKnic/go-powershell/pkg/powershell"
)

func executePowershellFullCreate(script string, namedArgs map[string]interface{}) (obj []string, err error) {

	runspace := powershell.CreateRunspaceSimple()
	defer runspace.Close()
	return executePowershellCommon(runspace, script, namedArgs)
}

func executePowershellBuffered(script string, namedArgs map[string]interface{}) (obj []string, err error) {
	runspace := <-runspaceInstances
	defer func() { runspaceInstances <- runspace }()

	return executePowershellCommon(runspace, script, namedArgs)
}

func executePowershellCommon(runspace powershell.Runspace, script string, namedArgs map[string]interface{}) (obj []string, err error) {

	results := runspace.ExecScript(script, true, namedArgs)
	defer results.Close()

	if !results.Success() {
		exception := results.Exception.ToString()
		panic(exception)
		//return []string{}, errors.New(exception)
	}

	strResults := make([]string, len(results.Objects), len(results.Objects))
	for i, result := range results.Objects {
		strResults[i] = result.ToString()
	}

	return strResults, nil
}

// after running a few tests, I do not believe that buffering powershell sessions is worth the complexity of the code.
// with a 10 unit worker pool, running 100 powershell tests took .65 seconds
// sequentially running 100 tests and creating and destroying the runspace every time it took 3.26 seconds

var runspaceInstances chan powershell.Runspace

func init() {
	runspaceInstanceCount := 10
	runspaceInstances = make(chan powershell.Runspace, runspaceInstanceCount)
	for i := 0; i < runspaceInstanceCount; i++ {
		runspaceInstances <- powershell.CreateRunspaceSimple()
	}
}

var powershellTestRunSize int = 100

var powershellCommand = "start-sleep -milliseconds 100"

// var powershellCommand = "$a = 1+1"

func Test_LotsPowershellParallelBuffered(t *testing.T) {

	var wg sync.WaitGroup
	simplePowershellCommand := func(wg *sync.WaitGroup) { executePowershellBuffered(powershellCommand, nil); wg.Done() }

	for i := 0; i < powershellTestRunSize; i++ {
		wg.Add(1)
		go simplePowershellCommand(&wg)
	}
	wg.Wait()

}
func Test_LotsPowershellParallel(t *testing.T) {

	var wg sync.WaitGroup
	simplePowershellCommand := func(wg *sync.WaitGroup) { executePowershellFullCreate(powershellCommand, nil); wg.Done() }

	for i := 0; i < powershellTestRunSize; i++ {
		wg.Add(1)
		go simplePowershellCommand(&wg)
	}
	wg.Wait()

}

func Test_LotsPowershellSequentialBuffered(t *testing.T) {

	for i := 0; i < powershellTestRunSize; i++ {
		executePowershellBuffered(powershellCommand, nil)
	}

}

func Test_LotsPowershellSequential(t *testing.T) {

	for i := 0; i < powershellTestRunSize; i++ {
		executePowershellFullCreate(powershellCommand, nil)
	}
}
