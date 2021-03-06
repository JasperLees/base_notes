## 进程与线程

### 进程和线程之间有什么区别？

- 进程是资源分配的最小单位，线程是`CPU`调度的最小单位
- 一个线程只能属于一个进程，而一个进程可以由多个线程
- 进程在执行过程中拥有独立的内存单元，而多个线程共享进程的内存
- 进程上下文切换比进程上下文切换要慢，开销大
- 进程间通信是在通过内核读写数据的，线程间可以直接读写进程数据段来进行通信(需要进程同步和互斥手段辅助，以保证数据的一致性)
- 进程间不会互相影响，一个线程挂掉将导致整个进程挂掉

------

### 进程间有哪些通信方式？

进程间的通信方式有：管道、消息队列、共享内存、信号量、信号、套接字

##### 管道

- 管道一种半双工的通信方式，数据只能单向流动，管道实质上是一个内核缓冲区，且以先进先出的方式存取数据。
- 管道分为命名管道和匿名管道
- 匿名管道它的优点是：简单方便；缺点是：1) 因为管道局限于单向通信且缓冲区有限， 所以它的通信效率低，不适合进程间频繁地交换数据，2) 只能在父子进程间使用，
- 命名管道，和匿名管道相比，它的优点是：可以实现任意关系的进程间的通信；缺点是：长期存在系统中，使用不当容易出错

##### 消息队列

- 消息队列是保存在内核中的消息链表
- 优点：可以实现任意进程间的通信，并且通过系统调用函数来实现消息发送和接收之间的同步，无需考虑同步问题
- 缺点：1) 消息队列不适合比较大数据的传输，因为每个消息体都有最大长度限制，同时全体消息也有总长度上限，2) 通信过程中，存在用户态与内核态之间的数据拷贝开销

##### 共享内存

- 共享内存就是映射一段能被其他进程所访问的内存，这段共享内存由一个进程创建，但多个进程都可以访问
- 优点：进程可以直接读写这块内存而不需要进行数据拷贝，提高效率
- 缺点：1）多个进程同时修改同一个共享内存，会发生冲突 2）共享内存只能在同一计算机系统中共享

##### 信号量

- 信号量是一个整型的计数器，主要用于实现进程间的互斥与同步，而不是用于缓存进程间通信的数据
- 优点：信号量解决了止多进程竞争共享资源，而造成数据的错乱的
- 缺点：信号量有限

##### 信号

- 信号是一种比较复杂的通信方式，用于通知接收进程某个事件已经发生

##### 套接字

- 套接字通信不仅可以跨网络与不同主机的进程通信，还可以在同主机进程通信
- 优点：1) 传输数据为字节级，传输数据可自定义 2) 适合客户端和服务端之间信息实时交互 3) 可以加密，数据安全性强
- 缺点：需对传输的数据进行解析，转化成应用级的数据
------

### 进程通信中的管道实现原理是什么？

------

### Linux 进程调度中有哪些常见算法以及策略？

Linux进程的调度分为实时调度(创建指定算法和优先级)和普通分时调度(创建指定算法和nice值)

#### 实现调度

- 实时调度优先级为 0〜99， 它主要包含两个算法SCHED_FIFO(先到先服务的实时进程)，SCHED_RR(时间片轮转的实时进程)

##### 这两个算法的相同点：

- 高优先级的进程会抢占低优先级的，高优先级的进程运行期间，低优先级进程不能能抢占，只能等到高优先级的进程主动退出

##### 不同点：

- SCHED_FIFO(先到先服务的实时进程)：对于同等优先级的进程，先运行的进程会一直占据CPU，只能等到先运行的进程主动退出，后续进程才能得到时间片

- SCHED_RR(时间片轮转的实时进程)：对于同等优先级的进程，各个进程会轮流运行一定的时间片(大约100ms)

#### 普通分时调度

- 普通分时调度根据nice值取判断优先级，nice值取值范围是(-20, +19)，对应进程的优先级为100〜139，普通分时调度主要基于CFS调度算法

##### CFS

- CFS通过进程物理运行时间和nice的加权值，计算得出虚拟运行时间 vruntime(vruntime=pruntime[物理运行时间]/weight[nice加权值]*1024)

- 物理运行时间越少、nice值越低的进程vruntime越小，调度程序时只需选择最小虚拟运行时间即可
- Linux CFS 通过红黑树维护可运行的任务，树的索引值为 vruntime
- 由于红黑树是平衡的，只需要找到最左侧的节点即为优先级最高的任务(O(logN))

##### 拓展：进程调度

- 就是从进程的就绪队列(阻塞)中按照一定的算法选择一个进程并将`CPU`分配给他运行，以实现进程的并发运行

##### 拓展：进程的三态模型



------

### 简单介绍进程调度的算法

------

### 进程有多少种状态？

##### 有三种状态

- 运行状态(`running`)：进程以获得`CPU`，并在`CPU`上运行
- 就绪状态(`ready`)：进程获得了除了`CPU`之外的所有的必要资源，只要获得CPU就可以立即执行 
- 阻塞态/等待态(`wait`)：进程由于发生某些事件而暂时无法继续执行，放弃处理机而处于暂停状态

