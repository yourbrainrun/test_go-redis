local member = KEYS[1]
local key = KEYS[2]
local flag
redis.log(redis.LOG_NOTICE,"member=",member)
flag = redis.call('EXISTS',member)
if flag < 1 then
    redis.call('ZREM',key,member)
end
