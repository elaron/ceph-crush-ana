package main

import (
	"testing"
	"io/ioutil"
	"encoding/json"
)

func Test_DoRule_1(t *testing.T)  {
	in, err := ioutil.ReadFile("test/rule_cluster1.json")
	if nil != err {
		t.Error("Can't read file test/rule_cluster1.json")
		return
	}

	var forest Forest
	err = json.Unmarshal(in, &forest)
	if nil != err {
		t.Error(err)
		return
	}

	var simpleRule SimpleRule
	simpleRule.Init(RuleRepresent{
		[]SelectStep{
			SelectStep{Type:"root", Num:1},
			SelectStep{Type:"room", Num:1},
			SelectStep{Type:"host", Num:3},
			SelectStep{Type:"osd", Num:1},
		},
	})

	selected := simpleRule.DoRule(&forest)
	if len(selected) != 3 {
		t.Error("Fail")
		return
	}
}