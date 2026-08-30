-- 加回库存（补偿脚本）
-- 参数：
-- KEYS[1]: 库存key
-- ARGV[1]: 要加回的数量

local stock_key = KEYS[1]
local add_quantity = tonumber(ARGV[1])

-- 检查数量是否合法
if add_quantity <= 0 then
    return -1  -- 数量必须为正整数
end

-- 关键：key 不存在说明库存已过期并重新预热，此时不该补偿（否则会凭空加库存导致超卖）
if not redis.call('EXISTS', stock_key) then
    return -2  -- 库存已重新预热，无需补偿
end

-- 加回库存
return redis.call('INCRBY', stock_key, add_quantity)
