package structmapping

import (
	"testing"
	"fmt"
)

func TestMakeMaps(t *testing.T) {

	o := NewWithMaps(new(OrderA)).(*OrderA)

	oi := OrderItemA{ItemPrice:MoneyA{krone:10,ore:0}}
	o.Items["not there"] = oi

	fmt.Println(o)
}
