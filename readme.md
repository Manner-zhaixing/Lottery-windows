## 建表
```sql
create table if not exists inventory(
    id int auto_increment comment '奖品id，自增',
    name varchar(20) not null comment '奖品名称',
    description varchar(100) not null default '' comment '奖品描述',
    picture varchar(200) not null default '' comment '奖品图片',
    price int not null default 0 comment '价值',
    count int not null default 0 comment '库存量',
    primary key (id)
)default charset=utf8mb4 comment '奖品库存表，所有奖品在一次活动中要全部发出去';
insert into inventory (id,name,picture,price,count) values (1,'谢谢参与','img/face.png',0,0);
insert into inventory (name,picture,price,count) values 
('篮球','img/ball.jpeg',100,1000),
('水杯','img/cup.jpeg',80,1000),
('电脑','img/laptop.jpeg',6000,200),
('平板','img/pad.jpg',4000,300),
('手机','img/phone.jpeg',5000,400),
('锅','img/pot.jpeg',120,1000),
('茶叶','img/tea.jpeg',90,1000),
('无人机','img/uav.jpeg',400,100),
('酒','img/wine.jpeg',160,500);
create table if not exists orders(
    id int auto_increment comment '订单id，自增',
    gift_id int not null comment '商品id',
    user_id int not null comment '用户id',
    count int not null default 1 comment '购买数量',
    create_time datetime default current_timestamp comment '订单创建时间',
    primary key (id),
	key idx_user (user_id)
)default charset=utf8mb4 comment '订单表';
```

- 指定一个特殊的id来标识“谢谢参与”，通过调节“谢谢参与”的count来控制它被抽中的概率。

## 后端接口
|请求路径|请求方式| 请求参数 |说明|
| :--- | :--- |:-----| :--- |
|/|GET|      |返回抽奖转盘页面|
|/gifts|GET|      |返回所有奖品的详细信息，用于往转盘里填充内容|
|/lucky|GET|      |返回抽中的奖品ID|  

## 前端展现
使用[luch-canvas](https://100px.net/usage/js.html)抽奖插件。

## 压测结果记录
windows系统，8核16G
- 100线程,1次请求

| time | 并发数 | 成功数 | 失败数 |  qps   | 最长耗时 | 最短耗时 | 平均耗时 |
|------|--------|--------|--------|--------|----------|----------|----------|
| 0s   |  100   |  100   |   0    | 2069.55|  77.55   |  25.24   |  48.32   |


- 400线程,1次请求

| time | 并发数 | 成功数 | 失败数 |   qps   |  最长耗时  |  最短耗时  |  平均耗时  |
|------|--------|--------|--------|---------|-----------|-----------|-----------|
| 3s   | 4000   | 4000   |  0     | 2777.84 | 2581.70   |  144.74   | 1439.97   |
