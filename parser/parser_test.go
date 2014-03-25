package parser

import (
    "testing"
    "reflect"
)


func TestEmptySource(t *testing.T) {
    New("")
    ast, _ := AST()

    if (reflect.TypeOf(ast) != reflect.TypeOf(new(NodeProgram))) {
        t.Errorf("AST should be a NodeProgram, found %s", reflect.TypeOf(ast).String())
    }

    if (ast.Child != nil) {
        t.Errorf("Program is not empty")
    }
}

func TestInvalidFirstNode(t *testing.T) {
    New("cloudson")
    _, error := AST()

    if error == nil {
        t.Errorf("Expected a syntax error")
    }
}

func TestValidFirstNode(t *testing.T) {
    New("select * from users")
    ast, _ := AST()

    if ast.Child == nil {
        t.Errorf("Program is empty")
    }
}

func TestUsingWildCard(t *testing.T) {
    New("select * from users")
    ast, error := AST()

    if error != nil {
        t.Errorf(error.Error())
    }

    if ast.Child == nil {
        t.Errorf("Program is empty")
    }

    selectNode := ast.Child.(*NodeSelect)
    if !selectNode.WildCard {
        t.Errorf("Expected wildcard setted")
    }
}

func TestUsingOneFieldName(t *testing.T) {
    New("select name from files")

    ast, error := AST()

    if error != nil {
        t.Errorf(error.Error())
    }

    selectNode := ast.Child.(*NodeSelect)
    
    if len(selectNode.Fields) != 1 {
        t.Errorf("Expected exactly one field and found %d", len(selectNode.Fields))
    }

    if selectNode.Fields[0] != "name" {
        t.Errorf("Expected param 'name' and found '%s'", selectNode.Fields[0])
    }
}

func TestUsingFieldNames(t *testing.T) {
    New("select name, created_at from files")

    ast, error := AST()

    if error != nil {
        t.Errorf(error.Error())
    }

    selectNode := ast.Child.(*NodeSelect)
    if len(selectNode.Fields) != 2 {
        t.Errorf("Expected exactly two fields and found %d", len(selectNode.Fields))
    }
}

func TestWithOneTable(t *testing.T) {
    New("select name, created_at from files")

    ast, error := AST()

    if error != nil {
        t.Errorf(error.Error())
    }

    selectNode := ast.Child.(*NodeSelect)
    if len(selectNode.Fields) != 2 {
        t.Errorf("Expected exactly two fields and found %d", len(selectNode.Fields))
    }

    if selectNode.Tables[0] != "files" {
        t.Errorf("Expected table 'files', found %s", selectNode.Tables[0])
    }
}

func TestErrorWithUnexpectedComma(t *testing.T) {
    New("select name, from files")

    _, error := AST()

    if error == nil {
        t.Errorf("Expected error 'Unexpected T_COMMA'")
    }
}

func TestErrorWithInvalidRootNode(t *testing.T) {
    New("name from files")

    _, error := AST()
    if error == nil {
        t.Errorf("Expected error 'EXPECTED T_SELECT'")
    }

}

func TestErrorSqlWithoutTable(t *testing.T) {
    New("select name from ")

    _, error := AST()
    if error == nil {
        t.Errorf("Expected error 'EXPECTED table'")   
    }
}

func TestWithLimit(t *testing.T) {
    New("select * from files limit 5")

    ast, error := AST()
    if error != nil {
        t.Errorf(error.Error())
    }

    selectNode := ast.Child.(*NodeSelect)
    if selectNode.Limit != 5 {
        t.Errorf("Limit should be 5, found %d!!!", selectNode.Limit)   
    }
}

func TestWithEmptyLimit(t *testing.T) {
    New("select * from files limit")

    _, error := AST()
    if error == nil {
        t.Errorf("Shoud throw error because limit has not value")
    }
}

func TestWithNonNumericLimit(t *testing.T) {
    New("select * from commits limit cloud")

    _, error := AST() 
    if error == nil {
        t.Errorf("Shoud throw error because limit is not a number")
    }
}

func TestWithWhereSimpleEqualComparation(t *testing.T) {
    New("select * from commits where hash = 'e69de29bb2d1d6434b8b29ae775ad8c2e48c5391' ")

    ast, err := AST() 
    if err != nil {
        t.Errorf(err.Error())
    }

    selectNode := ast.Child.(*NodeSelect)
    w := selectNode.Where
    if w == nil{
        t.Errorf("should has where node")
    }

    if reflect.TypeOf(w) != reflect.TypeOf(new(NodeEqual)) {
        t.Errorf("should be a NodeEqual")
    }

    lValue := w.LeftValue().(*NodeLiteral)
    rValue := w.RightValue().(*NodeLiteral)
    if (lValue.Value() != "hash") {
        t.Errorf("LValue should be 'hash'")
    }

    if (rValue.Value() != "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391") {
        t.Errorf("LValue should be 'e69de29bb2d1d6434b8b29ae775ad8c2e48c5391'")   
    }

}

