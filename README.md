#前端开发接口模拟小工具（hh_mock_api）

##工具诞生原因：
```

1、并行开发时，前期接口定义好，后面开发阶段中，前端同学不依赖，不堵塞。
2、目前调试采用电脑直连方式，那么在开发阶段中，会给前端页面造成阻断性问题，无法进行。
3、过程具有了流程化，预先定义 -> 各自开发 -> 接入调试 -> 提交测试。流程明确，同时也避免了很多问题的发生及多余的沟通。
4、前端可以自己伪造数据了。
```

##为什么要这么干？
```

1、就如前期先设计数据表是一样的，接口前期定义好，前后端同学也可从中发现更多的细节问题，以及梳理业务流程及之间的关联关系，后面开发思路会更加明确。
2、可提高开发同学对产品及技术的逻辑思维，对发现问题的能力有很大的帮助。
```

##为什么不用现有的mock-leason工具？
```

目前看来还没有人使用它（mock-leason），所以我们更需要的是简单、易用、直接、易懂，适合我们自己本身的工具。所以它来了。
```
##api文件按项目存放在各仓库
```
请求接口：/api/user/getuserinfo
api文件存在路径规则：
hh_auth_api：项目名称
	api：--|
	     user   --|
		   getuserinfo.json
	     --|				
	     account--|
		    add.json
		    del.json
```
##json文件定义规则：
```
{
    "method": "get",
    "parameters": {
        "id": "required",
        "name": "required"
    },
    "success_result": {
        "code": 200,
        "message": "success",
        "data": {
            "name": "xinyuan"
        }
    }
}
```
##配置文件：main.conf
```
#监听端口
listen = ":8081"
#api存放路径
api_path_dir = "/home/src/www/hh_test_api"

## 备注默认与执行程序在同一目录下 （可选参数 conf "/www/miss.conf" ）
```

##一键启动，即可调试～
``
$./hh_mock/hh_mock_api
``