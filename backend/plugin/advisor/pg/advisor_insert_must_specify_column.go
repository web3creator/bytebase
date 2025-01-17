package pg

// Framework code is generated by the generator.

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/bytebase/bytebase/backend/plugin/advisor"
	"github.com/bytebase/bytebase/backend/plugin/advisor/db"
	"github.com/bytebase/bytebase/backend/plugin/parser/sql/ast"
)

var (
	_ advisor.Advisor = (*InsertMustSpecifyColumnAdvisor)(nil)
	_ ast.Visitor     = (*insertMustSpecifyColumnChecker)(nil)
)

func init() {
	advisor.Register(db.Postgres, advisor.PostgreSQLInsertMustSpecifyColumn, &InsertMustSpecifyColumnAdvisor{})
}

// InsertMustSpecifyColumnAdvisor is the advisor checking for to enforce column specified.
type InsertMustSpecifyColumnAdvisor struct {
}

// Check checks for to enforce column specified.
func (*InsertMustSpecifyColumnAdvisor) Check(ctx advisor.Context, _ string) ([]advisor.Advice, error) {
	stmtList, ok := ctx.AST.([]ast.Node)
	if !ok {
		return nil, errors.Errorf("failed to convert to Node")
	}

	level, err := advisor.NewStatusBySQLReviewRuleLevel(ctx.Rule.Level)
	if err != nil {
		return nil, err
	}
	checker := &insertMustSpecifyColumnChecker{
		level: level,
		title: string(ctx.Rule.Type),
	}

	for _, stmt := range stmtList {
		checker.text = advisor.NormalizeStatement(stmt.Text())
		ast.Walk(checker, stmt)
	}

	if len(checker.adviceList) == 0 {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  advisor.Success,
			Code:    advisor.Ok,
			Title:   "OK",
			Content: "",
		})
	}
	return checker.adviceList, nil
}

type insertMustSpecifyColumnChecker struct {
	adviceList []advisor.Advice
	level      advisor.Status
	title      string
	text       string
}

// Visit implements ast.Visitor interface.
func (checker *insertMustSpecifyColumnChecker) Visit(in ast.Node) ast.Visitor {
	if node, ok := in.(*ast.InsertStmt); ok && len(node.ColumnList) == 0 {
		checker.adviceList = append(checker.adviceList, advisor.Advice{
			Status:  checker.level,
			Code:    advisor.InsertNotSpecifyColumn,
			Title:   checker.title,
			Content: fmt.Sprintf("The INSERT statement must specify columns but \"%s\" does not", checker.text),
			Line:    node.LastLine(),
		})
	}

	return checker
}
