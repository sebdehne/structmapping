# Struct Mapping library for Go

Small library which lets you copy data from one struct (tree) to another which may be of a different type. This can be useful for example when mapping between an external/internal domain model.

## Example

Given the following example types:

```Go
type OrderA struct {
	Items map[string]OrderItemA
	Tags []string
}

type OrderItemA struct {
	Quantity  int
	ItemPrice MoneyA
}

type MoneyA struct {
	krone int
	ore   int
}

type OrderB struct {
	Items map[string]OrderItemB
	Tags []string
}

type OrderItemB struct {
	Quantity  int
	ItemPrice MoneyB
}

type MoneyB struct {
	krone int
	ore   int
}
```

To convert from OrderA to OrderB:

```Go
// Setup the mapper
m := New(SrcFieldBased, false)

// with a customer mapping function for MoneyA->MoneyB
m.Add(func(src MoneyA, dst *MoneyB) {
    dst.krone = src.krone
    dst.ore = src.ore
})

// generate some sample data
src := genTestOrder()

// map now
result := new(OrderB)
m.Map(&src, result)
// "result" is now populated

```