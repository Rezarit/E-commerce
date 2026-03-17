-- 获取商品参数
local stock_key = KEYS[1]
local deduct_quantity=tonumber(ARGV[1])

-- 检查商品是否存在
local current_stock=redis.call('GET',stock_key)
if not current_stock then
    return -1  -- 商品不存在
end

-- 检查数量是否为正
if deduct_quantity<=0 then
    return -2  -- 数量必须为正整数
end

current_stock=tonumber(current_stock)

-- 检查库存是否充足
if current_stock<deduct_quantity then
    return -3
end

-- 减扣库存
local new_stock=redis.call('DECRBY',stock_key,deduct_quantity)

-- 返回新的库存数量
return new_stock