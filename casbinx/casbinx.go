package casbinx

import (
	"database/sql"
	"fmt"
	sqladapter "github.com/Blank-Xu/sql-adapter"
	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
)

type Casbinx struct {
	Sqladapter     *sqladapter.Adapter
	CasbinEnforcer *casbin.Enforcer
}

func NewCasbinx(db *sql.DB, model string, policy string) *Casbinx {
	a, err := sqladapter.NewAdapter(db, "mysql", policy)
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer(model, a)
	if err != nil {
		panic(err)
	}
	if err = e.LoadPolicy(); err != nil {
		zap.L().Error(fmt.Sprintf("Casbin LoadPolicy failed, err: %s", err.Error()))
	}
	return &Casbinx{
		Sqladapter:     a,
		CasbinEnforcer: e,
	}
}
