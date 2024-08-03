package database

// 该文件提供拦截功能的redis等服务,利用hash实现
//思路是：简历hash类型，key为ip或者userid，value为1，为每个key过期时间，
//每次请求时查询key值，如果不存在说明是这段时间内第一次访问，则可以抽奖，并将ip或者userid作为key，1作为value存入hash
//如果查询到的value值存在，则判断value是否超过限定的一定时间内的请求次数，如果超过则是为恶意请求，则返回错误，否则继续抽奖，value+1
