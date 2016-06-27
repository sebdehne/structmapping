package structmapping

import (
	"testing"
	"fmt"
	"encoding/json"
)

func TestStructMapper_Map(t *testing.T) {

	m := New(SrcFieldBased, false)
	m.Add(func(src MoneyA, dst *MoneyB) {
		dst.krone = src.krone
		dst.ore = src.ore
	})
	m.Add(func(src MoneyB, dst *MoneyA) {
		dst.krone = src.krone
		dst.ore = src.ore
	})
	fmt.Println("Using following mapper", m)

	src := genTestOrder()

	result := new(OrderB)
	m.Map(&src, result)

	back := new(OrderA)
	m.Map(result, back)

	srcB, err := json.Marshal(src)
	if err != nil {
		t.Fatal(err)
	}
	resultB, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	backB, err := json.Marshal(back)
	if err != nil {
		t.Fatal(err)
	}

	if string(srcB) != string(resultB) || string(resultB) != string(backB) || string(srcB) != string(backB) {
		t.Fatal("Structs are not the same")
	}
}

func genTestOrder() OrderA {
	o := OrderA{}
	o.Items = make(map[string]OrderItemA)
	o.Items["item1"] = OrderItemA{Quantity:1, ItemPrice:MoneyA{krone:100, ore:99}}
	o.Tags = make([]string, 0)
	o.Tags = append(o.Tags, "test")
	return o
}

type OrderA struct {
	Items map[string]OrderItemA
	Tags []string
}

type OrderItemA struct {
	Quantity  uint8
	ItemPrice MoneyA
}

type MoneyA struct {
	krone uint8
	ore   uint8
}

type OrderB struct {
	Items map[string]OrderItemB
	Tags []string
}

type OrderItemB struct {
	Quantity  uint8
	ItemPrice MoneyB
}

type MoneyB struct {
	krone uint8
	ore   uint8
}