# Database tester

Benchmark tests are:
 - Load JSON
 - Write JSON data into table
 - Extract record using JSON field

## Data

Dataset used for benchmark is taken from [Open Library data dump](https://openlibrary.org/data/ol_cdump_latest.txt.gz), using the first 10000000 records.
To convert it in usable JSON format, is needed to execute one of these command:

 * Powershell
 ```
 Get-Content -First 10000000 .\ol_cdump_latest.txt | ForEach-Object { $split = $_ -split "\t"; $content = $split[4]; Add-Content -Path .\ol_cdump_latest.json -Value $content }
 ```

* Bash
```
head -n 10000000 ol_cdump_latest.txt | cut -f 5 > ol_cdump_latest.json
```

## Tests

Database tested are:
 * Elasticsearch
 * MongoDB
 * MySQL (not MariaDB)
 * PostgresQL

### Tests on 2021-12-27

Hardware used for tests is:
 - Intel NUC i3 7th gen with 16Gb of RAM and Crucial SSD MX500
 - Hypervisor: Windows Server 2019 Standard with Hyper-V role enabled and Windows Defender active
 - Guest: CentOS 8.2 with 2 cores and 4Gb of RAM. SELinux and firewallD are disabled

Results:

| Database      | Version | Load data | Count w/o index | Find w/o index | Index | Find with index |
|---------------|-------: |----------:|---------------: |--------------: |-----: |---------------: |
| PostgreSQL    |    12.4 |    43m32s |             28s |            15s | 1m27s |              7s |
| MongoDB       |   4.4.0 |    58m48s |             17s |            19s |   37s |              5s |
| MySQL         |  8.0.17 |  1h27m29s |             15s |            19s | 1m18s |             11s |
| Elasticsearch |   7.9.0 |    41m33s |              1s |          917ms |  97ms |            22ms | 

### Tests on 2022-12-11

Hardware used for tests is:
- ASrock Deskmini A300 with AMD Ryzen 5 3400G, 32Gb of RAM and Kingston NVME SA2000M8500G
- Hypervisor: Proxmox 7.3 using LVM as storage data
- Guest: AlmaLinux 9.1 with 2 cores and 4Gb of RAM. SELinux is set on permissive and firewallD is disabled

Results:

| Database      | Version | Load data | Count w/o index | Find w/o index |      Index | Find with index |
|---------------|--------:|----------:|----------------:|---------------:|-----------:|----------------:|
| PostgreSQL    |    15.1 |    21m27s |         6s579ms |        6s137ms |   14s454ms |             2ms |
| MongoDB       |   6.0.3 |    41m28s |        12s630ms |       12s606ms |   28s036ms |            13ms | 
| MySQL         |  8.0.31 |    33m57s |         8s198ms |       25s096ms | 1m28s734ms |             4ms |
| Elasticsearch |   8.5.3 |    18m08s |           346ms |          151ms |      197ms |             7ms |

## License

License is MIT