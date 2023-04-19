# DB

The database is based on Postgres 14, deployed using the [Tanzu Postgres Operator](https://docs.vmware.com/en/VMware-SQL-with-Postgres-for-Kubernetes/index.html)

TanzuTrends are currently using one Table with the following schema

```
id BIGINT PRIMARY KEY,
time TEXT,
username TEXT,
text TEXT,
hashtags TEXT
```

Where the tweet id, is is being used as the primary key, to make sure there is no duplicate entrys.

Access to the database, is being handled by a Kubernetes Secret, that both [Frontend](frontend.md) and [Scrape](scrape.md) has mapped, to get the credentials needed to connect.
