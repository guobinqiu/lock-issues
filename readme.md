预备知识
---
隔离级别（低到高）

- Read Uncommitted
- Read Committed (Oracle, PostgreSQL, SQLServer...)
- Repeatable Read (MySQL)
- Serializable

性能要求：级别越低越好；数据一致性要求：级别越高越好

MySQL默认级别：Repeatable Read

```
mysql> select @@tx_isolation;
+-----------------+
| @@tx_isolation  |
+-----------------+
| REPEATABLE-READ |
+-----------------+
1 row in set, 1 warning (0.01 sec)
```

先调整到：Read Committed

```
mysql> set transaction_isolation='read-committed';
Query OK, 0 rows affected (0.00 sec)

mysql> select @@tx_isolation;
+----------------+
| @@tx_isolation |
+----------------+
| READ-COMMITTED |
+----------------+
1 row in set, 1 warning (0.00 sec)
```

乐观锁
---

[正确输出](/blob/master/optimistic_lock/correct/main.go)

```
余额：1000
[Transfer]查询到存款余额: 1000
[Withdraw]查询到存款余额: 1000
[Transfer]开始事务
[Transfer]存入100，余额变成 1100
[Transfer]提交事务
[Withdraw]开始事务
[Withdraw]已被transfer转账事务修改，需要重新执行（stale object error）
[Withdraw]查询到存款余额: 1100
[Withdraw]开始事务
[Withdraw]取出100，余额变成: 1000
[Withdraw]提交事务
余额: 1000
```

[错误输出](/blob/master/optimistic_lock/incorrect/main.go)

```
余额：1000
[Transfer]查询到存款余额: 1000
[Withdraw]查询到存款余额: 1000
[Transfer]开始事务
[Transfer]存入100，余额变成 1100
[Transfer]提交事务
[Withdraw]开始事务
[Withdraw]取出100，余额变成: 900
[Withdraw]提交事务
余额: 900
```

死锁
---

[输出](/blob/master/dead_lock/main.go)

```
[TransB]开始事务
[TransB]锁住ord.id=1这条记录
[TransA]开始事务
[TransA]锁住customer.id=1这条记录
[TransB]等待TransA释放锁
[TransA]Error 1213: Deadlock found when trying to get lock; try restarting transaction
[TransA]提交事务
[TransB]提交事务
```
