package mapper

import (
	"fmt"
	"log"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/rudeigerc/broker-gateway/model"
	"github.com/spf13/viper"
)

var DB *gorm.DB

func NewDB() {
	var err error

	mysqlAddr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.dbname"),
	)

	DB, err = gorm.Open("mysql", mysqlAddr)
	if err != nil {
		log.Fatalf("[mapper.database.NewDB] [FETAL] %s", err)
	}

	DB.AutoMigrate(&model.Firm{}, &model.Order{})
}

func NewEtcdClient() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   viper.GetStringSlice("etcd.endpoints"),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatalf("[mapper.database.NewEtcdClient] [FETAL] %s", err)
	}
	return cli
}
