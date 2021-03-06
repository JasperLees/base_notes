
# 索引

### `MySQL` 为什么使用 `B+` 树来作索引，对比 `B` 树它的优点和缺点是什么？

##### `B+`树和`B`数的定义

- `B`树是一种平衡多路查找树，它的每个子叶结点的高度是一样的，不管是叶子结点还是非叶子结点，都会保存数据；
- `B+`树是`B`树的一个变种，`B+`树和`B`树的不同主要在于：
  - `B+`树中的非叶子节点不存储数据，只存索引值，仅叶子结点存储数据
  - `B+`树通过链表的形式，依次按顺序将叶子结点连接

##### 为什么`B+`树比`B`树更适合应用于数据库索引？

- `B+`树更加适应磁盘的特性，相比`B`树减少了`I/O`读写的次数
  - 由于索引文件很大，因此索引文件存储在磁盘上；
  - `B+`树的非叶子结点只存关键字不存数据，因而单个页可以存储更多的关键字；
  - 即一次性读入内存的需要查询的关键字也就越多，磁盘的随机`I/O`读取次数相对就减少了
- `B+`树利于扫库和范围查询
  - `B+`树叶子结点之间用链表有序连接，所以扫描全部数据只需扫描一遍叶子结点
  - `B`树由于非叶子结点也存数据，所以只能通过中序遍历按序来扫描。也就是说，对于范围查询和有序遍历而言，`B+`树的效率更高

##### 拓展：`B+ Tree VS ` 红黑树

- 更少的查找次数：红黑树是二叉树，导致同样的数量的`data`，红黑树的高度会大于 `B+ Tree`
- 减少磁盘寻道：因为磁盘不是严格读取数据，而是会有预读的情况。顺序读取的过程中不需要磁盘寻道，速度更快

------

### 唯一索引与普通索引的区别是什么？使用索引会有哪些优缺点？

##### 区别

- 唯一索引和普通索引的区别是唯一索引定义了唯一性
- 从查询性能上对比，假如查询`k=5`的记录：
  - 对于普通索引来说，查到满足的第一个记录后，需要查找下一条记录，直到碰到第一个不满足`k=5`的条件，查询结束
  - 对于唯一索引，查到满足的第一个记录后，即可返回，查询结束
  - `InnoDB` 对于一条记录并不是将这条记录从磁盘中读出来，而是以页为单位，将整体读取到内存中，每个数据页的大小是`16kb`，那么对于一个`int`索引来说，一个数据存储的索引数据近千条，普通索引查询下一条需要读取两个数据页的几率很小。
  - 因此唯一索引和普通索引在查询的性能上差别不大
- 从更新性能上对比：
  - 当需要更新一个数据页时，如果数据页在内存中就直接更新
  - 如果数据页不在内存中，在不影响数据一致性的情况下，`InnoDB`会将这些更新操作缓存到 `change buffer`中，这样不需要从磁盘中读取数据页，减少磁盘读取，等到下一次查询时，先将原数据页的数据查到，然后再进行`merge`
  - 唯一索引，所有更新的操作都需要先判断操作是否违反了唯一性，所以必须将数据页读入内存才能判断，不能用到 `channge buffer`
  - 因此，在更新操作上，普通索引优于唯一索引

##### 使用索引会有哪些优缺点？

- 优点
  - 创建唯一性索引，保证数据库表中该列的唯一性
  - 大大加快数据的检索速度
  - 加速表和表之间的连接
  - 在使用分组和排序子句进行数据检索时，可以显著减少查询中分组和排序的时间
  - 通过使用索引，可以在查询的过程中使用优化隐藏器，提高系统的性能，例如：`COUNT(*)`

- 缺点
  - 创建和维护索引要耗费时间，这种时间随着数量的增加而增加
  - 索引需要占物理空间，除了数据表占数据空间外，每一个索引还要占一定的物理空间
  - 当对表中的数据进行增加、删除和修改时，索引也要动态维护，降低了数据的维护速度

------

### 数据库有哪些常见索引？

##### 从数据结构角度

- `B+Tree`索引：`MySQL`默认索引，可用于查找、分组、排序
- `Hash`索引：基于哈希表实现，适合新增和等值查询，不是有序的，范围查询很慢，必须全表扫描；
- `FULLTEXT`索引:大量文件检索时，需要用全文索引，因为它的速度是 `like`的 `N`倍(`MySQL 5.6`版本`InnoDB`支持)
- `R-Tree`索引：空间数据索引会从所有维度来索引数据，主要用于地理数据存储。

##### 从物理存储角度

- 聚簇索引：叶子结点存放一整行记录
- 非聚簇索引：叶子节点存放索引+指向行数据的指针

##### 从逻辑角度

- 主键索引：是一种特殊的唯一索引，不允许有空值
- 普通索引/单列索引
- 复合索引/多列索引：多个字段上创建的索引，在查询条件中，使用第一个字段，所以才会被使用，符合最左原则
- 唯一索引：该列值不能重复

------

### 数据库索引的实现原理是什么？

- 索引是一个排序的列表，在这个列表中存储着索引的值和包含这个值的数据所在行的物理地址，在数据十分庞大的时候，索引可以大大加快查询的速度，这是因为使用索引后可以不用扫描全表来定位某行的数据，而是先通过索引表找到该行数据对应的物理地址然后访问相应的数据。

------

### 聚簇索引和非聚簇索引有什么区别？什么情况用聚集索引？

#### 区别

- 聚簇索引与非聚簇索引的区别是：叶节点是否存放行数据

##### 聚簇索引

- 聚簇索引并不是一种单独的索引类型，而是一种数据存储方式。比如：`InnoDB`的聚簇索引使用`B+Tree`的数据结构存储索引和数据
- 当表有聚簇索引时，它的数据行实际存放在索引的叶子页。因为无法同时把数据行存放在两个不同的地方，所以一个表只能由一个聚簇索引。
- 聚簇索引的二级索引：叶子节点不会保存引用的行的物理位置，而是保存行的主键位置

##### 非聚簇索引

- 表数据存储顺序与索引顺序无关，叶结点包含索引字段值及指向数据页数据行的逻辑指针，其行数量与数据表行数据量一致。

##### 聚簇索引的优点

- 可以把相关数据保存在一起
- 访问数据更快
- 使用覆盖索引扫描的查询可以直接使用页借点的主键值

##### 聚簇索引的缺点

- 聚簇索引数据最大限度地提高`IO`密集性应用的性能，但如果把数据全部放在内存中，则访问的顺序没那么重要，聚簇索引就没什么优势了
- 插入速度严重依赖插入顺序
- 更新聚簇索引列的代价很高
- 基于索引的表插入新行，或主键被更新导致需要移动行的时候，可能面临页分裂
- 聚簇索引可能导致全表扫描变慢，尤其是行比较稀疏，或者由于页分裂导致数据存储不连续

#### 使用场合

##### 聚簇索引的使用场合

- 查询命令的回传结果是以该字段为排序依据的
- 查询的结果返回一个区间值
- 查询的结果返回某值相同的大量结果集

##### 非聚簇索引的使用场合

- 查询所获数据量较少时
- 某个字段中的数据的唯一性比较高时

#### 拓展：`InnoDB`的主键列

