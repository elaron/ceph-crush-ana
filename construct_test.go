package main

import (
	"testing"
	"fmt"
	"encoding/json"
)

func Test_Construct(t *testing.T)  {
	var forest Forest
	_, rootName:= forest.AddNode("room1", "room")
	for i:=0; i < 3; i++ {
		_, bucketName := forest.AddNode(fmt.Sprintf("host%d",i), "host")
		err := forest.MoveNode(bucketName,rootName)
		if nil != err {
			t.Error(err)
		}
		for j := 0; j < 10; j++ {
			_,osdName := forest.AddNode(fmt.Sprintf("osd.%s", i*j), "osd")
			err := forest.MoveNode(osdName, bucketName)
			if nil != err {
				t.Error(err)
			}
		}
	}

	b, err := json.Marshal(&forest)
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(string(b))
}