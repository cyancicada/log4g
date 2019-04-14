# log4g

# 在项目如何中使用呢？


# 如果要把日志输出到文件<br>
 只要一步那就是初始化log4g<br>
 在项目启动时初始化<br>

### 添加一行以下代码即可

```javascript 1.8

log4g.Init(
    Config{
    		LogMode:varMode, // 1
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

1. 日志模式 提供三种日志模式 默认给参数值时为普通模式 <br>
    1.1 `varMode` 变量模式 ，如果配置此模式会在 `Path` 参数目录下生成 一个随机的字符串目录<br>
        目录会放所有的日志文件<br>
    1.2 `consoleMode` 控制台模式 ，如果配置此模式，会在控制台进行日志输出<br>
    1.3 如果以上两种模式都没有配置 ，则以最普通形式输出在`Path` 参数目录中生成日志文件<br>
    
2. 日志文件存放的目录<br>
    2.1 如果要输出到文件 这个参数必须配置<br>
        如果只是打印到控制台则不配置，不配置些参数时如同fmt.Println功能一样,只做控制台打印


3. 日志文件所存放的空间，如果配置，会生成  `Path`+`NameSpace` 的两级目录<br>
  如上配置 会生成`logs/knowing/`目录 
