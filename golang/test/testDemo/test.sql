事务语句
    假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）
    和 transactions 表（包含字段 id 主键， from_account_id 转出账户ID，
    to_account_id 转入账户ID， amount 转账金额）。
要求 ：
    编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。
    在事务中，需要先检查账户 A 的余额是否足够，如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，
    并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

-- 1. 开启事务
START TRANSACTION;

-- 2. 锁定账户 A 并获取余额（防止并发）
SELECT balance INTO @balance FROM accounts WHERE id = 1 FOR UPDATE;

-- 3. 判断余额是否足够（这步是关键）
-- 如果余额不足，立即回滚
IF @balance < 100 THEN
    ROLLBACK;
    LEAVE; -- 跳出当前循环
END IF;

-- 4. 扣除账户 A 的余额
UPDATE accounts SET balance = balance - 100 WHERE id = 1;

-- 5. 增加账户 B 的余额
UPDATE accounts SET balance = balance + 100 WHERE id = 2;

-- 6. 插入转账记录
INSERT INTO transactions (from_account_id, to_account_id, amount)
VALUES (1, 2, 100);

-- 7. 提交事务
COMMIT;

