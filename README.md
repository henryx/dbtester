# Database tester

Benchmark are:
 - Load JSON
 - Write JSON data into table
 - Extract record using JSON field

To generate data, you need to download complete records of [Open Library data dumps](https://openlibrary.org/data/ol_cdump_latest.txt.gz) and extract first 10000000 records.
To convert it in usable JSON format you need to use this command:

 * Powershell
 ```
 Get-Content -First 10000000 .\ol_cdump_latest.txt |ForEach-Object {$split = $_ -split "\t"; $content = $split[4]; Add-Content -Path .\ol_cdump_latest.json -Value $content }
 ```

* Bash
```
head -n 10000000 ol_cdump_latest.txt | cut -f 5 > ol_cdump_latest.json
```

Database tested are:
 * Elasticsearch
 * MongoDB
 * MySQL (not MariaDB)
 * PostgresQL

 License is MIT