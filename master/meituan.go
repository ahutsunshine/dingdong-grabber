package main

import (
	"fmt"
	"github.com/dingdong-grabber/meituan"
	"sync"
	"time"
)

func mtMain(waitGroup *sync.WaitGroup) {
	waitGroup.Add(1)
	defer waitGroup.Done()

	session := meituan.MeiTuanSession{}
	session.InitSession(meituan.GetUserInfo())
	for true {

		// 0-4 || 7-23 跳过
		if time.Now().Hour() > 6 || time.Now().Hour() < 5 {
			continue
		}

		// < 5.58 跳过
		if time.Now().Hour() == 5 && time.Now().Minute() < 58 {
			continue
		}

		// > 6.10 跳过
		if time.Now().Hour() == 6 && time.Now().Minute() > 30 {
			continue
		}

		time.Sleep(200 * time.Millisecond)
		fmt.Printf("########## 获取购物车中有效商品【%s】 ###########\n", time.Now().Format("15:04:05"))
		err := session.CheckCart()
		if err != nil {
			fmt.Sprint(err.Error())
			continue
		}

		fmt.Printf("########## 生成订单信息【%s】 ###########\n", time.Now().Format("15:04:05"))
		time.Sleep(100 * time.Millisecond)
		viewResult, err := session.PreView()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		fmt.Printf("########## 提交订单中【%s】 ###########\n", time.Now().Format("15:04:05"))
		time.Sleep(100 * time.Millisecond)
		err = session.Submit(viewResult)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		break
	}
	fmt.Println("抢购成功，请前往app付款！")
}
