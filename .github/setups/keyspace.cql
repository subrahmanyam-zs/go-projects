CREATE KEYSPACE IF NOT EXISTS test WITH REPLICATION = {'class':'SimpleStrategy','replication_factor' : 1};
use test;
CREATE TABLE IF NOT EXISTS persons (id int PRIMARY KEY, name text, age int, state text );