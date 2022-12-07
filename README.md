# dis-weather-collection
Data-Intensive Systems project.

## Prototype implementation notes

#### Databases

The prototype is implemented using three PostgreSQL databases each on their own Docker container that are run on the same network with Docker Compose. The server is accessible on the internet. 

The pg_hba.conf file was modified to allow the databases to communicate with each other without password. The postgresql.conf file was modified to enable table-specific replication between the databases and to configure the pg_cron extension. 

The overall schema generally follows the original design document, however, there are some differences in the details of the implementation.

- Each database corresponds to one country. For this implementation "db1" was chosen as Finland, "db2" as Sweden, and "db3" as Norway. 

- Instead of partitioning all tables, data fragmentation is done so that only the "organization" table is partitioned based on the "country" column. Other fragmented tables have a local table for all the local data and additional foreign tables of the same format for the data of each of the other databases. These include the "station", "device", and "measurement" tables. For combining the local and foreign tables, there is an additional view for each of these tables that queries the union of the local and foreign data. Since globally unique identifiers can't be guaranteed across the tables from all the databases, the views also have a "country" column where the country of the data is automatically inserted based on which database it is from. Using the combination of "id" and "country" then guarantees global uniqueness (this is also required for the "organization" partitioned table). 

- The reason that data fragmentation is done like this instead of partitioning all tables is that it would require knowing the "country" value of each data row directly. Since this information is only known based on the "organization" table, the column would need to be repeated for the data of all of these tables. Moreover, primary and foreign key constraints aren't supported for partitioned tables, so implementing the tables as above allows for better control over data integrity. While the "country" column is indeed repeated in the views combining the data from all the databases, since it is in a view, it isn't physically saved in the data rows. 

- The postgres_fdw extension is used for foreign tables. The naming uses suffixes "_local" for local tables, "_dbx" for foreign tables, and "_all" for the views of all data. 

- The "users" table only has a local table in each database, as this information doesn't need to be accessible from the other sites.

- The "phenomenon_type" and "unit" tables have to replicated and be exactly the same in all of the databases. This was implemented using PostgreSQL's logical replication. Since this functionality requires one of the databases to be the master, "db1" was chosen for this implementation as controlling the master data for these tables. Thus, when updating the data of these tables in "db1", the changes are automatically replicated in "db2" and "db3". These tables should consequently only be updated for "db1". 

- For performance reasons there are also material views aggregating the measurement data for each day and accessing relevant station information related to that aggregated data. There is a local "station_info" and "meas_min_max_day" material view in each of the databases. For replicating the material views globally, the "meas_min_max_day" material views of the other databases are added as foreign tables. These foreign tables are then combined with the local one using unions to create an overall materialized view for the "meas_min_max_day" data. Even though having the data of the other databases as foreign tables means that they wouldn't be saved locally, combining all the data again in another materialized view ensures that it is saved anyway. As with the normal views, the "country" value is inserted for globally unique identification. The materialized view for the "station_info" is queried from the other locally accessible views combining all data instead of foreign tables. The reason for using foreign tables for the "meas_min_max_day" data is that then only the aggregated data needs to be sent over the network and not all measurement data. 

- The pg_cron extension is used for periodically refreshing all the material views. All of the view are set to refresh every minute, set with a query like "select cron.schedule('* * * * *', $$refresh materialized view meas_min_max_day_local$$);"


#### Clients

The client desktop application includes a user interface for getting data from the system in the form of charts, and a station simulator that inputs simulated measurement data into the system based on a separate configuration file. The application is implemented with Go. 
