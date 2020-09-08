package lib

import (
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

//IQiuniu .
type IQiuniu interface {
	RemoveFile(keys []string) error
}

//Qiuniuer .
type Qiuniuer struct {
}

var qiuniu IQiuniu

//Qiuniu .
func Qiuniu() IQiuniu {
	if qiuniu == nil {
		qiuniu = &Qiuniuer{}
	}
	return qiuniu
}

//SetQiuniu .
func SetQiuniu(q IQiuniu) {
	qiuniu = q
}

//RemoveFile 删除七牛云上的文件
func (qiuniu *Qiuniuer) RemoveFile(keys []string) error {
	//七牛云
	conf := Conf()
	bucket := conf.GetConf("qiniu", "bucket")
	mac := qbox.NewMac(conf.GetConf("qiniu", "ak"), conf.GetConf("qiniu", "sk"))
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	for _, key := range keys {
		err := bucketManager.Delete(bucket, key)
		if err != nil {
			return err
		}
	}
	return nil
}