- 如果没有定义主键，`InnoDB`会选择一个唯一的非空索引代替
- 如果没有这样的索引，`InnoDB`会隐式定义一个主键作为聚簇索引

#### 拓展：为什么主键通常建议使用自增 `id`

- 聚簇索引的数据的物理存放顺序与索引的顺序一致。即：只要索引是相邻的，那么对应的数据一定也是相邻地存放在磁盘上
- 如果主键不是自增 `id`，那么需要不断调整数据的物理地址、分页等；
- 如果是自增的，它只需要一页一页地写，索引结构相对紧凑，磁盘碎片少，效率也高
- 
------

### 索引覆盖、最左原则、索引下推

#### 简述什么是覆盖索引？

##### 概念

- 覆盖索引指的是，在使用索引查询时，需要查询的值已经在索引树上，因此可以直接提供查询结果，不需要回表。
- 由于覆盖索引可以减少树的搜索次数，显著提升查询性能，所以使用覆盖索引是一个常用的性能优化手段

##### 例子

```sql
/* id 为主键，c 为索引 */
SELECT id FROM T WHERE k BETWEEN 3 AND 5;
```

- 因为索引`k`的数据时`k+id`，包含了需要查询的值，不需要再回表，所以符合覆盖索引查询。


#### 简述什么是最左匹配原则？

##### 概念

- 最左匹配原则是指，在执行查询过程中，不需要索引的全部定义，只需要满足最左前缀，就可以利用索引来加速检测。
- 最左前缀可以是联合索引的最左`N`个字段，也可以是字符串索引的最左`M`个字符。

##### 例子

```sql
/* 联合索引，a, b, c */
SELECT * FROM T WHERE A=val1 /* A 符合最左匹配 */
SELECT * FROM T WHERE A=val1 AND B=val2 /* AB 符合最左匹配 */
SELECT * FROM T WHERE A=val1 AND C=val3 /* 该SQL会选择 A 作为索引去查询 */
SELECT * FROM T WHERE B=val1 AND C=val3 /* 没有用到索引，不符合 */

/* 假设 a 是字符串 */
SELECT * FROM T WHERE A="A%" /* 也是符合最左匹配原则的 */
```

##### 拓展

- 一个索引可以最多包含16列
- 取字符串的最左部分字符可以作为前缀索引，但是不能使用覆盖索引，因为索引查询时必须回表。


#### 假设建立联合索引 (a, b, c) 如果对字段 a 和 c 查询，会用到这个联合索引吗？

- 准确来说，查询字段`(a, c)` 的话，其中一部分字段会用到联合索引，也就是字段`a`，但`c`的部分没有使用索引。

- 从定义来说，我们说查询用到联合索引的话，一般指所有字段的查询都用到了该索引。也就是 `(a), (a,b),(a,b,c)` 这三种情况

  
#### 简述什么是索引下推？

##### 概念

索引下推是 `MySQL5.6` 引入的特性，指在索引遍历的过程中，在符合最左匹配原则情况下，对索引中包含的字段先进行判断，直接过来掉不满足条件的记录，减少回表次数。

##### 例子

```sql
/* 联合索引 (name, age) */
SELECT * FROM T WHERE name like 'A%' AND age = 10
```