##### 状态转换

- 就绪-> 执行：进程调度程序为位置分配了处理机
- 执行-> 就绪：分配的时间片以用完
- 执行-> 阻塞：因等待某件事发生而无法继续执行，例如：等待`IO`完成
- 阻塞-> 就绪：等待的事件以及发生，例如：`IO`已经完成

------

### 线程有多少种状态，状态之间如何转换

##### 状态

- 6种
- 新创建(`New`)、可运行(`Runnable`)、被阻塞(`Blocked`)、等待(`Waiting`)、计时等待(`Timed Waiting`)、被终止(`Terminated`)

------

### 线程间有哪些通信方式？

------

### 两个线程交替打印一个共享变量

------

### 什么情况下，进程会进行切换？

------

### 进程空间从高位到低位都有些什么？

------

## 多路复用(Redis)

### 简述 `socket` 中 `select、poll、epoll` 的使用场景以及区别，各自支持的最大描述符上限以及原因是什么？

##### 区别 和 使用场景

|            | select                                                       | poll     | epoll                                                        |
| ---------- | ------------------------------------------------------------ | -------- | ------------------------------------------------------------ |
| 底层实现   | 数组                                                         | 链表     | 红黑树(未就绪fd) + 双向链表(就绪fd)[epoll_create创建]        |
| 最大连接数 | 1024(32位) 或2048(64位)                                      | 无上限   | 无上限                                                       |
| IO效率     | 每次调用都进行线性遍历(轮询)，时间复杂度为O(n)               | 同select | 基于事件驱动实现的，每当fd就绪，系统注册的回调函数就会被调用，将就绪的fd放到双向链表中，epoll_wait直接从链表中取节点，时间复杂O(1) |
| fd拷贝     | 每次调用select，都需要把fd集合从用户态拷贝到内核态           | 同select | 调用epoll_ctl时拷贝进内核并保存，之后每次epoll_wait不拷贝    |
| 适用场景   | 连接数不多的场景，连接数不多适用epoll并为其建立文件系统、红黑树和链表，效率反而不高 | 同select | 对于连接特别多，活跃的连接较少的场景，例如为一个需要处理上万的连接的服务器。如果过多的活跃连接，因为处理问题不及时，性能有所下降 |


##### 各自支持最大描述符上限以及原因

- `Select`：最大描述符上限和系统内存有关，基于数组，32位机器1024，64位默认2048
- `poll`  ：没有最大连接数限制，因为它基于链表存储
- `epoll` : 没有最大连接的上限(`1G`的内存上能监听约10万个端口)

##### 拓展：概念

- `select`：监控三种文件描述符，读(`write_fds`)、写(`read_fds`)、异常(`except_fds`)
  - 调用后阻塞，直到有描述符就绪，或超时返回
- `poll`:查询每个`fd`对应的设备状态
  - 如果就绪，则加入等待队列，否则挂起等待
- `epoll`:通过`epoll_ctl`注册`fd`：
  - 一旦该`fd`就绪，内核就会采用类似`callback`的机制来激活`fd`, `epoll_wait`便可以收到通知

------

### `epoll` 中水平触发以及边缘触发有什么不同 

- 当`epoll_wait`检测到描述符事件发生此事件，通知应用程序时；
- 水平触发(`LT level trigger`)
  - 应用程序可以不立即处理该事件，下次调`epoll_wait时`，会再次响应应用程序并通知此事件
  - 默认模式
- 边缘触发(`ET edge trigger`)
  - 应用程序必须立即处理该事件，如果不处理，下次调用`epoll_wait`时，不会再次响应应用程序并通知此事件

------

### 简述同步与异步的区别，阻塞与非阻塞的区别

------

## 内存

### 操作系统如何申请以及管理内存的？

------

### 简述 Linux 零拷贝的原理

------

### 如何调试服务器内存占用过高的问题？

------

### 简述操作系统如何进行内存管理

------

### 操作系统如何申请以及管理内存的？

------

### 简述操作系统中 malloc 的实现原理

------

### Linux 中虚拟内存和物理内存有什么区别？有什么优点？

------

### 操作系统中，虚拟地址与物理地址之间如何映射？

------

### 简述操作系统中的缺页中断

------

## 应用

### 简述几个常用的 Linux 命令以及他们的功能。

------

### Linux 下如何排查 CPU 以及 内存占用过多？

------

### Linux 下如何查看 CPU 荷载，正在运行的进程，某个端口对应的进程？

------

### Linux 如何查看实时的滚动日志？

------

### 简述 Linux 系统态与用户态，什么时候会进入系统态？

------

### 简述 LRU 算法及其实现方式

------

### 什么时候会由用户态陷入内核态？

------

### Linux 下如何查看端口被哪个进程占用？

------

### BIO、NIO 有什么区别？怎么判断写文件时 Buffer 已经写满？简述 Linux 的 IO模型

------

### 简述自旋锁与互斥锁的使用场景

------

### 简述 mmap 的使用场景以及原理

------

### 简述 traceroute 命令的原理

------
