package oss2bdc

import (
	"log"

	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func init() {
	config := GetConfig()
	if config.SQL.Driver == "" || config.SQL.Source == "" {
		engine = nil
		log.Println("放弃数据库模块.")
		return
	}
	newOrm, err := xorm.NewEngine(config.SQL.Driver, config.SQL.Source)
	if err != nil {
		log.Fatalln("orm.NewEngine:", err)
	}
	engine = newOrm
	if _, err := engine.Exec("CREATE OR REPLACE VIEW view_team_member AS (SELECT team.name, team.root_phone, member.phone FROM team JOIN member ON team.id = member.team_id)"); err != nil {
		log.Println("engine.Exec:", err)
	}
}

// Phone2Team 按照手机号查找
func Phone2Team(phone string) string {
	if engine == nil {
		return phone
	}
	res, err := engine.Query("SELECT name, root_phone FROM view_team_member WHERE phone = ?", phone)
	if err != nil {
		log.Fatalln("engine.Query:", err)
	}
	name := string(res[0]["name"])
	rootPhone := string(res[0]["root_phone"])
	return name + "_" + rootPhone
}

// ShutDown 切断数据库
func ShutDown() {
	if engine == nil {
		return
	}
	if _, err := engine.Exec("DROP VIEW view_team_member"); err != nil {
		log.Fatalln("engine.Exec:", err)
	}
}
