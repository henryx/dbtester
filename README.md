# Database tester

Benchmark tests are:
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

Hardware used for test is:
 - Intel NUC i3 7th gen with 16Gb of RAM and Crucial SSD MX500
 - Hypervisor: Windows Server 2019 Standard with Hyper-V role enabled and Windows Defender active
 - Guest: CentOS 8.2 with 2 cores and 4Gb of RAM. SELinux and firewallD are disabled

Results:

| Database      | Version | Load     | Count w/o index | Find w/o index | Index | Find with index |
|---------------|-------: |--------: |---------------: |--------------: |-----: |---------------: |
| PostgreSQL    |    12.4 |   43m32s |             28s |            15s | 1m27s |              7s |
| MongoDB       |   4.4.0 |   58m48s |             17s |            19s |   37s |              5s |
| MySQL         |  8.0.17 | 1h27m29s |             15s |            19s | 1m18s |             11s |
| Elasticsearch |   7.9.0 |   41m33s |              1s |          917ms |  97ms |            22ms | 

License is MIT