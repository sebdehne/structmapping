package structmapping

import (
	"testing"
	"fmt"
)

func TestStructMapper_Map(t *testing.T) {

	m := New()
	m.Add(func(src MoneyA, dst *MoneyB) {
		dst.krone = src.krone
		dst.ore = src.ore
	})
	fmt.Println("Using following mapper", m)

	src := genTestOrder()
	result := new(OrderB)
	m.Map(&src, result)
	fmt.Println(result)
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