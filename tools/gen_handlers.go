package main

import (
	"github.com/dave/jennifer/jen"
	"github.com/df-mc/dragonfly/server/player"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"io"
	"os"
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
				code := gen(getDoc(m), f)
				w, err := os.OpenFile("./mhandler/generated.go", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
				if err != nil {
					panic(err)
				}
				if _, err := io.WriteString(w, code); err != nil {
					panic(err)
				}
				w.Close()
				break
			}
		}
	}
}

func gen(d *doc.Package, f *jen.File) string {
	reflectionIface := reflect.TypeOf(player.NopHandler{})

	for _, t := range d.Types {
		if t.Name == "Handler" {
			ifaceType := t.Decl.Specs[0].(*ast.TypeSpec).Type.(*ast.InterfaceType)
			genFromInterface(ifaceType, reflectionIface, f)
			break
		}
	}
	return f.GoString()
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
			newFieldName := "_" + newInterfaceName

			f.Func().
				Params(jen.Id("h").Id("*MultipleHandler")).Id(originalMethodName).
				Params(typedIn...).
				Block(
					jen.For(
						jen.List(jen.Id("_"), jen.Id("hdr")).Op(":=").Range().Id("h." + newFieldName),
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

	var fields []jen.Code
	var clearFields []jen.Code
	for _, method := range ifaceType.Methods.List {
		for _, methodName := range method.Names {
			originalMethodName := methodName.Name
			newInterfaceName := strings.TrimPrefix(originalMethodName, "Handle") + "Handler"
			fields = append(fields, jen.Id("_"+newInterfaceName).Id("[]"+newInterfaceName))
			clearFields = append(clearFields, jen.Id("h").Dot("_"+newInterfaceName).Op("=").Nil())
		}
	}
	f.Type().Id("MultipleHandler").Struct(fields...)

	blocks := []jen.Code{
		jen.Id("reg").Op(":=").False(),
		jen.Var().Id("funcs").Id("[]func()"),
	}
	for _, method := range ifaceType.Methods.List {
		for _, methodName := range method.Names {
			originalMethodName := methodName.Name
			newInterfaceName := strings.TrimPrefix(originalMethodName, "Handle") + "Handler"
			newFieldName := "_" + newInterfaceName
			blocks = append(blocks, jen.If(
				jen.List(jen.Id("hdr"), jen.Id("ok")).Op(":=").Op("hdr").Assert(jen.Id(newInterfaceName)),
				jen.Id("ok"),
			).Block(
				jen.Id("h").Dot(newFieldName).Op("=").Append(jen.Id("h").Dot(newFieldName), jen.Id("hdr")),
				jen.Id("reg").Op("=").True(),
				jen.Id("funcs").Op("=").Append(jen.Id("funcs"),
					jen.Func().Params().Block(
						jen.Id("h").Dot(newFieldName).Op("=").Id("deleteVal").Call(
							jen.Id("h").Dot(newFieldName),
							jen.Id("hdr"),
						),
					),
				),
			))
		}
	}
	blocks = append(blocks,
		jen.If(jen.Id("!reg")).
			Block(
				jen.Panic(jen.Lit("not a valid handler")),
			),
		jen.Return(jen.Func().Params().Block(
			jen.For(
				jen.List(jen.Id("_"), jen.Id("f")),
			).Op(":=").Range().Id("funcs").Block(jen.Id("f").Call()),
		)),
	)
	f.Func().
		Params(jen.Id("h").Id("*MultipleHandler")).Id("Register").
		Params(jen.Id("hdr").Any()).Id("func()").
		Block(blocks...)
	f.Func().
		Params(jen.Id("h").Id("*MultipleHandler")).Id("Clear").
		Params().
		Block(clearFields...)

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
		if paramType == "*event.Context[*github.com/df-mc/dragonfly/server/player.Player]" {
			paramType = "*event.Context[*player.Player]"
		}
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
