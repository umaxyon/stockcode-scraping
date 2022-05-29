package db

import (
	"fmt"
	"github.com/awslabs/goformation/v6"
	"log"
	"os"
	"os/exec"
	"testing"
)

var (
	WinStartCmd = []string{"/C", "docker-compose -f .\\test\\docker-compose.yml up -d"}
	WinStopCmd  = []string{"/C", "docker-compose -p test stop"}
)

func ReadDDL() {
	tmp, err := goformation.Open("../../dynamo_ddl.yaml")
	if err != nil {
		log.Fatalf("There was an error processing the template: %s", err)
	}
	fmt.Println(tmp)
}

func TestMain(m *testing.M) {
	// Before
	ReadDDL()
	err := exec.Command("cmd", WinStartCmd...).Start()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	//After
	_ = exec.Command("cmd", WinStopCmd...).Start()

	os.Exit(code)
}

func TestDynamo(t *testing.T) {
	t.Run("dynamo test", func(t *testing.T) {
		fmt.Println("hoge")
	})
}
