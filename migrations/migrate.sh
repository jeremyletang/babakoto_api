
HOST="192.168.99.100"
PORT="3307"

mysql -v --host=$HOST -P $PORT -u root --password=root < 0_create_db.sql
mysql -v --host=$HOST -P $PORT -u root --password=root < 1_create_users.sql
