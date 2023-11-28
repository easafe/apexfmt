package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/antlr4-go/antlr/v4"
	"github.com/octoberswimmer/apexfmt/parser"
)

type TreeShapeListener struct {
	indentLevel int
	*parser.BaseApexParserListener
}

type Visitor struct {
	indentLevel int
	parser.BaseApexParserVisitor
}

func NewVisitor() *Visitor {
	return &Visitor{}
}

func (v *Visitor) visitRule(node antlr.RuleNode) interface{} {
	return node.Accept(v)
}

func (v *Visitor) VisitCompilationUnit(ctx *parser.CompilationUnitContext) interface{} {
	fmt.Fprintln(os.Stderr, "HERE WE GO!")
	t := ctx.TypeDeclaration()
	// TODO: Handle the typeDeclaration modifiers
	switch {
	case t.ClassDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "IT'S A CLASS")
		return v.visitRule(t.ClassDeclaration()).(string)
	case t.InterfaceDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "IT'S AN INTERFACE")
		return v.visitRule(t.InterfaceDeclaration()).(string)
	case t.EnumDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "IT'S AN ENUM")
	}
	return ""
}

func (v *Visitor) VisitClassDeclaration(ctx *parser.ClassDeclarationContext) interface{} {
	fmt.Fprintln(os.Stderr, "IN THE CLASS DECLARATION")
	// TODO: handle extends and implements
	return fmt.Sprintf("class %s {\n%s\n}", ctx.Id().GetText(), v.visitRule(ctx.ClassBody()))
}

func (v *Visitor) VisitInterfaceDeclaration(ctx *parser.InterfaceDeclarationContext) interface{} {
	fmt.Fprintln(os.Stderr, "IN THE INTERFACE DECLARATION")
	// TODO: handle extends
	return fmt.Sprintf("interface %s {\n%s\n}", ctx.Id().GetText(), v.visitRule(ctx.InterfaceBody()))
}

func (v *Visitor) VisitInterfaceBody(ctx *parser.InterfaceBodyContext) interface{} {
	fmt.Fprintln(os.Stderr, "IN THE INTERFACE BODY")
	declarations := []string{}
	for _, d := range ctx.AllInterfaceMethodDeclaration() {
		declarations = append(declarations, v.visitRule(d).(string))
	}
	return strings.Join(declarations, "\n")
}

func (v *Visitor) VisitClassBody(ctx *parser.ClassBodyContext) interface{} {
	fmt.Fprintln(os.Stderr, "NEED TO DEAL WITH THE CLASS BODY")
	var cb []string
	for _, b := range ctx.AllClassBodyDeclaration() {
		cb = append(cb, v.visitRule(b).(string))
	}
	return strings.Join(cb, "\n")
}

func (v *Visitor) VisitClassBodyDeclaration(ctx *parser.ClassBodyDeclarationContext) interface{} {
	switch {
	case ctx.SEMI() != nil:
		return ";"
	case ctx.Block() != nil:
		// TODO: Handle static
		fmt.Fprintln(os.Stderr, "GOT A BLOCK")
	case ctx.MemberDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "IT'S A MEMBER")
		mods := []string{}
		for _, m := range ctx.AllModifier() {
			mods = append(mods, m.GetText())
		}
		return fmt.Sprintf("%s %s", strings.Join(mods, " "), v.visitRule(ctx.MemberDeclaration()))
	}
	return ""
}