- 在无索引下推的情况下，找到符合 `like A%` 的行，直接回表查主键的行记录，然后在判断 `age=10`
- 在有索引下推的情况下，先判断 `age = 10，再回表。

------

### `MySQL` 的索引什么情况下会失效？

- 联合索引，没有使用最左匹配原则
  -例如：联合索引`(a, b)`，`SELECT * FROM T WHERE b = 1`
  - 因为索引都是先基于`a`排序的，`a`相同的情况，再基于`b`排序，所以查询条件没有`a`，对于索引来说，`b`是无序的，所以索引失效；

- 范围查询，右边失效
  -  例如：联合索引`(a, b)`, `SELECT * FROM T WHERE a > 1 AND b = 1`
  -  字段`a` 在 `B+ Tree`上是有序的，可以用二分查找去定位，`a`的索引能被用到，`b`有序的前提是 `a`是确定值，在之前逻辑中，`b`是无序的，所以`b`用不到索引

- `like`的 `%` 放在左边，索引失效
  - 例如：`SELECT * FROM T WHERE a LIKE "%a"`
  - 字符串的排序方式是，先按照第一个字母排序，如果第一个字母相同，就按照第二个字母排序，依次类推

- 条件字段做函数操作，导致索引失效
  - 对索引字段做函数操作，可能会破坏索引的有序性，优化器一刀切，不管有没有破坏有序性，都不会考虑使用索引
  
- 隐式转换，导致索引失效
  - 数据类型转换，如果字符串和数字比较，则将字符串转换成数字
  - 如果条件条件字段是字符串，查询条件是数字，则会隐式的将字符串转为数字，即对字段做函数操作

------

# 事务

### 什么是数据库事务，简述数据库中的 `ACID` 分别是什么？

##### 事务的概念

- 数据库事务是指一个逻辑单元执行的一系列操作，一个逻辑工作单元必须有四个属性，称为 `ACID`（原子性、一致性、隔离性和持久性）属性

##### `ACID`

- `ACID`是指数据库在写入或更新资料的过程中，为了保证事务时正确可靠的，所必须具备的四个特性：`A`代表原子性，`C`代表一致性，`I`代表隔离性，`D`代表持久性
- `A`：(`Atomicity`)原子性
  - 一个事务中的所有操作，要么全部完成，要么全部不完成，不会结束在中间某个环节；
  - 事务在执行过程中发生错误，会被回滚(`Rollback`)到事务开始前的状态；
- `C`:(`Consistency`)一致性
  - 在事务开始之前和事务结束以后，数据库的完整性没有被破坏；
- `I`:(`Isolation`)隔离性
  - 数据库允许多个并发事务同时对其数据进行读写和修改，隔离性可以防止多个事务并发执行时由于交叉执行而导致数据不一致。
  - 事务的隔离分为不同级别：读未提交、读提交、可重复读、串行化
- `D`:(`Durability`)耐用性
  - 事务处理结束后，对数据的修改是永久的，即使系统故障也不会丢失。

------

### 数据库的事务隔离级别有哪些？各有哪些优缺点？

##### 数据库的事务隔离级别有哪些？

事务隔离级别主要有四种：

- 读未提交(`READ UNCOMMITED`)
  - 定义：一个事务可以读取另一个事务已修改但未提交的数据
  - 存在问题：脏读
- 读已提交(`READ COMMITED`)
  - 定义：一个事务只能读取另一个事务已经提交的数据
  - 存在问题：不可重复读
- 可重复读(`REPEATABLE READ`)，`MySQL` 默认隔离级别
  - 定义：在一个事务中多次读取同一条记录，结果一致，无论其他事务是否对这条记录做了修改
  - 存在问题：幻读
- 串行化(`SERIALIZABLE`)
  - 定义：所有事务顺序执行，并发性差
  - 不存在脏读、不可重复读、幻读等问题

##### 各有哪些优缺点？

- 隔离级别从上到下，并发性能越来越差，但对于数据的隔离性一致性保证程度越好

##### 脏读定义

- 一个事务读到另一个事务已修改未提交的数据，如果前一个事务回滚或修改之前的值，读到的数据就是错误的

##### 不可重复读定义

- 事务`A`修改某条数据，事务`B` 在事务`A`提交前读取到的数据和提交后读取到的数据不一致

##### 幻读定义

- 在可重复度隔离级别下，普通的查询时快照度，是不会看到别的事务插入的数据的，幻读在当前读下才会出现，幻读仅专指新插入的行

------

### 简述脏读和幻读的发生场景，`InnoDB` 是如何解决幻读的？

##### 脏读发生的场景

- `MySQL`在读未提交的事务隔离级别下，可能会发生脏读
- 当把事务的隔离级别提升到读提交，则可以避免脏读

##### 幻读发生的场景

- 在可重复度隔离级别下，普通的查询时快照度，是不会看到别的事务插入的数据的，因此幻读在当前读下才会出现
- 幻读仅专指新插入的行，更新或删除的行在当前读下出现，不算幻读
- 当前读出现的方式
  - `update` 和 `delete` 语句
  - `select` 加读锁 `lock in share mode`
  - `select` 加写锁 `for update`

##### `InnoDB` 是如何解决幻读的？

- 产生幻读的原因是，行锁只能锁住行，但是新插入记录这个动作，要更新的是记录之间的间隙
- 为了解决幻读的问题，`InnoDB` 在引入间隙锁，间隙锁锁住了记录之间的间隙，在事务的可重复度隔离级别下，间隙锁才会生效
- 间隙锁和行锁合称临键锁( `Netx-Key Locks`)，是`InnoDB` 解决幻读的手段

------

### 并发事务会引发哪些问题？如何解决？

##### 引发的问题

- 脏读
- 丢失修改：撤销一个事务时，把其他事务已提交的更新数据覆盖
- 不可重复读
- 幻读

##### 解决方案

- 解决脏读：将事务隔离提高到读提交级别
- 解决数据丢失：悲观锁
- 解决不可重复读：将事务隔离提高到可重复读级别
- 解决幻读：可重复读事务隔离级别下的临键锁

------

# 锁

### 简述乐观锁以及悲观锁的区别以及使用场景

#### 区别

##### 乐观锁

- 乐观锁在数进行提交更新的时候，才会正式对数据的冲突与否进行检测；
- 在对数据库进行处理的时候，乐观锁并不会使用数据库的锁机制；
- 一般实现乐观锁的方式就是记录版本号，更新数据同步更新版本号，如果发现版本号不一致，则说明数据被其他线程修改了，本次更新不成功
- 乐观锁不会产生死锁，拥有更好的性能

```sql
UPDATE t SET k = 1, version = version + 1 WHERE id = ? AND version = ?
```

##### 悲观锁

- 当我们要对一个数据库中的一条数据进行修改时，为了避免同时被其他人修改，最好的办法就是直接对该数据进行加锁以防止并发
- 悲观锁并发控制实际上是，先取锁再访问的保守策略，为数据处理的安全提供了保证；
- `InnoDB`的行锁分为共享锁和排他锁，两种锁都是悲观锁；
- 共享锁：又称读锁。当一个事务对一条数据加了读锁后，其他事务也能来读数据，可以共享一把锁 
```sql
SELECT * FROM t WHERE id = 1 lock in share mode;
```

- 排他锁：又称写锁。当一个事务对一条数据加了写锁后，其他事务想来访问这条数据只能阻塞等待锁的释放，具有排他性
  - `MySQL`的`SELECT FOR UPDATE`、`INSERT`、`UPDATE`、`DELETE` 都会加锁写锁
```sql
SELECT * FROM t WHERE id = 1 for update;
```

#### 使用场景

##### 乐观锁的使用场景

- 比较适合读取操作比较频繁的场景；
- 如果出现大量的写入操作，数据发生冲突的可能性就会增大：
  - 为了保证数据的一致性，应用层需要不断的重新获取数据
  - 这样会增加大量的查询操作，降低了系统的吞吐量
- `MVCC`通过保存数据的历史版本，根据比较版本号来处理数据是否显示
  - 从而达到读取数据时不需要加锁，就可以保证事务隔离性的效果
  - 体现了乐观锁的思想
  
##### 悲观锁的使用场景

- 比较适合写入操作比较频繁的场景；
- 如果出现大量的读取操作，每次读取的时候都会进行加锁：
  - 这样会增加大量的锁的开销，降低了系统的吞吐量
------

### 什么情况下会发生死锁，如何解决死锁？

#### 发生死锁的情况

- 当并发系统中不同线程出现循环资源依赖，且都在等待别的线程释放资源，从而导致这几个线程都进入无限等待的状态，即出现死锁。

#### 解决死锁

##### 方式一：超时机制

- 两个事务互相等待时，当一个等待时间超过设置的某一个阈值，其中一个事务进行回滚，另一个等待的事务就能继续进行。
- `InnoDB`引擎中，参数`innodb_lock_wait_timeout`用来设置事务的超时时间。
- 缺点：根据`FIFO`的顺序选择回滚对象，无法自主选择
  - 如果超时的事务所占权重比较大，如果事务操作更新了很多行，占用了较多的`undo log`；
  - 这时候回滚这个事务的时间相对另一个事务所占用的时间可能会更多

##### 方式二：`wait-for graph`(等待图)

- `InnoDB`采用的死锁检测方式，是一种较为主动的死锁检测方式；
- 在每个事务请求锁并发生等待时，都会判断是否存在回路，如果存在则会有死锁；
- `InnoDB`引擎会选择`undo`量最小的事务进行回滚。

------

### 简述 `MySQL` 的间隙锁

- 间隙锁，顾名思义，锁的就是两个值之间的空隙，`InnoDB` 是为了解决幻读而引入的新锁
- 间隙锁是在可重复隔离级别下才会生效，如果把事务的隔离级别设置未读提交，就没有间隙锁
- 间隙锁之间不互锁，因为它们都是保护间隙，不允许锁住的间隙里插入值
- 间隙锁和行锁合称临键锁( `Netx-Key Locks`)，每个间隙锁是开区间的，`Netx-Key Locks` 是前开后闭区间

------

# 应用

### 什么是 `SQL` 注入攻击？如何防止这类攻击？

##### 概念

- 通过执行恶意 `SQL`，进而将任意`SQL`代码插入数据库查询，从而使攻击者完全控制 `Web` 应用程序后台的数据库服务器；
- 攻击者可以使用 `SQL` 注入漏洞绕过应用程序验证，比如绕过登录验证登录`Web`身份验证和授权页面；也可以绕过网页，直接检索数据库的所有内容；
- 还可以恶意修改、删除和增加数据库内容

##### 防止

- 不要信任用户的输入，对用户的输入进行校验，可以通过正则表达式，或限制长度；对单引号和双"-"进行转换等。
  - `SQL`对前`''`作为值，后`''`可以插入另外一个`SQL`
  - 双`-`则会注释后面的查询条件
- 不要使用动态拼装`SQL`，可以使用参数化的`SQL`或者直接使用存储过程进行数据查询存取
- 不要使用管理员权限的数据库连接，为每个应用使用单独的权限、有限的数据库连接
- 不要把机密信息直接存放，加密或`HASH`掉密码和敏感的信息；
- 应用的异常信息应该给出尽可能少的提示，最好使用自定义的错误信息对原始错误信息进行包装。

------

### `MySQL` 中 `join`、`left join`、`right join` 的区别是什么？

##### `join/inner join`：内连接

```sql
SELECT * FROM a JOIN a ON a.id = b.id
```

- `join`是`inner join`的简写
- 表示以两个表的交集为主，查出来的是两个表有交集的部分

###### `left join`：左连接

```sql
SELECT * FROM a LEFT JOIN a ON a.id = b.id
```

- 表`a`左连接表`b`，以左表`a`为主，关联上右表`b`，查出来的结果显示左边的所有数据，然后右边显示的是和左边有交集部分的数据

###### `right join`：右连接

```sql
SELECT * FROM a RIGHT JOIN a ON a.id = b.id
```

- 表`a`右连接表`b`，以右表`b`为主，关联上左表`a`，查出来的结果显示右边的所有数据，然后左边显示的是和右边有交集部分的数据

------

### 模糊查询是如何实现的？

##### 方式一：`%`

- `%`表示任意0个或多个字符。可匹配任意长度的字符串

```sql
/* 既有a又有b */
SELECT * FROM t WHERE k LIKE '%a%' AND k LIKE '%b%';
/* xaxbx, xbxax */
```

##### 方式二：`-`

- 表示任意单个字符。匹配单个任意字符，它常用来限制表达式字符长度语句

```sql
SELECT * FROM t WHERE k LIKE '_a_';
/* xax */
```

##### 方式三：`[]`

- 表示括号内所列字符中的一个(类似正则表达式)。
- 指定一个字符、字符串或范围，要求所匹配对象为它们中的任一个

```sql
SELECT * FROM t WHERE k LIKE '[abc]d'; /*因为是连续的，可以写成 [a-c] */
/* ad, bd, cd */
```

##### 方式四：`[^]`

- 表示不在括号所列的单个字符

```sql
SELECT * FROM t WHERE k LIKE '[^abc]d'; /*因为是连续的，可以写成 [a-c] */
/* 非ad, 非bd, 非cd */
```

------

### 索引设计有什么原则？(一般索引需要怎么设计)

#### 建立索引的场景

##### 为常做查询条件的字段建立索引

- 如果某个字段经常用来做查询条件，那么该字段的查询会影响整个表的查询速度
- 为这样的字段建立索引，可以提高整个表的查询速度

##### 为经常需要排序、分组和联合操作的字段建立索引

- 排序操作会浪费很多时间，如果为其建立索引，有效地避免排序操作
- `order by、group by、distinct、union`

##### 尽量使用数据小的列作为索引，长字段使用前缀索引

- 因为`InnoDB`是通过数据页的方式存储数据的，如果某个列数据较小，一个页存储的键值会更多，查询可以减少随机`IO`的读取
- 前缀索引使截取字符串的前一部分建立索引，例如：邮件字段可以先择前7位作为前缀索引(后面的`@xxx`都是相同的)

##### 经常使用多个条件查询时，使用联合索引而不是多个单列索引，同时将离散度高的列放在前面

- 在多个索引时，无论是使用联合索引还是使用最左部分索引，都能使用联合索引；
- 同时最左前缀、索引下推都可以加速查询
- 离散度：`count(distinct(col) : count(*)` 列的全部不同值个数：全部数据行行数

#### 不建立索引的场景

##### 限制索引的数目(索引不是越多越好)

- 浪费磁盘空间
- 更新缓慢：修改表内容，可以会设计到索引的变更，导致页的分裂和合并，因此，索引越多，更新表的时间就越长
- 查询缓慢：创建多余的索引导致优化器花费更多时间选择索引，也可能选择不到所要使用的最佳索引

##### 数据量小的表最好不要使用索引

- 由于数据较小，查询花费的时间可能比遍历索引的时间还要短，索引可能不会产生优化效果

##### 离散度低的列不建议建立索引

- 列的重复值越多，分散度越低，`MySQL`优化器发现索引遍历与使用全表扫描没有太大区别，即使建立了索引，也不一定会遍历

##### 不建议使用无序值(例如：`UUID`)作为索引

- 无序的值作为索引，导致`InnoDB`的缓冲区频繁的加载新的页和刷脏页，导致了大量的随机`IO`
- 无序的值作为索引，写入的值是随机的，可能出现大量页的分裂，分裂会造成索引碎片，得到不够紧凑的索引结构。

##### 尽量少使用经常更新的值用作主键或索引 (主键不推荐有业务含义)

- 和上面一样的问题

##### 删除不再使用或者很少使用的索引

- 表中的数据被大量更新，或者数据的使用方法被改变后，原有的一些索引可能不再需要。
- 应该定期找出这些索引，将它们删除，从而减少索引对更新操作的影响

------

### `MySQL` 有什么调优的方式？

- `MySQL`常见的优化手段分为三个层面：`SQL`和索引优化、数据库结构优化、系统硬件优化等；
- 然而每个大的方向又包含多个小的优化，具体分为：

##### `SQL`和索引优化

此优化方案指的是通过优化`SQL`语句以及索引提高`MySQL`数据库的运行效率，如：

- 使用正确的索引 
  - 避免在`WHERE`查询使用`!=`或`<>`操作符，因为这些操作符会导致查询引擎放弃索引而进行全表扫描
  - 适当使用前缀索引，即定义字符串的一部分来作为索引。因为索引占用磁盘空间越小，相同数据页中能放下的索引值越多，查询效率越快；
- 查询具体的字段，而非全部字符
  - 避免使用`SELECT *`，而是查询需要的字段，这样可以提升速度，以及减少网络传输的带宽压力
- 优化子查询
  - 尽量使用`Join`语句来替代子查询
  - 因为子查询是嵌套查询，而嵌套查询会创建一张临时表，而临时表的创建与销毁会占用一定的系统资源以及花费一定的时间，`Join`语句不会创建临时表，因此性能会更高
- 注意查询结果集
  - 要尽量使用小表驱动大表的方式进行查询
  - `SELECT * FRIM t1 LEFT JOIN t2 ON t1.k = t.k` ；`t1`有`N`条记录，`t2`有`M`条记录；
  - 假设`k`在两个表都要索引，对于`t1`来说是全表扫描，对于`t2`来说是走树搜索，总扫描行数是：`N + N*2*log2M`
  - 驱动表对扫描行数影响更大，因此应该让小表来做驱动表
- 不要在列上进行运算操作
  - 对列字段进行运算操作，会进行隐式转换，优化器放弃使用索引
- 适当增加冗余字段
  - 增加冗余字段可以减少大量的连表查询
  - 因为多张表的连表查询性能很低，适当增加冗余字段，可以减少多张表的关联查询，这是以空间换时间的优化策略
  
##### 数据库结构优化

- 最小数据长度
  - 一遍来说数据库的表越小，它的查询速度就越快；
  - 因此为了提高表的效率，应该将表的字段设置的尽可能小；
- 使用最简单数据类型
  - 能用`int`就不要用`varchar`，因为`int`查询效率更高
- 尽量少定义`text`类型
  - `text`类型的查询效率很低；
  - 如果必须要使用`text`定义字段，可以把字段分离成子表，需要查询此字段时使用联合查询，这样可以提供主表的查询效率
- 适当分表、分库策略
  - 将比较高频的主信息放入主表，其他的放入子表，从而有效提高了查询的效率
  - 把一个库的读和写的压力，分摊给多个库，从而提高了数据库整体的运行效率
  
##### 硬件优化

`MySQL`对硬件的要求主要体现在三个方面：磁盘、网络、内存

- 磁盘
  - 磁盘应该尽量使用有高性能读写能力的磁盘；
  - 比如固态硬盘，这样就可以减少`I/O`运行时间，从而提高了`MySQL`整体的运行效率
- 网络
  - 保证网络带宽的通畅(低延迟)以及够大的网络带宽是`MySQL`正常运行的基本条件；
  - 如果条件允许的化，可以设置多个网卡，以提高网络高峰期`MySQL`服务器的运行效率。
- 内存
  - `MySQL`服务器的内存越大，那么存储和缓存的信息就越多，而内存的性能是非常高的，从而提高了整个`MySQL`的运行效率

------

### `SQL`优化的方案有哪些，如何定位问题并解决问题？

#### `SQL`的优化方案

- 上面的`SQL`和索引优化

#### 定位问题并解决问题

- 慢查询通常的排查手段是先使用慢查询日志功能，查询出比较慢的 `SQL`语句；
- 然后再通过`explain`来查询`SQL`语句执行计划，最后分析并定位出问题的根源，再进行处理

##### 开启慢查询功能

- `set global show_query_log=ON` 开启慢查询，`set global long_query_time=1` 设置慢查询数据，单位秒

```shell
# 或修改 my.cnf
slow_query_log = ON            
slow_query_log_file = /usr/local/mysql/data/slow.log # 日志路径
long_query_time = 1
```

- 开启慢查询日志会对`MySQL`的性能造成一定影响，因此在生产环境慎用此功能；

##### `EXPLAIN`重要字段

- `type`列：表的连接类型，一个好的`SQL`语句至少要达到`range`(索引范围查找)，杜绝出现`all`(全表扫描)
- `key`列：使用到的索引名，如果没有选择索引，值是`NULL`，可以采取强制索引方式
- `rows`列：扫描行数，该值是个预估值
- `extra`列：详细说明，常见的不友好值由：`Using filesort`(需要排序)、`Using temporary`(使用临时表)

##### 解决问题的方案

- 在 `MySQL`中，会引发性能问题的查询，大体分为三种：索引没设计好、`SQL`语句没写好、`MySQL`先错索引

- 索引没有设计好
  - 应该在开发阶段就定位出来，并且修改，可以避免
  - 万一到了线上，可以使用紧急处理，面对高峰期数据库以及被这个语句打挂的情况：
  - 如果是`MySQL 5.6`版本，创建索引都支持`Online DDL`，最高效的做法是直接执行 `alter table` 语句；

- 语句没写好，导致慢查询
  - 如果是平时做修改的话，可以将`SQL`优化，可以避免
  - 如果需要紧急处理，可使用`MySQL 5.7` 提供 `query_rewrite`功能，可以把输入的一种语句改写成另一种模式
```sql
INSERT INTO query_rewrite.rewrite_rules(pattern, replacemenet, pattern_database) VALUES ("SELECT * FROM t WHERE id + 1 = 10000", "SELECT * FROM t WHERE id = 10000 - 1", "db1");

call
query_rewrite.flush_rewrite_rules() /* 使插入的新规则生效，查询重写 */
```

- `MySQL` 选错了索引
  - 应急方案就是给这个语句加上 `force index()`，然后使用查询重新功能，给原来语句加上该操作
  - 后期可以考虑修改语句或新键一个更合适的索引来提供给优化器做选择

------

# 特性

### MySQL 有哪些常见的存储引擎？它们的区别是什么？

- `InnoDB`存储引擎，支持事务，行锁设计、支持外键、多版本并发控制
- `MyISAM`存储引擎，不支持事务，表锁设计，支持全文索引
- `Memory`存储引擎，不支持持久化，适合存储临时表
- `Archive`存储引擎，非常适合存储归档数据，仅支持`SELECT` 和 `INSERT`

------

### `MySQL` 为什么会使用 `InnoDB` 作为默认选项

- 选择的原因是因为`InnoDB`功能方面有较多的优势，包括：支持事务、灾难恢复性好、使用行级锁、实现了缓冲处理、支持外键
- 支持事务
  - `InnoDB`最重要的一点就是支持事务，可以说这是`InnoDB`成为`MySQL`默认存储引擎的一个非常重要原因
  - 此外`InnoDB`还实现4种事务隔离级别，使得对事务的支持更加灵活；
- 灾难恢复性好
  - `InnoDB`通过 `commit、rollback、crash-safe`来保障数据的安全
  - `crash-safe`就是指如果服务器因为硬件或软件的问题而崩溃，不管当时数据是怎样的状态，在重启`MySQL`后，`InnoDB` 都会自动恢复到发生崩溃之前的状态
- 使用行锁
  - `InnoDB`改变了`MyISAM`，实现了行锁。
  - 虽然`InnoDB`的行锁机制是通过索引来完成的，但毕竟在数据库中大部分的`SQL`语句都要使用索引来检索数据。
  - 行锁机制也为`InnoDB`在承受高并发的环境下增强了不小的竞争力
- 实现了缓冲处理
  - `InnoDB`提供了专门的缓冲池，实现了缓冲管理，不仅能缓冲索引也能缓冲数据，常用的数据可以直接从内存中处理，比从磁盘获取数据处理速度更快
- 支持外键
  - `InnoDB`支持外键约束，检查外键、插入、更新和删除，以确保数据的完整性。存储表中的数据时，每张表的存储都按主键顺序存放。

###### 拓展：插入缓冲(`Insert Buffer`)

- 解决的问题：对于非聚簇的辅助索引叶子节点的插入是随机`IO`，导致插入性能下降
- 解决的方案：
  - 对于非聚簇索引的插入或更新操作，不是每一次直接插入到索引页中，而是先判断插入的非聚簇索引是否在缓冲池中，若在，则直接插入；
  - 若不在则先放入到一个`Insert Buffer`对象中。
  -再以一定的频率，将多个插入和辅助索引页子节点的`merge`(合并)操作，提供非聚簇索引插入的性能
- `Insert Buffer`的使用条件：
  - 索引是辅助索引
  - 索引不是唯一的
- `Change Buffer`:`Insert Buffer`的升级版
  - 对`DML`(`INSERT、DELETE、UPDATE`)都进行缓冲，分别为`Insert Buffer`、`Delete Buffer`、`Purge Buffer`
- `Insert/Change Buffer`是一棵`B+`树
  
##### 拓展：两次写(`doublewrite`)

- 解决问题：部分写失效，即`InnoDB`存储正在写某个页到表的过程中发生宕机，导致该页数据部分丢失
-`doublewrite`的组成：内存中的`doublewrite buffer`大小`2MB` + 磁盘共享表空间中连续`128`个页大小`2MB`
- 执行过程：
  - 在对缓冲池的脏页进行刷新时，并不直接写磁盘，而是会通过`memcpy`函数将脏页复制到内存中的`doublewrite buffer`；
  - 在分两次将`doublewrite buffer`写入共享表空间中(共享表空间的页是连续的，开销不大)，然后调用`fsync`函数，同步磁盘；
  - `doublewrite`页写入后，在将其页写入各个表空间文件中，此时写入时离散的；
  - 如果在写入磁盘的过程中发生了崩溃，在恢复过程中，`InnoDB`存储引擎可以从共享表空间中的`doublewrite`中找到该页的一个副本，将其复制到表空间文件，再应用重做日志。

##### 拓展：自适应哈希索引

- `InnoDB`存储引擎会监控对表上各索引页的查询。如果观察到建立哈希索引可以带来速度提交，则建立哈希索引，称自适应哈希(`AHI`)
- 自适应哈希是通过缓冲池的`B+`数页构造而来的，不需要对整张表构建哈希索引
- 自适应哈希要求：
  - 对这个页的连续访问模式(如：`k=xx`)必须是一样的
  - 以该模式访问了100次
  - 页通过该模式访问了`N`(=页记录/16)次
  
##### 拓展：`InnoDB`和`MyISAM`的区别

- `InnoDB`支持事务，`MyISAM` 不支持事务，这是`MySQL`将默认存储引擎从`MyISAM`变成`InnoDB`的重要原因之一；
- `InnoDB`支持外键，`MyISAM`不支持。对一个包含外键的`InnoDB`表转为`MyISAM`会失败；
- `InnoDB`是聚簇索引，`MyISAM`是非聚簇索引。
  - 聚簇索引的文件存放在主键索引的叶子节点上，因此`InnoDB`必须要有主键，通过主键索引效率很高；
  - 但是`InnoDB`的辅助索引需要两次查询，先查询主键，然后再通过主键查询到数据，因此主键不应该过大；
  - 非聚簇索引数据文件是分离的，索引保存的是数据文件的指针。
- `InnoDB`不保存具体行数，`COUNT(*)`需要全表扫描，`MyISAM`虽然记录了行数，但是存在`WHERE`条件时，和`InnoDB`的查询一样
  - 优化器对`InnoDB`的`COUNT(*)`做了优化，会选择扫描成本最小的索引来查询(索引小，扫描的数据页少)，所以效率`COUNT(*) > COUNT(id)`
  - `InnoDB COUNT(*)`不保存行数的原因：多版本并发控制(`MVCC`)
- `InnoDB`最小的锁粒度是行锁，`MyISAM`最小的锁粒度是表锁，`InnoDB`的并发访问性能更好。
  - `InnoDB`的行锁是现实在索引上的，而不是锁在物理行记录上，所以如果访问没有命中索引，也无法使用行锁，将要退化为表锁。 

------

### 数据库设计的范式是什么？

- 第一范式：每个列都不可以再拆分
- 第二范式：在第一范式基础上，非主键列完全依赖于主键，而不能依赖主键的一部分
- 第三范式：在第二范式基础上，非主键只依赖主键，不依赖于其他非主键，即没有传递关系

------

### 数据库反范式设计会出现什么问题？

##### 反范式设计会出现的问题

- 存在大量冗余数据，数据维护成本更高，可能伴随者删除异常、插入异常、更新异常

##### 不符合第二范式的例子

- 例如：在选课关系表(学号，课程号，成绩，学分)，主键为(学号，课程号)
- 但是非主属性学分仅仅依赖于课程，对主键仅依赖一部分，不是完全依赖，不符合第二范式(`2NF`)
- 存在问题：
  - 数据冗余：每条记录都含有相同信息(学分)
  - 删除异常：删除所有的学生成绩，就把课程信息全删除了
  - 插入异常：学生未选课，则没有课程信息
  - 更新异常：调整课程学分，所有行都调整
- 调整方法
  - 选课关系表(学号，课程号，分数)
  - 课程表(课程号，学分)

##### 不符合第三范式的例子

- 例如：学生表(学号、姓名、年龄、学院名称、学院电话)
- 该表存在学院电话不依赖学号，依赖学院名称，存在传递关系，所以不符合第三范式
- 可能存在的问题
  - 数据冗余：每条记录都含有相同信息
  - 更新异常：修改学院电话时，要更新所有包含该学院的记录
- 调整后的表
  - 学生表：(学号，姓名，年龄，所在学院)
  - 学院表：(学院，电话)
  
##### 反范式化的优点

- 数据冗余将带来很好的读取性能(以空间换时间，因为不需要`join`很多表)
- 例如：学生表与选课表，假定选课表要经常被查询，而且在查询中要显示学生姓名，如果这个需要被大范围、高频率的执行，可能会因为表关联造成一定程度的影响
- 现在评估学生改名的需求很少，那么可以把学生姓名冗余到选课表中。

------
### `MySQL` 中 `varchar` 和 `char` 的区别是什么？

##### 区别一：定长和变长

- `char` 定义的列为长度固定的字符串
  - 当所插入的字符串长度超出定义的长度时，如果时严格模式，则会拒绝插入并提示错误信息，如果是非严格模式，则会截断然后插入；
  - 当所插入的字符串长度小于定义的长度是，则将在它们的右边填充空格以达到指定的长度，在检索到的值，拖尾的空格被删除；
- `varchar` 定义的列为长度可变字符串
  - 占用的字节空间分为长度前缀+字符实际大小，长度前缀为当字符串实际的字节数，当字节数大于255时，则用2个字节去存储其长度，否则使用1个字节去存储；（`M+1/2`)

##### 区别二：存储的容量不同

- `char`存储的字符个数在 `0～255`之间，和编码无关
- `varchar` 存储的字符个数在 `0～65532`，`varchar`的最大有效长度由最大行大小和使用字符集确定
  - 行长度限制，`MySQL`要求一个行的定义长度不能超过65535字节
  - `varchar`字符类型使用1个字节来保存控制信息(65532=65535-1-2[长度信息])
  - 编码长度限制
    - 字符类型是`gbk`，每个字符最多占2个字节，最大长度不能能超过`32766`；
    - 字符类型是`utf8`,每个字符最多占3个字节，最大长度不能超过`21845`；
    - 定义时超过限制，则`varchar`字段会被强行转为`text`类型

##### 拓展：使用场景

- `char`
  - 存储的数据长度基本一致，不需要空格
  - 例如：手机号、`UUID`、密码加密后的密文
- `varchar`
  - 数据长度不一定，长度范围变化较大的场景

------

# 持久化

### 简述 MySQL 三种日志的使用场景(作用)

- 我们经常讲的 `MySQL`三种日志是：`binlog`、`redo log`、`undo log`，其中`redo log`和`undo log` 是`InnDB`的实现关键

##### `binlog`的使用场景

- `binlog` 是 `MySQL Server` 层维护的一种二进制日志，其主要是用来记录`MySQL`数据更新或潜在发生更新的`SQL`语句，并以"事务"的形式保存在磁盘中；
- `binlog`的使用场景是：
  - 主从复制：`MySQL` 复制是 `Master`端开启`binlog`，`Master`把它的二进制日志传递给`Slaves`，并在`Slaves`端回放来达到`Master-Slave`数据一致的目的
  - 数据恢复：通过`mysqlbinlog` 工具恢复数据
  - 增量备份：`binlog` 是通过追加的方式进行写入的，可以通过 `max_binlog_size`参数设置每个`binlog`文件的大小，当文件大小到达给定值后，会生成新的文件来保存日志；

#### `redo log` 的使用场景

- 在 `InnoDB`中，数据一致性由`redo log`来保证，使用的是`WAL(Write-Ahead Logging)`机制，即先写日志再写数据；
- `InnoDB` 使用这种方式在进行故障恢复时，会将 `redo log`中的日志重做一遍，也就是将系统中未提交的事务重新执行；
- 默认情况下，`redo log`记录在`ib_logfile0` 和 `ib_logfile1`两个文件，分别用`write pos` 记录写入位置；
- 用`checkpoint`记录整个系统当前日志已经同步的位置，`checkpoint`保证了未提交的事务重新执行。

##### `undo log` 的使用场景

- `undo log`是事务原子性的保证，主要作用是回滚和多版本控制(`MVCC`)
- `undo log`主要记录了数据的逻辑变化，比如：一条`INSERT`语句，对应一条`DELETE`的`undo log`；
- 如果用户执行的事务或语句由于某种原因失败了，可以利用这些`undo`信息将数据回滚到修改前的数据状态；
- `MVCC`也是通过`undo log`来保证快照读的逻辑。

##### 拓展：`binlog`相关

- 查询 `binlog`日志的两种方式
  - 原因：`binlog`日志是二进制格式，无法直接进行查看
  - `mysqlbinlog`:`/usr/bin/mysqlbinlog mysql-bin.000009`
  - 命令解析 `SHOW BINLOG EVENTS [IN 'binlog_name']`
- `MySQL` 通过`sync_binlog`控制刷盘时机
  - 0(系统自行判断)，1(每次`commit`都会写入)，`N`(每`N`次事务才写入)
  - `MySQL 5.7.7`后默认设置为`1`

##### 拓展：`redo log` 和 `binlog` 的区别

- `redo log` 是 `InnoDB` 引擎特有；`binlog`是`MySQL`的`Server`层实现的，所有引擎都可以用
- `redo log` 是物理日志，记录的是"在某个数据页上做了什么修改"; `binlog` 是逻辑日志，记录的是这个语句的原始逻辑，比如："给`ID=2`的行`c`字段加1"
- `redo log` 是循环写的，空间固定会用完；`binlog`是可以追加写的，当文件大小到达`max_binlog_size`值后，会生成新的文件来保存日志
- `redo log`适用于崩溃恢复(`crash-safe`), `binlog`日志适用于主从复制和数据恢复
  - `binlog`没有崩溃恢复的能力，没有能力恢复数据页;
  - 假设现在数据更新后，脏页还没刷盘，出现宕机，在`binlog`中该事务已经提交了，无法恢复数据页


##### 拓展：`redo log` 和`undo log`的区别

- `redo log` 和 `undo log` 都是用来恢复日志，但不是逆向过程
- `redo log`是物理日志，记录的是数据页的物理修改；`undo log`是逻辑日志，用来回滚行记录到某个版本，根据每行记录进行记录
- `redo log` 是读顺序写的， `undo log`是随机读写

##### 拓展：数据恢复

##### 拓展：正常运行中的实例，数据写入的最终落盘，是从`redo log`更新过来的吗？

- `redo log`并没有记录数据页的完整数据，所以它并没有能力自己去更新磁盘数据页，也就是不存在：数据最终落盘，是由`redo log`更新过去的情况；
- 如果是正常运行的实例
  - 数据页被修改后，跟磁盘的数据页不一致，称为脏页；
  - 最终数据落盘，就是把内存中的数据页写盘；
- 在崩溃恢复场景中
  - `InnoDB` 如果判断一个数据页可能在崩溃的时候丢失了更新，就会将它读到内存，然后让`redo log`更新内存内容
  - 更新完成后，内存页变成脏页，脏页旧回到上面的情况
  
##### 拓展：什么时候写`redo log buffer`，什么时候写`redo log file`

- `redo log buffer` 就是一块内存，用来先存`redo`日志的，也就是说，在数据的内存被修改了(执行写入操作)，`redo log buffer` 也写入日志；
- `redo log file` 是在执行 `commit`语句是做的

##### 拓展：怎样刷脏页

- 写数据
  - `buffer pool`维护着一个脏页列表，假设现在 `redo log`的 `checkpoint`记录的`LSN`为10
  - `write pos`为 11，修改后该页的`LSN`为12；
  - 则写`redo log`同时该页也会被标记为脏页，记录到脏页列表中
- 刷脏页
  - 假设内存不足，需要淘汰该页，进行`flush`操作，将该页刷新到磁盘中并从脏页列表中移除；
  - `redo log` 没有改变
  - 当`redo log`执行`checkpoint`时，如果查到这个刚才淘汰的`LSN`没有在脏页列表中，就跳过不做该数据页的`flush`操作
- `LSN`：日志序列号
  - 每个数据页头部有`LSN`，8字节，每次修改都会变大
  - 对比这个`LSN`跟`checkpoint`的`LSN`,比`checkpoint`小的一定是干净页 

------

### 简述 `MySQL MVCC` 的实现原理

##### 概念

- `MVCC`，多版本并发控制。在`MySQL InnoDB` 中主要是为了提高数据库并发性能，用更好的方式去处理读-写冲突，做到即使有读写冲突时，也能做到不加锁，非阻塞并发读

##### 实现原理

- 在可重复读隔离级别下，事务启动相当于拍了个快照；
- 为了实现快照功能，`InnoDB`内部为每一行添加了两个隐藏字段：该行最后修改的事务`ID`(`DB_TRX_ID`) 和 回滚指针(`DB_ROLL_PTR`);
- `InnDB`里每个事务都要一个严格递增的事务`ID`，每次事务更新时，都会生成一个新的数据版本，把事务`ID`保存在 `DB_TRX_ID`中；
- 旧的数据保存`undo log`中，并把该`undo log`的指针保存在`DB_ROLL_PTR`；

##### 对于 `SELECT` 的逻辑

- 开启事务`A`事务`ID`为 100，保存未提交事务ID数组，在执行`SELECT k FROM T WHERE id = 1` 查到 `k = 10`；
- 这时候另一个客户端开始事务`B` 事务`ID`为101，执行更新 `UPDATE T SET k = k+1 WHERE id = 1`，然后提交事务：
  - 同时会把`DB_TRX_ID`更新为101，同时`undo log`保存`id=1,k=10`的回滚段，并把其指针保存在行记录中；
- 这时候事务`A`再执行`SELECT k`：
  - 先会找到`id=1`的记录，然后判断`DB_TRX_ID`比自己保存的最后修改事务ID要大或保存的ID在未提交数组中，则通过`DB_ROLL_PTR`找到上一个版本；
  - 在上一个版本中发现`DB_TRX_ID` 小于或等于自己的事务`ID`，则是自己需要找的值，返回`k=10`
- 在同一个事务中，两次找到的数据是一致的，这就是`MVCC`的快照功能，即快照读。

##### 对于 `UPDATE` 的逻辑

- 需要说明的是，在`MVCC`中，更新数据都是先读后写的，这个读，只能是读取当前的值，即当前读；
- 还是上面的场景，事务`A`在第二次`SELECT k` 后，进行 `UPDATE T SET k = k+1 WHERE id = 1`;
- 在更新时，当前读拿到的数据时 `k=11`，更新后`k=12`，并把当前事务`ID`记为`100`
  - 同时也因为事务`B`已经提交了，不会因为事务`B(TRX_ID=102)`比`事务A(TRX_ID=101)`而造成脏读问题；
- 这时候如果再次执行 `SELECT k` 时，因为最后一次修改的事务`ID`是事务`A`自己的，不会再去找`undo log`，返回`k=12`，逻辑无误。
- 再事务100未提交前，再开启事务103，保存的未提交事务未[100, 103]，在查询`SELECT k`时先找到记录为100的数据，发现它在未提交事务中，再往前找到101记录 

##### 拓展：`InnoDB`向数据库中存储的每一行添加三个字段
- 6字节的`DB_TRX_ID`:标识插入或更新该行的最后一个事务的事务标识(删除被视为更新，行中的特殊位被设置位删除标记)
- 7字节的`DB_ROLL_PTR`:回滚指针，指向写入回滚端的`undo log`记录，回滚段包含重建之前行所需的内容
- 6字节的`DB_ROW_ID`：随着新行插入而单调增加的行`ID`，如果 `InnoDB`自动生成聚簇索引，则该索引包含行`ID`值。否则，该 `DB_ROW_ID`列不会出现在任何索引中

##### 拓展：当前读

- 概念
  - 读取的是记录的最新版本，读取时还要保证其他并发事务不能修改当前记录，会对读取的记录进行加锁
- 当前读使用场景
  - `select lock in share mode`(共享锁)
  - `select for update`(排他锁)
  - `update、insert、delete`

##### 拓展：快照读

- 不加锁的 `select`操作就是快照读，即不加锁的非阻塞读；
- 当事务级别是串行级别下的快照读退化成为当前读；
- 快照读的实现是基于多版本并发控制(`MVCC`)，快照读可能读到的并不是数据的最新版本

------

### 简述什么是两阶段提交？

##### 概念

- 为了保持分布式事务的原子性，事务管理器使用了一个标准的两层恢复机制，称为两阶段提交协议(`2PC`);
- 这种两阶段协议可以确保事务的更新可以提交给所以参与的资源，或确保更新的资源完全回滚到事务开始的状态；
- 通过这种方式，提交协议就可以保证数据的完整性。

##### `MySQL`的两阶段提交

- `MySQL`使用两阶段提交主要解决：`binlog`和`redo log`的数据一致性的问题
- 步骤：
  - 引擎将新的数据更新到内存中，同时将这个更新操作记录到`redo log`里面，此时`redo log` 处于 `prepare`状态。然后告知执行器执行完成了；
  - 执行器生成这个操作的`binlog`，并把`binlog`写入磁盘；
  - 执行器调用引擎的提交事务接口，引擎把刚刚写入的`redo log`改成提交(`commit`)状态，更新完成。

##### 拓展：崩溃恢复时的判断规则

- 如果`redo log`里面的事务是完整的，就是已经有了`commit`标识，则直接提交；
- 如果`redo log`里面的事务只有完整的`prepare`，则判断对应的事务`binlog`是否存在并完整：
  - 如果是，则提交事务；
  - 否则，回滚事务
  - 每个事务`binlog`末尾，会记录一个 `XID event`，标志着事务是否提交成功
  
##### 拓展：`redo log` 是先写 `buffer`还是`file`

- `buffer`在`SQL`的`INSERT`等语句执行过程中写入；
- `file` 是在执行`SQL`的`commit`时候写入的。

------

### 简述 `MySQL` 的主从同步机制，如果同步失败会怎么样？

##### `MySQL`的主从同步机制

- 1-主库`A`接收到客户端的更新请求后，执行内部事务的更新逻辑，同时写`binlog`;
- 2-在备库`B`上通过 `change master`命令，设置主库`A`的`IP`、端口、用户名、密码，以及要从哪个位置开始请求`binlog`，这个位置包含文件名和日志偏移量；
- 3-在备库`B`上执行`start slave`命令，这时候备库会启动两个线程，即`IO`线程和`SQL`线程，其中`IO`线程负责与主库建立连接；
- 4-主库`A`开启`binlog dump`线程，校验完用户名、密码后，开始按照备库`B`传过来的位置，从本地读取`binlog`，发给`B`；
- 5-备库`B`拿到`binlog`后，写到本地文件，这个文件称为中转日志(`relay log`)；
- 6-备库的`SQL`线程读取中转日志，解析出日志里的命令，并执行。

##### 拓展：`MySQL`主从复制模式

##### 拓展：`MySQL`主从复制方式

- 异步模式:默认模式，这种模式下，主节点不会主动 `push bin log`到从节点
  - 当从节点连接主节点时，会主动从主节点出获取最新的`bin log`文件；
  - 这种模式下可能导致`failover`(故障转移)，因为从节点没有即时将最新的`binlog`同步
- 半同步模式
  - 这种模式下主节点只需接收到其中一台从节点的返回信息，就会`commit`；
  - 否则需要等待直到超时时间，然后切换成异步模式再提交；
  - 可以时主从数据库的数据延迟缩小，提高安全性，性能上会有一定的降低，响应时间会变长
- 全同步模式
  - 主节点和从节点全部执行了`commit`并确认会给客户端返回成功
- GTID

##### 拓展：`binlog` 日志格式

- `STATMENT` ：基于 `SQL` 语句的复制，每一条会修改数据的 `SQL` 语句会记录到 `binlog` 中
  - 优点：不需要记录每一行的变化，减少了日志量
  - 缺点： 在某些情况下会导致主从数据不一致
- `ROW`：基于行的复制，仅需要记录哪条数据被修改了
  - 优点： 不会出现某些特定情况下无法被正确复制的问题
  - 缺点： 会产生大量的日志，尤其是 alter table 的时候会让日志暴涨
- `MIXED`：混合模式

------

### 数据库的读写分离的作用是什么？如何实现？

##### 作用

- 分摊主库压力，
- 负载均衡

TODO

------

### 简述数据库中什么情况下进行分库，什么情况下进行分表？

- 分表分库有两种方式：垂直切分和水分切分

##### 垂直分库时机

- 垂直分库是基于业务分类的，每一个独立的服务都有拥有自己的数据库；
- 当业务过于庞大，数据库压力大，且业务不能互相影响时；
- 通过垂直拆分，降低业务的耦合性，同时也能缓解服务器压力；

##### 垂直分表时机

- 某个表字段较多，且每次经常查询的仅是部分数据
- 可以新建一张拓展表，将不经常用或字段长度较大的字段拆分出去
- 大表拆小表，更便于开发与维护，也能避免跨页问题
  - 一条记录占用空间过大会导致跨页，造成额外的性能开销   
  
##### 水平切分

- 水平分表分为：库内分表和分库分表
- 水平分表时机
  - 垂直分表仅能解决多字段问题，但是仍会存在单表数据量过大、单表读写、存储的瓶颈问题
  - 通过库内分表，可以解决单一表数据量过大的问题 
- 分库分表时机
  - 库内分表仍存在竞争同一个物理机的`CPU`、内存、网络`IO`
  - 通过分库分表，可以使切分出的子表，分散到不同的数据库中，从而使得单个表的数据量表小，达到分布式的效果
