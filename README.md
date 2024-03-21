# 影子文件同步服务
提供几个接口用于添加元数据，读取数据集，查询索引等功能
# 接口
## /api/v1/addItem (POST)
~~~
此方法给管理服务器调用
此方法用户向服务器添加元数据，元数据可以是一个任意结构形式的json，但在根节点必须包含"id":"xxx"的键值对，一般用于标识元数据的唯一标识，
可以用数据库的索引，或者系统生成的UUID。
在/removeItem,可以指定id来删除元数据
调用方式：
/api/v1/addItem?entry=xx&name=xx
{json}
entry表示数据集的名称，比如guest,global
name表示元数据所在分组，比如face,card等
{json}：表示元数据json，举例：
{
   "endTime":"UTC time",
   "faceFeature":"zzz",
   "id":"2",  #这个很关键
   "name":"xxx",
   "no":"000001C-01B-01U-01F-0101R",
   "startTime":"UTC time"
}
~~~

## /api/v1/removeItem (GET)
~~~
此方法给管理服务器调用
删除一个元数据，调用方法:
/api/v1/removeItem?entry=xx&name=xx&id=xx
entry:表示数据集名称
name:表示元数据所在分组
id:元数据的唯一标识，如果为空，则删除 entry/name所在的分组
~~~

## /api/v1/getIndexes (GET)
~~~
获取全部数据集的索引，格式如下：
[
   {
      "name":"guest",
      "hash":"d4fa3be3211d05b94ab46cc5b0adfcaf"
   },
   {
      "name":"0001C-01B",
      "hash":"d4fa3be3211d05b94ab46cc5b0adfcaf"
   }
]
hash为数据集的md5校验码
~~~

# /api/v1/getEntry (GET)
~~~
获取指定名称的数据集,调用方法:
/api/v1/getEntry?entry=guest
返回数据集，比如:
{
   "card":[
      {
         "endTime":"UTC time",
         "faceFeature":"zzz",
         "id":"10",
         "name":"xxx",
         "no":"000001C-01B-01U-01F-0101R",
         "startTime":"UTC time"
      }
   ],
   "face":[
      {
         "endTime":"UTC time",
         "faceFeature":"zzz",
         "id":"10",
         "name":"xxx",
         "no":"000001C-01B-01U-01F-0101R",
         "startTime":"UTC time"
      }
   ]
}
~~~

# /api/v1/removeEntry (GET)
~~~
此方法给管理服务器调用，用于删除某个数据集，调用方法:
/api/v1/removeEntry?entry=guest
~~~

# 接口返回值
~~~
200：标识成功
其他均为失败
~~~