func TestWhereWithNotEqualCompare(t *testing.T) {
    New("select * from commits where author != 'cloudson'")

    ast, err := AST() 
    if err != nil {
        t.Errorf(err.Error())
        return
    }

    selectNode := ast.Child.(*NodeSelect)
    w := selectNode.Where
    if w == nil{
        t.Errorf("should has where node")
    }

    if reflect.TypeOf(w) != reflect.TypeOf(new(NodeNotEqual)) {
        t.Errorf("should be a NodeNotEqual")
    }

    lValue := w.LeftValue().(*NodeLiteral)
    rValue := w.RightValue().(*NodeLiteral)
    if (lValue.Value() != "author") {
        t.Errorf("LValue should be 'author'")
    }

    if (rValue.Value() != "cloudson") {
        t.Errorf("LValue should be 'cloudson'")   
    }
}

func TestWhereWithGreater(t *testing.T) {
    New("select * from commits where date > '2014-05-12 00:00:00' ")

    ast, err := AST() 
    if err != nil {
        t.Errorf(err.Error())
        return
    }

    selectNode := ast.Child.(*NodeSelect)
    w := selectNode.Where

    if reflect.TypeOf(w) != reflect.TypeOf(new(NodeGreater)) {
        t.Errorf("should be a NodeGreater")
    }
}

func TestWhereWithSmaller(t *testing.T) {
    New("select * from commits where date <= '2014-05-12 00:00:00' ")

    ast, err := AST() 
    if err != nil {
        t.Errorf(err.Error())
        return
    }

    selectNode := ast.Child.(*NodeSelect)
    w := selectNode.Where

    if reflect.TypeOf(w) != reflect.TypeOf(new(NodeSmaller)) {
        t.Errorf("should be a NodeSmaller")
    } 
}

func TestWhereWithNumeric(t *testing.T) {
    New("select * from anything where blah >= 8")

    ast, err := AST()
    if err != nil {
        t.Errorf(err.Error())
        return
    }   

    selectNode := ast.Child.(*NodeSelect)
    w := selectNode.Where

    if reflect.TypeOf(w) != reflect.TypeOf(new(NodeGreater)) {
        t.Errorf("should be a NodeGreater")
    }

    greater := w.(*NodeGreater)
    if !greater.Equal {
        t.Errorf("Greater shoud be marked with equal")
    }
}

func TestWhereWithOR(t *testing.T) {
    New("select * from commits where hash = 'e69de29' or date > 'now' ")

    ast, err := AST()
    if err != nil {
        t.Fatalf(err.Error())
    }

    selectNode := ast.Child.(*NodeSelect)
    w := selectNode.Where 

    if reflect.TypeOf(w) != reflect.TypeOf(new(NodeOr)) {
        t.Errorf("should be a NodeOr")
    }

    lValue := w.LeftValue().(*NodeEqual)
    rValue := w.RightValue().(*NodeGreater)

    if reflect.TypeOf(lValue) != reflect.TypeOf(new(NodeEqual)) {
        t.Errorf("should be a NodeEqual")
    }

    if reflect.TypeOf(rValue) != reflect.TypeOf(new(NodeGreater)) {
        t.Errorf("should be a NodeGreater")
    }    
}

func TestErrorIfOperatorBeforeEOF(t *testing.T) {
    New("select * from commits where hash = 'e69de29' and date > ")
    _, err := AST()
    if err == nil {
        t.Fatalf("Shoud be error with T_GREATER before T_EOF")
    }


    New("select * from commits where hash = 'e69de29' and   ")
    _, err = AST()
    if err == nil {
        t.Fatalf("Shoud be error with T_AND before T_EOF")
    }

    // @todo should this sql throws an error ? 
    New("select * from commits where hash = 'e69de29' and date  ")
}

func TestConditionWithNoPrecedentParent(t *testing.T) {
    New("select * from commits where hash = 'e69de29' and date > 'now' or hash = 'fff3331'")

    ast, err := AST()
    if err != nil {
        t.Fatalf(err.Error())
    }

    selectNode := ast.Child.(*NodeSelect)
    w := selectNode.Where 
    if reflect.TypeOf(w) != reflect.TypeOf(new(NodeOr)) {
        t.Errorf("should be a NodeOr")
    }
    lValue := w.LeftValue().(*NodeAnd)
    rValue := w.RightValue().(*NodeEqual)

    if reflect.TypeOf(lValue) != reflect.TypeOf(new(NodeAnd)) {
        t.Errorf("should be a NodeAnd")
    }

    if reflect.TypeOf(rValue) != reflect.TypeOf(new(NodeEqual)) {
        t.Errorf("should be a NodeEqual")
    }     
}

func TestConditionWithPrecedentParent(t *testing.T) {
    New("select * from commits where hash = 'e69de29' and (date > 'now' or hash = 'fff3331')")

    ast, err := AST()
    if err != nil {
        t.Fatalf(err.Error())
    }

    selectNode := ast.Child.(*NodeSelect)
    w := selectNode.Where 
    if reflect.TypeOf(w) != reflect.TypeOf(new(NodeAnd)) {
        t.Errorf("should be a NodeAnd")
    }
    lValue := w.LeftValue().(*NodeEqual)
    rValue := w.RightValue().(*NodeOr)

    if reflect.TypeOf(lValue) != reflect.TypeOf(new(NodeEqual)) {
        t.Errorf("should be a NodeEqual")
    }

    if reflect.TypeOf(rValue) != reflect.TypeOf(new(NodeOr)) {
        t.Errorf("should be a NodeOr")
    }     
}