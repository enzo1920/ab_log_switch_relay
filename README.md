# ab_log_inserter

##Create db and table

```
 CREATE DATABASE ab_log_db;

CREATE TABLE light (
    light_id serial not null primary key,
    light_val float  NOT NULL,
    light_date  timestamp default NULL
);
```
## Crontab

SHELL=/bin/sh
PATH=/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin
# m h  dom mon dow   command
*/2 * * * *  /home/ab_log_inserter/bin/abloginsert
*/4 * * * *  /home/ab_log_dht22/bin/ablogdht22
*/4 * * * *  /home/tars_dht22/bin/tarsdht22
*/10 * * * *  /home/ab_log_bmx280/bin/ablogbmx280
0,15,30,45  17,18,19,20,21,22,23,0,1,2,3,4,5,6,7,8 * * *  /home/ab_log_parserelay/bin/ab_log_parsestate
0,15,30,45  18,19,20,21,22,23,0,1,2,3,4,5,6,7,8 * * * /home/ab_log_switch_relay/bin/ablogswitch_ON
34 8 * * *  /usr/bin/python3.4 /home/TARS/check_cmd_local.py

