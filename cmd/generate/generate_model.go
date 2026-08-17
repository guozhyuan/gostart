package main

import (
	"gostart/internal/config"
	"gostart/internal/pkg"

	"gorm.io/gen"
)

func main() {
	generateModel()
}

// 根据表生成model
func generateModel() {

	// GenerateAllTable() → 只生成模型结构体
	// ApplyBasic() → 生成指定表的查询文件（CRUD、gen.go 等）并应用到指定的模型实例

	config.ReadConfigByViper()
	pkg.ConnectDB()

	g := gen.NewGenerator(gen.Config{
		OutPath: "internal/dao/query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithGeneric, // generate mode
	})
	g.UseDB(pkg.DB)
	// g.GenerateModel("user")
	// g.GenerateAllTable()
	g.ApplyBasic(g.GenerateAllTable()...)
	// g.ApplyBasic(g.GenerateModel("platform"), g.GenerateModel("streamer_online"))
	// g.ApplyBasic(g.GenerateModel("user"))
	g.Execute()
}

// 直接传入模型实例
// g.ApplyBasic(&model.User{}, &model.Article{})
// // 或配合 GenerateModel 动态生成
// user := g.GenerateModel("user")
// article := g.GenerateModel("article")
// g.ApplyBasic(user, article)
// // 或生成所有表
// g.ApplyBasic(g.GenerateAllTable()...)

// 使用
// query.SetDefault(db)
