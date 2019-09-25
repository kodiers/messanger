-- create new database
DROP DATABASE IF EXISTS {{.DbName}};
CREATE DATABASE {{.DbName}} OWNER {{.DbUser}} IF NOT EXISTS;