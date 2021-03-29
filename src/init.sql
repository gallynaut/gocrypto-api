CREATE TABLE exchanges (
    full_name varchar(255) NOT NULL PRIMARY KEY,
    exchange_name varchar(255) NOT NULL,
    end_point varchar(255) NOT NULL
);

INSERT INTO exchanges(full_name, exchange_name, end_point) VALUES 
('BYBIT-MAIN', 'BYBIT', 'api.bybit.com'), 
('BYBIT-TEST', 'BYBIT', 'test.api.bybit.com');

ALTER DEFAULT PRIVILEGES [ FOR ROLE me] GRANT ALL ON TABLES TO bspu;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO me;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO me;