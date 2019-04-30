# log4g

# 在项目如何中使用呢？


# 如果要把日志输出到文件<br>
 只要一步那就是初始化log4g<br>
 在项目启动时初始化<br>

### 添加一行以下代码即可

```javascript 1.8

log4g.Init(
    Config{
    		Stdout:true 1
    		Path: "logs", //2 
    		NameSpace:"knowing",//3 
    	},
)

//记录一般日志
log4g.Info("hjhjhjhj") ; //

//记录错误日志
log4g.Error("hjhjhjhj") ;//

result:=map[int]string{1:"a",2:"v"}

//记录一般日志并且格式化，%s 表示 会替代为字符串，%+v 表示将变量以json形式打印
log4g.InfoFormat("info %s,%+v","ss",result);
//记录错误日志并且格式化，%s 表示 会替代为字符串，%+v 表示将变量以json形式打印
log4g.ErrorFormat("error %s,%+v","ss",result)
```

1. 日志是否控制台输出 <br>
    1.1 当为`true` 时会同时在控制台输出<br>
    1.2 当为`false` 只会在文件中输出<br>
    
2. 日志文件存放的目录<br>
    2.1 如果要输出到文件 这个参数必须配置<br>
        如果只是打印到控制台则不配置，不配置些参数时如同fmt.Println功能一样,只做控制台打印

3. 日志文件所存放的空间，如果配置，会生成  `Path`+`NameSpace` 的两级目录<br>
  如上配置 会生成`logs/knowing/`目录 
  
4.与gin 框架结合
 ```js
    log4g.Init(log4g.Config{Path:"logs",Stdout:true})
    gin.DefaultWriter = log4g.InfoLog
    gin.DefaultErrorWriter = log4g.ErrorLog

```  
        
 
        
