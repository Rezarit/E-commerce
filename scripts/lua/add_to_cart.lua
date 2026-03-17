-- 购物车添加商品Lua脚本
-- 参数：
-- KEYS[1]: 购物车key (cart:用户ID)
-- ARGV[1]: 商品ID
-- ARGV[2]: 要增加的数量
-- ARGV[3]: 最大允许数量（防止异常）

-- 使用HINCRBY实现真正的原子操作
local new_quantity = redis.call('HINCRBY', KEYS[1], ARGV[1], tonumber(ARGV[2]))

-- 验证数量范围
if new_quantity <= 0 then
    -- 如果数量小于等于0，回滚操作
    redis.call('HINCRBY', KEYS[1], ARGV[1], -tonumber(ARGV[2]))
    return -1  -- 数量必须大于0
end

-- 限制最大数量（防止异常数据）
local max_quantity = 10000  -- 硬编码最大数量
if new_quantity > max_quantity then
    -- 如果超过最大数量，回滚到最大值
    local excess = new_quantity - max_quantity
    redis.call('HINCRBY', KEYS[1], ARGV[1], -excess)
    new_quantity = max_quantity
end

-- 设置购物车过期时间（24小时）
redis.call('EXPIRE', KEYS[1], 86400)

-- 返回新数量
return new_quantity