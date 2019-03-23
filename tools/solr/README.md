# OER Collector Apache Solr Readme

## Example queries:


All program entries with "title:Tagesschau"

```
http://localhost:8983/solr/oer-server/select?fq=title%3Atagesschau&q=*%3A*&rows=0
```


Most used words in description fields via simple facet search and date range on **start_date_time**:

```
http://localhost:8983/solr/oer-server/select?facet.field=description&facet.mincount=1&facet=on&fq=start_date_time:{2017-01-01T00%3A00%3A00.0Z%20TO%202018-01-01T00%3A00%3A00.0Z}&q=*%3A*&rows=0
```