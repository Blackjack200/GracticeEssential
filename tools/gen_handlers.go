package main

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/df-mc/dragonfly/server/player"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
)

func main() {
	f := jen.NewFile("mhandler")

	if i, ok := debug.ReadBuildInfo(); ok {
		for _, m := range i.Deps {
			if m.Path == "github.com/df-mc/dragonfly" {
				gen(getDoc(m), f)
				break
			}
		}
	}
}

func gen(d *doc.Package, f *jen.File) {
	reflectionIface := reflect.TypeOf(player.NopHandler{})

	for _, t := range d.Types {
		if t.Name == "Handler" {
			ifaceType := t.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType)
			genFromInterface(ifaceType, reflectionIface, f)
			break
		}
	}
	fmt.Printf("%#v", f)
}

func genFromInterface(ifaceType *ast.InterfaceType, reflectionIface reflect.Type, f *jen.File) {
	for _, method := range ifaceType.Methods.List {
		for _, methodName := range method.Names {
			originalMethodName := methodName.Name
			typedIn, _ := getFuncIn(method, reflectionIface, originalMethodName)
			newInterfaceName := strings.TrimPrefix(originalMethodName, "Handle") + "Handler"
			var stmt []jen.Code
			for _, comm := range method.Doc.List {
				stmt = append(stmt, jen.Comment(comm.Text))
			}

			stmt = append(stmt, jen.Id(originalMethodName).Params(typedIn...))
			f.Type().Id(newInterfaceName).Interface(stmt...)
		}
	}

	for _, method := range ifaceType.Methods.List {
		for _, methodName := range method.Names {
			originalMethodName := methodName.Name
			typedIn, paramIn := getFuncIn(method, reflectionIface, originalMethodName)
			newInterfaceName := strings.TrimPrefix(originalMethodName, "Handle") + "Handler"

			f.Func().
				Params(jen.Id("h").Id("*MultipleHandler")).Id(originalMethodName).
				Params(typedIn...).
				Block(
					jen.For(
						jen.List(jen.Id("_"), jen.Id("hdr")).Op(":=").Range().Id("h.handlers"),
					).Block(
						jen.If(
							jen.List(jen.Id("hdr"), jen.Id("ok")).Op(":=").Op("hdr").Assert(jen.Id(newInterfaceName)),
							jen.Id("ok"),
						).Block(
							jen.Id("hdr").Dot(originalMethodName).Call(paramIn...),
						),
					),
				)
		}
	}
}

func getFuncIn(method *ast.Field, reflectionIface reflect.Type, originalMethodName string) ([]jen.Code, []jen.Code) {
	params := method.Type.(*ast.FuncType).Params.List
	var typedIn []jen.Code
	var paramIn []jen.Code
	seq := 1
	for _, param := range params {
		var names []jen.Code
		for _, pN := range param.Names {
			names = append(names, jen.Id(pN.Name))
		}
		reflectionMethod, found := reflectionIface.MethodByName(originalMethodName)
		if !found {
			panic(originalMethodName)
		}
		paramType := reflectionMethod.Type.In(seq).String()
		typedIn = append(typedIn, jen.List(names...).Id(paramType))
		paramIn = append(paramIn, names...)
		seq += len(names)
	}
	return typedIn, paramIn
}

func getDoc(m *debug.Module) *doc.Package {
	fset := token.NewFileSet()
	pkg, _ := parser.ParseFile(fset, filepath.Join(build.Default.GOPATH, "/pkg/mod/github.com/df-mc/dragonfly@"+m.Version+"/server/player/handler.go"), nil, parser.ParseComments)
	d, _ := doc.NewFromFiles(fset, []*ast.File{pkg}, filepath.Join(build.Default.GOPATH, "/pkg/mod/github.com/df-mc/dragonfly@"+m.Version+"/server/player"), doc.AllDecls)
	return d
}
