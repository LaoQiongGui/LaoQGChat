CREATE SEQUENCE chat_content_seq;

-- 创建触发器函数
CREATE OR REPLACE FUNCTION increment_chat_content_serial()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.serial_number := nextval('chat_content_seq');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

ALTER SEQUENCE public.chat_content_seq
    OWNER TO laoqionggui;

-- 创建触发器
CREATE TRIGGER chat_content_serial_trigger
    BEFORE INSERT ON chat_content
    FOR EACH ROW
EXECUTE FUNCTION increment_chat_content_serial();

ALTER FUNCTION public.increment_chat_content_serial()
    OWNER TO laoqionggui;
