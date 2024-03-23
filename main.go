package main

import (
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/microcosm-cc/bluemonday"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"regexp"
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

func main2() {
	client := resty.New()

	persona, err := load()
	if err != nil {
		log.Fatal(err)
	}

	if len(persona) == 0 {
		log.Fatal("没有找到用户画像")
	}

	p := "我的用户画像是：" + persona + "。以我的用户画像为基准，提炼给我可能感兴趣的2024年3月22日发生的真实新闻，要求分点，简洁"
	news := aiReq(client, &p)
	fmt.Println(news)
}

func main() {
	client := resty.New()

	//url := "https://mp.weixin.qq.com/s/RKH4uwmlvI4A4QKBsmsQiQ"
	//url := "https://mp.weixin.qq.com/s/pvVf7KUpUVCAkFVU1oDRVQ"
	//url := "https://mp.weixin.qq.com/s/vwRW6Y-ID6d7nqPo3srUkA"
	//url := "https://mp.weixin.qq.com/s/JmnPO1TR8mUX656P7rt64w"
	url := "https://top.baidu.com/board?tab=realtime"
	resp, err := client.R().
		Get(url)
	if err != nil {
		log.Fatal("连接URL异常: %w", err)
	}

	policy := bluemonday.StrictPolicy()
	content := policy.Sanitize(resp.String())
	content = formatText(content)

	var p string

	//persona, err := load()
	//if err != nil {
	//	log.Fatal(err)
	//}
	p = "帮我阅读下面这篇文章，并提炼出重点，按点列出来：" + content
	//if len(persona) > 0 {
	//	p = "我之前的画像是：" + persona + "。帮我阅读下面这篇文章，并以此分析我的阅读兴趣画像，且要和上面所说已经有的画像进行合并：" + content
	//} else {
	//	p = "帮我阅读下面这篇文章，并以此分析我的阅读兴趣画像：" + content
	//}
	log.Debug(p)

	newPersona := aiReq(client, &p)
	newPersona = formatText(newPersona)
	log.Debug(newPersona)
	//
	//err = save(newPersona)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func formatText(content string) string {
	return regexp.MustCompile(`\s+`).ReplaceAllString(content, "")
}

func save(content string) error {
	err := os.WriteFile("persona.data", []byte(content), 644)
	if err != nil {
		return errors.New("保存数据异常")
	}
	return nil
}

func load() (string, error) {
	// 不存在文件
	_, err := os.Stat("persona.data")
	if err != nil {
		return "", nil
	}

	file, err := os.Open("persona.data")
	if err != nil {
		return "", errors.New("加载数据异常")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("加载数据异常:", err)
		}
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		return "", errors.New("加载数据异常")
	}

	return string(content), nil
}
