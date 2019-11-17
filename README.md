What I did here was to repurpose an older tech challenge, which was a full crud with tests and a mongo database - you can find it on master).

Sqlite proved useful instead of mongo, since sqlite can also be run in memory for realistic testing. 

A translation layer (importer) between the json format from the gists and relational sqlite was also needed and it kinda works (m:n relations with categories are still @todo so here's why the app is not finished).

For the endpoint I'm using gin. As an orm I chose gorm, but since I'm not familiar enough with it I got a little stuck with the many to many relations.

Everything has automatically generated uuids as primary keys, so employees could each receive their own individual links over mail.

I chose to [translate](https://github.com/florinutz/form3-api/tree/everphone/importer) your gist data into a relational format because that would have easily facilitated smooth dialogue about any product issue (like how to choose the gifts).

### Usage
```bash
make binary
./bin
```
