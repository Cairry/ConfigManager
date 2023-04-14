package services

import (
	"ConfigManager/db"
	"ConfigManager/kube"
	"ConfigManager/models"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	options2 "go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func (c ConfigMap) ListConfig(configName string, pageSize, pageNum int) (interface{}, error) {

	// 页面大小
	if pageSize < 1 {
		pageSize = 1
	}

	// 页面数
	if pageNum < 1 {
		pageNum = 1
	}

	// 定义分页
	options := options2.FindOptions{}
	options.SetSkip(int64(pageSize * (pageNum - 1)))
	options.SetLimit(int64(pageSize))
	options.SetSort(bson.D{})

	config := bson.D{}
	if configName != "" {
		config = bson.D{{"configname", configName}}
	}

	// 根据分页查询数据
	cur, err := db.ConfigCollection.Find(db.GvaMongoCtx, config, &options)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 定义切片，用于收集多组数据
	var comments []*models.ConfigureStruct

	// 遍历结果
	for cur.Next(db.GvaMongoCtx) {

		// 将查询结果写入自定义结构体
		comment := &models.ConfigureStruct{}
		if err := cur.Decode(comment); err != nil {
			log.Println(err)
			return nil, err
		}
		// 结构体内容加入切片中
		comments = append(comments, comment)

	}

	return comments, err
}

/*
	GetConfig
	获取指定的配置, 判断配置是否存在
*/
func (c ConfigMap) GetConfig(configName string) (models.ConfigureStruct, error) {

	var req models.ConfigureStruct
	config := bson.D{{"configname", configName}}
	err := db.ConfigCollection.FindOne(db.GvaMongoCtx, config).Decode(&req)
	if err != nil {
		return req, err
	}

	return req, nil

}

func (c ConfigMap) CreateConfig(req models.ConfigureStruct, isCreateCm string) error {

	var (
		err error
	)
	resp, _ := c.GetConfig(req.ConfigName)
	if resp.ConfigName == req.ConfigName {
		return errors.New("配置名称已存在")
	}

	_, err = db.ConfigCollection.InsertOne(db.GvaMongoCtx, models.ConfigureStruct{
		ConfigName: req.ConfigName,
		Docs:       req.Docs,
		Version:    "latest",
		Content:    req.Content,
		Created:    time.Time{},
		Deleted:    time.Time{},
	})
	if err != nil {
		return err
	}

	if isCreateCm == "true" {
		if req.Content == "" {
			log.Println("content 为空")
			return errors.New("err")
		}
		_, err = c.CreateConfigMap(req.Content)
		if err != nil {
			_ = c.DeleteConfig(req.ConfigName, "")
			return err
		}

	}

	return nil

}

func (c ConfigMap) UpdateConfig(newReq *models.ConfigureStruct, oldReq *models.ConfigureStruct, isUpdateCm string) error {

	_, err := db.HisConfigCollection.InsertOne(db.GvaMongoCtx, oldReq)
	if err != nil {
		log.Println("历史配置创建失败 ->", err)
		return err
	}

	options := bson.M{"$set": bson.M{
		"docs":    newReq.Docs,
		"content": newReq.Content,
		"created": time.Now(),
	}}

	_, err = db.ConfigCollection.UpdateMany(db.GvaMongoCtx, bson.M{"configname": oldReq.ConfigName}, options)
	if err != nil {
		log.Println("配置更改失败 ->", err)
		return err
	}

	fmt.Println("---", isUpdateCm)
	if isUpdateCm == "true" {
		_, err = c.UpdateConfigMap(newReq.Content)
		if err != nil {
			return err
		}
	}

	return nil

}

func (c ConfigMap) DeleteConfig(configName, configNamespace string) error {

	if _, err := c.GetConfig(configName); err != nil {
		return err
	}

	config := bson.D{{"configname", configName}}
	_, err := db.ConfigCollection.DeleteOne(db.GvaMongoCtx, config)
	if err != nil {
		return err
	}

	err = c.DeleteConfigMap(configName, configNamespace)
	if err != nil {
		return err
	}

	return nil

}

func (c ConfigMap) HistoryConfig(configName string, pageSize, pageNum int) (interface{}, error) {

	if _, err := c.GetConfig(configName); err != nil {
		return nil, err
	}

	// 定义分页
	options := options2.FindOptions{}
	options.SetSkip(int64(pageSize * (pageNum - 1)))
	options.SetLimit(int64(pageSize))
	options.SetSort(bson.D{})

	config := bson.D{{"configname", configName}}
	resp, _ := db.HisConfigCollection.Find(db.GvaMongoCtx, config, &options)

	// 定义切片，用于收集多组数据
	var historys []*models.ConfigureStruct

	for resp.Next(db.GvaMongoCtx) {
		var history = &models.ConfigureStruct{}
		err := resp.Decode(&history)
		if err != nil {
			return nil, err
		}
		historys = append(historys, history)
	}

	return historys, nil

}

func (c ConfigMap) RollbackConfig(configName, version string, oldConfigMap *models.ConfigureStruct) (info interface{}, err error) {

	options := bson.D{{"version", version}, {"configname", configName}}
	resp, err := db.HisConfigCollection.Find(kube.KubeCtx, options)
	if err != nil {
		log.Println(err)
		return
	}

	var comment *models.ConfigureStruct
	for resp.Next(db.GvaMongoCtx) {
		if err = resp.Decode(&comment); err != nil {
			log.Println(err)
			return nil, err
		}
	}

	err = c.UpdateConfig(comment, oldConfigMap, "true")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return comment, nil

}
