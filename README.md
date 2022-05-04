# 叮咚抢菜助手(dingdong-grabber)

- 支持多策略抢菜
- 支持多种运行方式
- 支持ios原生API请求

# 更新

目前项目仍在调试阶段，更新了部分代码，此项目目前仅供参考，对ios设备可用，目前暂不支持android，后续可能会更新， 也可能会静默。

1. 支持ios原生api请求
2. 更新签名算法
3. 添加测试模式

# 05-05 风控

- 风控越来越严格，大家尽量少运行，可以事先使用测试模式测试，如果不通，则不要运行其他策略以避免可能的风控。
- **项目可能静默，也可能不定期更新**

# 05-04 升级程序避免风控

正在升级程序避免`405 AssertError`风控

# 05-03 风控问题

叮咚升级了风控策略，很容易被风控，出现`405 AssertError`问题, 所以运行程序每天最多运行2次。 为了避免被风控，只能完全获取用户运行环境的参数动态填写，正在升级程序。

```
{
    "success":null,
    "error":"AssertError",
    "code":"405",
    "message":"",
    "msg":"",
    "data":"-405"
}
```

# 05-01 重大更新

本次更新兼容了叮咚小程序最新版本`9.50.2`，参考了[Runc2333](https://github.com/Runc2333)和[IMLR](https://github.com/IMLR)提供的签名算法，同时感谢
[longIvan](https://github.com/longIvan)和[dodobel](https://github.com/dodobel)两位童鞋的帮忙和协作。

- **需要安装`node.js`环境**: https://www.runoob.com/nodejs/nodejs-install-setup.html

# 问题issue或者需求

大家如果遇到问题或者有更合适的需求的话，可直接在Github的Issues提问或者给出建议，我会及时关注，尽量解决和满足合理的需求。

# 运行策略

本程序暂时只提供两种策略。

1. 人工策略: 程序运行即开始抢菜，此策略下程序默认出于保护只会跑2分钟，如果没有商品库存，则会立即停止
2. 定时策略: 定时抢菜，事先订好时间，叮咚默认是早上5:59:50和8:29:50开始抢菜，这种策略要避免启动过早导致用户登录信息过期。
3. 哨兵策略: 捡漏模式，长期运行捡漏可配送时间, 不错过任何叮咚可配送时间。
4. 测试策略: 测试抢菜配置是否正确，抢菜流程是否跑通。如果失败，就不要选择其他策略再跑了。

# 使用教程

`dingdong-grapper`需要用户提供`Cookie`才可运行， 所以第一步用户需要通过抓包软件抓取相关的API提取`Cookie`。相关的API:

```
获取用户信息: https://sunquan.api.ddxq.mobi/api/v1/user/detail/
获取用户买菜地址: https://sunquan.api.ddxq.mobi/api/v1/user/address/    
勾选购物车所有商品地址: https://maicai.api.ddxq.mobi/cart/allCheck
获取购物车商品地址: https://maicai.api.ddxq.mobi/cart/index
预约送达时间地址: https://maicai.api.ddxq.mobi/order/getMultiReserveTime
获取确认订单地址: https://maicai.api.ddxq.mobi/order/checkOrder
提交订单地址: https://maicai.api.ddxq.mobi/order/addNewOrder
```

## 1. 获取`Cookie`

新版本`2.85.x`改变了原来请求地址，但是获取收获地址的API并没有改变，所以无论何种客户端，可以点击微信叮咚小程序`我的`-`收获地址`， 然后在抓包软件中输入`user/address`获取`Cookie`

- [iPhone获取`Cookie`](教程/cookie/iphone.md)
- [Android获取`Cookie`](教程/cookie/android.md)
- [Mac获取`Cookie`](教程/cookie/mac.md)
- [Windows获取`Cookie`](教程/cookie/windows.md)

## 2. 填写`Cookie`

- 将`Cookie`填入`config.yaml` cookie参数中

## 3. 运行

### IDE直接运行

可以使用[Goland](https://www.jetbrains.com/go/download/#section=mac) 或者[VS Code](https://code.visualstudio.com/download)
等IDE运行。填写完用户参数后，直接运行main.go, 默认在5:59:50和08:29:50开始抢菜，长时间运行一定要注意用户登录信息过期

- 如果没有安装Golang环境，请根据[教程/安装Go环境](教程/安装Go环境)安装
- 定时策略: 默认即为定时策略
- 人工策略: 运行此策略需要在`config.yaml`修改`strategy`字段值0，此策略下程序默认出于保护只会跑2分钟，如果没有商品库存，则会立即停止。
- 哨兵策略: 运行此策略需要在`config.yaml`修改`strategy`字段值2，此策略下程序会长期运行，直到无商品库存。

### Docker运行

Docker运行隔离了对Go等其他环境的依赖，可以直接运行

- 后续将更新上传docker image

# 注意！注意！注意！

1. 一定要设置买菜地址为默认地址，否则程序无法正常工作