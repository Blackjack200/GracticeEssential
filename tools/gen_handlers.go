package main

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/df-mc/dragonfly/server/player"
	"reflect"
	"strings"
)

func main() {
	t := reflect.TypeOf(player.NopHandler{})
	f := jen.NewFile("handler")
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		originalName := method.Name
		newInterfaceName := strings.TrimPrefix(method.Name, "Handle") + "Handler"

		f.Type().Id(newInterfaceName).Interface(
			jen.Id(originalName).Params(scanIn(method, f)...).Params(scanOut(method, f)...),
		)
	}
	for i := 0; i < t.NumMethod(); i++ {
		method := t.Method(i)
		originalName := method.Name
		newInterfaceName := strings.TrimPrefix(method.Name, "Handle") + "Handler"

		f.Func().
			Params(jen.Id("h").Id("*MultipleHandler")).Id(originalName).
			Params(scanIn(method, f)...).
			Params(scanOut(method, f)...).
			Block(
				jen.For(
					jen.List(jen.Id("hdr"), jen.Id("_")).Op(":=").Range().Id("h.handlers"),
				).Block(
					jen.If(
						jen.List(jen.Id("hdr"), jen.Id("ok")).Op(":=").Op("hdr").Assert(jen.Id(newInterfaceName)),
						jen.Id("ok"),
					).Block(
						jen.Id("hdr").Dot(originalName).Call(scanInParams(method)...),
					),
				),
			)
	}

	fmt.Printf("%#v", f)
}

func getQual(i int, mm reflect.Type, file *jen.File) jen.Code {
	return jen.Id(getParamName(i)).Id(mm.String())
}

func getParamName(i int) string {
	return fmt.Sprintf("param%v", i)
}

func scanIn(method reflect.Method, file *jen.File) []jen.Code {
	var arr []jen.Code
	for j := 1; j < method.Type.NumIn(); j++ {
		mm := method.Type.In(j)
		arr = append(arr, getQual(j, mm, file))
	}
	return arr
}

func scanInParams(method reflect.Method) []jen.Code {
	var arr []jen.Code
	for j := 1; j < method.Type.NumIn(); j++ {
		arr = append(arr, jen.Id(getParamName(j)))
	}
	return arr
}

func scanOut(method reflect.Method, file *jen.File) []jen.Code {
	var arr []jen.Code
	for j := 0; j < method.Type.NumOut(); j++ {
		mm := method.Type.Out(j)
		arr = append(arr, getQual(j, mm, file))
	}
	return arr
}
