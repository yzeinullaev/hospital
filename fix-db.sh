#!/bin/bash
echo "Создание базы данных hospital_feedback..."
docker exec -i hospital-bot-mysql-1 mysql -u root -ppassword -e "CREATE DATABASE IF NOT EXISTS hospital_feedback;"
docker restart hospital-bot-app-1
echo "Готово!" 