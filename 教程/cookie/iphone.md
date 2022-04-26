### iPhone获取`Cookie`教程

安装[Stream](https://blog.csdn.net/qq_36502272/article/details/117341718)(免费)

1. Stream安装完成后，点击设置下的`HTTPS抓包`，根据提示安装CA证书

   ![](../images/stream/HTTPS抓包.jpg)

2. 安装证书后点击开始抓包

   ![](../images/stream/开始抓包.jpg)

3. 打开微信-叮咚买菜小程序-登录-购物车，经过这些步骤后，点击stream首页`抓包历史`

   ![](../images/stream/抓包日期.jpg)

4. 点击最新的抓包日期显示全部请求

   ![](../images/stream/抓包历史.jpg)

5. 点击全部请求页面右上角搜索，如果是旧版本例如`2.83.x`,可输入`cart/index`后确定。如果是新版本`2.85.x`，可以点击微信叮咚小程序`我的`-`收获地址`，搜索`user/address`(亲测有效)
   。或者根据以下API的关键字搜索，如`api.ddxq.mobi`, `user/detail`, `user/address`, `allCheck`
   , `getMultiReserveTime`等
   ```
   获取用户信息: https://sunquan.api.ddxq.mobi/api/v1/user/detail/
   获取用户买菜地址: https://sunquan.api.ddxq.mobi/api/v1/user/address/    
   勾选购物车所有商品地址: https://maicai.api.ddxq.mobi/cart/allCheck
   获取购物车商品地址: https://maicai.api.ddxq.mobi/cart/index
   预约送达时间地址: https://maicai.api.ddxq.mobi/order/getMultiReserveTime
   获取确认订单地址: https://maicai.api.ddxq.mobi/order/checkOrder
   提交订单地址: https://maicai.api.ddxq.mobi/order/addNewOrder
   ```

   ![](../images/stream/搜索.jpg)

6. 即可找到对应的请求地址，点击进去详情，在请求头部中查找参数`Cookie`

   ![](../images/stream/用户参数.jpg)
