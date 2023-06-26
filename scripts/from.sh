#!/bin/bash 

Host=127.0.0.1
Port=3306
User='root'
DB=exer

mysqldump -u$User -p --host=$Host --port=$Port --databases $DB > ./sqls/$DB.sql
echo "success!"
exit 
