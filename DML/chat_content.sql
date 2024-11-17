CREATE SEQUENCE chat_content_seq;

-- 创建触发器函数
CREATE OR REPLACE FUNCTION increment_chat_content_serial()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.serial := nextval('chat_content_seq');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
CREATE TRIGGER chat_content_serial_trigger
    BEFORE INSERT ON chat_content
    FOR EACH ROW
EXECUTE FUNCTION increment_chat_content_serial();