package viper


import (
	"fmt"
	"time"
	"testing"
)

func intTest(t * testing.T) {
	val := GetInt("test_int")
	exp := 5
	if val != exp {
		t.Fatalf("test_int was [%v], expected [%v]", val, exp)
	}
}
func boolTest(t * testing.T) {
	val := GetBool("test_bool")
	exp := false
	if val != exp {
		t.Fatalf("test_bool was [%v], expected [%v]", val, exp)
	}
}

func stringTest(t * testing.T) {
	val := GetString("test_string")
	exp := "zaphod"
	if val != exp {
		t.Fatalf("test_string was [%s], expected [%s]", val, exp)
	}
}

func stringSliceTest(t * testing.T) {
	val := GetStringSlice("test_string_slice")
	t.Logf("%v", val)
	if val[0] != "foo" {
		t.Fatalf("expected foobar, got %s", val[0])
	}
	if val[1] != "bar" {
		t.Fatalf("expected foobar, got %s", val[1])
	}
	if val[2] != "quux" {
		t.Fatalf("expected foobar, got %s", val[2])
	}
}

func assertData(t * testing.T) {
	stringTest(t)
	boolTest(t)
	intTest(t)
	stringSliceTest(t)
}

func TestJSONConfig(t * testing.T) {
	SetConfigFile("test.json")
	err := ReadInConfig()
	if err != nil {
		t.Fatalf("problem with ReadInConfig for viper json %s", err)
	}

	// Test rest of reader methods here
	assertData(t)
}

func TestConsulConfig(t * testing.T) {
	SetConfigType("json") // basically unused now
	SetRemoteProvider("consul", "127.0.0.1:8500", "test/config")
	// Test watcher
	watcherChan, err := StartWatcher("127.0.0.1:8500", "test/config", 5)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
				case err := <-watcherChan:
					  if err != nil {
						  fmt.Println("error:")
						  fmt.Println(err)
					  } else {
						  fmt.Println("config should be reloaded now?")
						  assertData(t)
					  }
				case <-time.After(time.Second * 4):
					//close(watcherChan)
					return
			}
		}
	}()

	// Test rest of reader methods here
}
