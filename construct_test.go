package main

import (
	"testing"
	"fmt"
	"io/ioutil"
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
	in, err := ioutil.ReadFile("test/construct1.json")
	if nil != err {
		t.Error("Can't read file test/construct1.json")
		return
	}

	var standard Forest
	err = json.Unmarshal(in,&standard)
	if nil != err {
		t.Error(err)
		return
	}
	if false == equal(forest.Roots, standard.Roots){
		t.Error("Fail")
		return
	}
	t.Log(forest.ToJson())
}

func Test_Construct_json(t *testing.T)  {
	rep := &ConstructRepresent{
		AddOps: []AddBucketRep{
			AddBucketRep{"root",[]string{"EBS_SHB"}},
			AddBucketRep{"room",[]string{"room1", "room2", "room3"}},
			AddBucketRep{"rack",[]string{"vrack1","vrack2"}},
			AddBucketRep{"host",[]string{"hostA","hostB"}},
		},
		MoveOps: []MoveRep{
			MoveRep{"root","EBS_SHB",[]string{"room1", "room2", "room3"}},
			MoveRep{"room","room1", []string{"vrack1"}},
			MoveRep{"room","room2",[]string{"vrack2"}},
			MoveRep{"rack","vrack1",[]string{"hostA"}},
			MoveRep{"rack","vrack2",[]string{"hostB"}},
		},
	}

	addOsdRepHostA := &AddOsdRep{"hostA",3}
	addOsdRepHostB := &AddOsdRep{"hostB",5}
	var forest Forest
	forest.New(rep)
	forest.AddOsds(addOsdRepHostA)
	forest.AddOsds(addOsdRepHostB)

	in, err := ioutil.ReadFile("test/construct2.json")
	if nil != err {
		t.Error("Can't read file test/construct2.json")
		return
	}

	var standard Forest
	err = json.Unmarshal(in,&standard)
	if nil != err {
		t.Error(err)
		return
	}

	if false == equal(forest.Roots, standard.Roots){
		t.Error("Fail")
		return
	}
	t.Log(forest.ToJson())
}

func equal(v1 interface{}, v2 interface{}) bool {
	b1, _ := json.Marshal(v1)
	b2,_ := json.Marshal(v2)
	if string(b1) != string(b2) {
		return false
	}
	return true
}