func (v *Visitor) VisitMemberDeclaration(ctx *parser.MemberDeclarationContext) interface{} {
	fmt.Fprintln(os.Stderr, "IN MEMBER DECLARATION")
	switch {
	case ctx.MethodDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "FOUND METHOD DECLARATION")
		return v.visitRule(ctx.MethodDeclaration())
	case ctx.FieldDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "FOUND FIELD DECLARATION")
		return v.visitRule(ctx.FieldDeclaration())
	case ctx.ConstructorDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "FOUND CONSTRUCTOR DECLARATION")
		return v.visitRule(ctx.ConstructorDeclaration())
	case ctx.InterfaceDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "FOUND INTERFACE DECLARATION")
		return v.visitRule(ctx.InterfaceDeclaration())
	case ctx.ClassDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "FOUND CLASS DECLARATION")
		return v.visitRule(ctx.ClassDeclaration())
	case ctx.EnumDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "FOUND ENUM DECLARATION")
		return v.visitRule(ctx.EnumDeclaration())
	case ctx.PropertyDeclaration() != nil:
		fmt.Fprintln(os.Stderr, "FOUND PROPERTY DECLARATION")
		return v.visitRule(ctx.PropertyDeclaration())
	}
	fmt.Fprintln(os.Stderr, "FOUND UNEXPECTED DECLARATION")
	return ""
}

func (v *Visitor) VisitInterfaceMethodDeclaration(ctx *parser.InterfaceMethodDeclarationContext) interface{} {
	fmt.Fprintln(os.Stderr, "IN INTERFACE METHOD DECLARATION")
	returnType := "void"
	if ctx.TypeRef() != nil {
		returnType = v.visitRule(ctx.TypeRef()).(string)
	}
	mods := []string{}
	for _, m := range ctx.AllModifier() {
		mods = append(mods, m.GetText())
	}
	// TODO: formalParameters
	return fmt.Sprintf("%s %s %s();", strings.Join(mods, " "), returnType, ctx.Id().GetText())
}

func (v *Visitor) VisitFieldDeclaration(ctx *parser.FieldDeclarationContext) interface{} {
	fmt.Fprintln(os.Stderr, "IN FIELD DECLARATION")
	return fmt.Sprintf("%s %s;", v.visitRule(ctx.TypeRef()), v.visitRule(ctx.VariableDeclarators()))
}

func (v *Visitor) VisitVariableDeclarators(ctx *parser.VariableDeclaratorsContext) interface{} {
	fmt.Fprintln(os.Stderr, "IN VARIABLE DECLARATORS")
	vars := []string{}
	for _, vd := range ctx.AllVariableDeclarator() {
		vars = append(vars, v.visitRule(vd).(string))
	}
	return strings.Join(vars, ", ")
}

func (v *Visitor) VisitVariableDeclarator(ctx *parser.VariableDeclaratorContext) interface{} {
	fmt.Fprintln(os.Stderr, "IN VARIABLE DECLARATOR")
	decl := ctx.Id().GetText()
	if ctx.Expression() != nil {
		decl = fmt.Sprintf("%s = %s", decl, v.visitRule(ctx.Expression()))
	}
	return decl
}

func (v *Visitor) VisitMethodDeclaration(ctx *parser.MethodDeclarationContext) interface{} {
	fmt.Fprintln(os.Stderr, "IN METHOD DECLARATION")
	returnType := "void"
	if ctx.TypeRef() != nil {
		returnType = v.visitRule(ctx.TypeRef()).(string)
	}
	// TODO: formalParameters
	return fmt.Sprintf("%s %s() {\n%s\n}\n", returnType, ctx.Id().GetText(), "")
}

func (v *Visitor) VisitTypeRef(ctx *parser.TypeRefContext) interface{} {
	typeNames := []string{}
	for _, t := range ctx.AllTypeName() {
		// TODO: format typeList
		typeNames = append(typeNames, t.GetText())
	}

	return fmt.Sprintf("%s%s", strings.Join(typeNames, "."), ctx.ArraySubscripts().GetText())
}

func main() {
	input, _ := antlr.NewFileStream(os.Args[1])
	lexer := parser.NewApexLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)

	p := parser.NewApexParser(stream)
	// p.AddErrorListener(antlr.NewDiagnosticErrorListener(true))

	// fmt.Println(TreesIndentedStringTree(p.CompilationUnit(), "", nil, p))
	// antlr.ParseTreeWalkerDefault.Walk(NewTreeShapeListener(), p.CompilationUnit())
	v := NewVisitor()
	out, ok := p.CompilationUnit().Accept(v).(string)
	if !ok {
	}
	fmt.Println(out)
}
