package tableify

import (
	"os"
	"testing"
)

type TestStruct struct {
	Value1 string
	Value2 string
}

func initTestStructTable() *StructTable {
	st, err := NewStructTable(os.Stdout, TestStruct{})
	if err != nil {
		panic(err)
	}

	return st
}

func TestHeader(t *testing.T) {
	st := initTestStructTable()

	expectedHeaders := []string{"Value1", "Value2"}
	if st.headers[0] != expectedHeaders[0] {
		t.Errorf("Expected %v, got %v", expectedHeaders, st.headers)
	}
}

func TestValues(t *testing.T) {
	st := initTestStructTable()

	st.Append(TestStruct{"A", "B"})
	if st.rows[0][0] != "A" {
		t.Errorf("Error: expected %v, got %v", "A", st.rows[0][0])
	}
}

func TestAppendBuild(t *testing.T) {
	st := initTestStructTable()
	st.AppendBulk([]TestStruct{
		TestStruct{"A", "B"},
		TestStruct{"C", "D"},
	})

	if st.rows[0][0] != "A" || st.rows[1][0] != "C" {
		t.Errorf("Error: \n\t->expected %v, got %v\n\t->expected %v, got %v", "A",
			st.rows[0][0], "C", st.rows[1][0])
	}
}
