package main

import (
	"github.com/go-resty/resty/v2"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})

	// 初始化配置管理器
	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Info("Error reading config file:", err)
	} else {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}
}

func main() {
	client := resty.New()

	const userId int64 = 1

	dsn := "root:123456@tcp(127.0.0.1:3306)/news?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database: %w", err)
	}

	//url := "https://mp.weixin.qq.com/s/RKH4uwmlvI4A4QKBsmsQiQ"
	//url := "https://mp.weixin.qq.com/s/pvVf7KUpUVCAkFVU1oDRVQ"
	url := "https://mp.weixin.qq.com/s/vwRW6Y-ID6d7nqPo3srUkA"
	resp, err := client.R().
		Get(url)
	if err != nil {
		log.Fatal("连接URL异常: %w", err)
	}

	policy := bluemonday.StrictPolicy()
	content := policy.Sanitize(resp.String())
	content = formatText(content)

	var p string

	var user User
	result := db.Where("id = ?", userId).First(&user)
	if result.RowsAffected != 0 || len(user.Persona) > 0 {
		p = "我之前的画像是：" + user.Persona + "。帮我阅读下面这篇文章，并以此分析我的阅读兴趣画像，且要和上面所说已经有的画像进行合并：" + content
	} else {
		p = "帮我阅读下面这篇文章，并以此分析我的阅读兴趣画像：" + content
	}

	newPersona := aiReq(client, &p)
	newPersona = formatText(newPersona)
	log.Debug(newPersona)
	user.Persona = newPersona

	db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"persona"}),
	}).Create(&user)
}

func formatText(content string) string {
	content = strings.ReplaceAll(content, " ", "")
	content = strings.ReplaceAll(content, "\n", "")
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, " ", "")
	return content
}
