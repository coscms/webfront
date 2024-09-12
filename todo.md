# 迁移说明


执行以下正则替换
## 替换 1
```
github\.com/admpub/webx/application/(library|cmd|dbschema|initialize|listener|middleware|model|registry|request|response|transform|version)(/|")
```
替换为
```
github.com/coscms/webfront/$1$2
```
