package main

import (
	"gostart/internal/config"
	"gostart/internal/model"
	"gostart/internal/service"

	"gorm.io/gorm"
)

// 常见的gorm标签

// primaryKey			指定字段为主键。
// autoIncrement		指定字段为自增。
// unique / uniqueIndex	创建唯一索引，确保值唯一。
// index			  	创建普通索引，加速查询。
// not null			  	设置字典非空约束。
// default:			  	设置字典默认值，如 default:0。
// autoCreateTime / autoUpdateTime	分别用于记录创建和更新时间。
// type:			  	手动指定数据库中的字段类型，如 type:varchar(100)。
//
//

//Go 与 MySQL 类型映射速查表

// int, uint		INT / INT UNSIGNED	int 有符号，uint 无符号。默认长度是 11。
// int8, uint8		TINYINT / TINYINT UNSIGNED	常用于表示小整数，或配合 bool 使用。
// int32, uint32	INT / INT UNSIGNED	32位整数。
// int64, uint64	BIGINT / BIGINT UNSIGNED	64位大整数，适合存储如 雪花ID 等字段。
// float32			FLOAT	单精度浮点数。
// float64			DOUBLE	双精度浮点数，精度更高。
// string			VARCHAR(255)	默认长度 255 的变长字符串。可以用 gorm:"size:100" 指定长度。
// bool	TINYINT(1)	MySQL 没有布尔类型，用 0/1 表示真假。
// []byte			BLOB	存储二进制数据，如图片或文件。也可用 gorm:"type:mediumblob" 指定更精确的类型。
// time.Time		DATETIME	存储日期和时间。
// *time.Time		DATETIME NULL	使用指针类型，对应数据库中可以存储 NULL 的时间字段。
// gorm:"type:decimal(10,2)"	DECIMAL(10,2)	通过 type 标签指定，适合存储金额等需要精确计算的数字。
// gorm:"type:text"	TEXT	通过 type 标签指定，用于长文本内容。
// gorm:"type:json"	JSON	通过 type 标签指定，MySQL 5.7+ 支持
func main() {
	DB.AutoMigrate(&model.UserDO{})
}

var DB *gorm.DB

func init() {
	config.ReadConfigByViper()
	service.ConnectDB()
	DB = service.DB
}
