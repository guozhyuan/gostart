
### Docker



### Elasticsearch
#### 概念
```
Elasticsearch是一个分布式搜索引擎，用于全文搜索、分析和存储数据。它支持多种数据类型，如文本、数字、日期等，并提供丰富的查询功能，如全文搜索、过滤、排序等。Elasticsearch还支持多种分词器和分析器，可以根据不同的需求进行分词和分析。Elasticsearch还支持多种索引和文档操作，如创建索引、查询文档、删除文档、修改文档等。Elasticsearch还支持多种分片

对比于关系型数据库
Index     -- Tabel    -- 表
Document  -- Row      -- 记录
Field     -- Column   -- 字段
Mapping   -- Schema   -- 约束 

````


#### Index操作

```


创建索引库：PUT /索引库名

    mappings：定义字段类型等约束.常见属性如下

      - type：字段类型，如text、keyword、integer、boolean、date、object  
             #keyword类型只能整体搜索，不支持搜索部分内容
      - properties：该字段的子字段
      - analyzer：分词器，如standard、ik_max_word、ik_smart
      - index：是否索引，如true、false
      - store：是否存储，如true、false
      - boost：权重，如1.0
      - copy_to：复制字段，如name_copy
      - ignore_above：忽略长度超过的字段，如1024
      - ignore_missing：忽略缺失的字段，如true

    settings：设置分片数和副本数
      - number_of_shards：分片数，如3
      - number_of_replicas：副本数，如2

    示例：
    PUT /my_index
    {
        "settings": {
            "number_of_shards": 3,
            "number_of_replicas": 2
        },
        "mappings": {
            "ip":{
                type: "keyword",  # keyword类型不需要分词器
                index: true,
            },
            "address":{
                type: "text",
                analyzer: "ik_smart"
            },
            "properties": {
                "name": {
                    "type": "text"
                },
                "age": {
                    "type": "integer"
                }
            }
        }
    }

    查询索引库：GET /索引库名


    删除索引库：DELETE /索引库名


    修改索引库（添加字段）：PUT /索引库名/_mapping

    示例：
    PUT /my_index/_mapping   # 只允许添加字段，不能修改字段
    {
        "properties": {
            "email": {
                "type": "text"
            }
        }
    }

```

#### Document操作

```
Inex结构:
{
    mappings: {
        ip: {type: "keyword"},
        address: {type: "text", analyzer: "ik_smart"},
        name: {type: "text"},
        age: {type: "integer"}
    }
}

创建文档：POST /{索引库名}/_doc/文档id
    示例：
    POST /my_index/_doc/1
    {
        "ip": "192.168.1.1",
        "address": "北京市海淀区",
        "name": "张三",
        "age": 25
    }

查询文档：GET /{索引库名}/_doc/文档id
    示例：
    GET /my_index/_doc/1
    //批量查询：查询该索引库下的全部文档
    GET /my_index/_search

删除文档：DELETE /{索引库名}/_doc/文档id
    示例：
    DELETE /my_index/_doc/1

修改文档：
    全量修改：PUT /{索引库名}/_doc/文档id
    增量修改：POST /{索引库名}/_update/文档id { "doc": {字段}}

```

### Go RestAPI
```
官方文档
    https://www.elastic.co/docs/reference/elasticsearch/clients/go
Github
    https://github.com/elastic/go-elasticsearch

```


### Kibana