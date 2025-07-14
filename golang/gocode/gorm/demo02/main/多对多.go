package main

import (
	"gorm.io/gorm"
	"time"
)

type Tag struct {
	ID   uint
	Name string
}

type Article struct {
	ID    uint
	Title string
	Tags  []Tag `gorm:"many2many:article_tags;"`
}

type ArticleTags struct {
	ArticleID uint      `gorm:"primaryKey"`
	TagID     uint      `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate 实现hook函数，在create之后自动给CreatedAt赋当前时间
func (at *ArticleTags) BeforeCreate(tx *gorm.DB) (err error) {
	at.CreatedAt = time.Now()
	return err
}

func main() {
	// 设置 Article 的 Tag 表为  ArticleTags
	//DB.SetupJoinTable(&Article{}, "Tags", &ArticleTags{})   // 添加和更新表数据得用 SetupJoinTable，查询不用
	//DB.AutoMigrate(&Article{}, &Tag{}, &ArticleTags{})

	//err := DB.Debug().AutoMigrate(&Article{}, &Tag{})
	//if err != nil {
	//	panic(err)
	//}

	//// 添加文章，关联已有标签
	//DB.SetupJoinTable(&Article{}, "Tags", &ArticleTags{})
	//
	//var tags []Tag
	//DB.Find(&tags, []int{1, 2})
	//
	//DB.Debug().Create(&Article{
	//	Title: "django基础课程",
	//	Tags:  tags,
	//})

	//// 添加文章，并创建标签
	//DB.Debug().Create(&Article{
	//	Title: "python基础课程",
	//	Tags: []Tag{
	//		{Name: "python"},
	//		{Name: "基础课程"},
	//		{Name: "编程"},
	//	},
	//})

	////添加文章，选择已有标签
	//var tags []Tag
	//DB.Find(&tags, "name in ?", []string{"基础课程", "编程"})
	//DB.Create(&Article{
	//	Title: "golang基础",
	//	Tags:  tags,
	//})

	// 添加文章，选择已有的标签并创建新的表中不存在的标签

	//// 多对多的更新
	////	先删除原有的
	//var article Article
	//DB.Preload("Tags").Take(&article, 5)
	//DB.Debug().Model(&article).Association("Tags").Delete(article.Tags)
	//// 再添加新的
	//var tag Tag
	//DB.Take(&tag, 6)
	//DB.Model(&article).Association("Tags").Append(&tag)

	//var article Article
	//var tags []Tag
	//DB.Find(&tags, []int{6, 7})
	//
	//DB.Preload("Tags").Take(&article, 5)
	//DB.Model(&article).Association("Tags").Replace(tags)

}
