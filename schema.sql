CREATE OR REPLACE FUNCTION reinitialize_schema() RETURNS void AS $$
BEGIN
    DROP TABLE IF EXISTS transaction_tags CASCADE;
    DROP TABLE IF EXISTS transactions CASCADE;
    DROP TABLE IF EXISTS tags CASCADE;
    DROP TABLE IF EXISTS categories CASCADE;
    DROP TABLE IF EXISTS accounts CASCADE;
    DROP TABLE IF EXISTS user_secrets CASCADE;
    DROP TABLE IF EXISTS users CASCADE;

    DROP DOMAIN IF EXISTS user_id_type CASCADE;
    DROP DOMAIN IF EXISTS currency_code CASCADE;
    DROP DOMAIN IF EXISTS transaction_type CASCADE;

    CREATE DOMAIN user_id_type AS VARCHAR(30)
    CHECK (
        attribute IS NOT NULL
        -- ensure no empty whitespace value
        AND length(trim(attribute)) != 0
        -- ensure lowercase alphanumeric and hyphen only and between 2 to 30 in length
        AND attribute ~ '^[a-z0-9][a-z0-9-]{0,28}[a-z0-9]$'
        -- ensure no multiple hyphens together
        AND attribute !~ '--+'
    );
    COMMENT ON DOMAIN user_id_type IS 'Valid user ID format: lowercase, alphanumeric, hyphens, 2-30 chars, no double hyphens';

    CREATE DOMAIN currency_code AS CHAR(3)
    CHECK (
        attribute IS NOT NULL
        AND length(attribute) = 3
        AND attribute = upper(attribute)
    );
    COMMENT ON DOMAIN currency_code IS 'ISO 4217 Currency Code (3 chars, uppercase)';

    CREATE DOMAIN transaction_type AS VARCHAR(20)
    CHECK (
        attribute IS NOT NULL
        AND attribute IN ('income', 'expense', 'transfer')
    );
    COMMENT ON DOMAIN transaction_type IS 'Type of transaction: income, expense, or transfer';


    CREATE TABLE users (
        id user_id_type PRIMARY KEY, -- username (uses custom domain for validation)
        email VARCHAR(255) NOT NULL UNIQUE,
        full_name VARCHAR(255) NOT NULL,
        display_name VARCHAR(255),
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

        CONSTRAINT users_email_validation CHECK (
            -- ensure no empty whitespace value
            length(trim(email)) != 0

            -- ensure valid email format:
            -- - lowercase only
            -- - local part starts/ends with alnum, allows ._%+-
            -- - domain part allows alnum, dot, hyphen
            -- - TLD at least 2 chars
            AND email ~ '^[a-z0-9](?:[a-z0-9._%+-]*[a-z0-9])?@[a-z0-9](?:[a-z0-9.-]*[a-z0-9])?\.[a-z]{2,}$'
        ),

        CONSTRAINT users_full_name_validation CHECK (
            -- ensure no empty whitespace value
            length(trim(full_name)) != 0

            -- ensure no leading and trailing spaces
            AND full_name = btrim(full_name)

            -- forbid tabs, newlines, carriage returns
            AND full_name !~ '[\t\n\r]'
        ),

        CONSTRAINT users_display_name_validation CHECK (
            display_name IS NULL OR (
                -- ensure no empty whitespace value
                length(trim(display_name)) != 0

                -- ensure no leading and trailing spaces
                AND display_name = btrim(display_name)

                -- forbid tabs, newlines, carriage returns
                AND display_name !~ '[\t\n\r]'
            )
        )
    );

    CREATE TABLE user_secrets (
        id user_id_type NOT NULL,
        value TEXT NOT NULL,

        CONSTRAINT user_secrets_pk PRIMARY KEY (id),
        CONSTRAINT user_secrets_id_fk FOREIGN KEY (id)
            REFERENCES users(id)
            ON UPDATE CASCADE
            ON DELETE CASCADE,

        CONSTRAINT user_secrets_value_format_check CHECK (
            length(trim(value)) != 0
        )
    );

    CREATE TABLE accounts (
        id SERIAL PRIMARY KEY,
        user_id user_id_type NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        name VARCHAR(255) NOT NULL,
        type VARCHAR(50) NOT NULL,
        currency currency_code NOT NULL DEFAULT 'USD',
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

        CONSTRAINT accounts_name_validation CHECK (
            length(trim(name)) != 0
            AND name = btrim(name)
        ),

        CONSTRAINT accounts_type_validation CHECK (
            length(trim(type)) != 0
        )
    );

    CREATE TABLE categories (
        id SERIAL PRIMARY KEY,
        user_id user_id_type NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        name VARCHAR(100) NOT NULL,
        type transaction_type NOT NULL CHECK (type IN ('income', 'expense')), -- Reuse transaction_type but limit to income/expense
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

        CONSTRAINT categories_name_validation CHECK (
            length(trim(name)) != 0
            AND name = btrim(name)
        )
    );

    CREATE TABLE tags (
        id SERIAL PRIMARY KEY,
        user_id user_id_type NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        name VARCHAR(50) NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE(user_id, name),

        CONSTRAINT tags_name_validation CHECK (
            length(trim(name)) != 0
            AND name = btrim(name)
        )
    );

    CREATE TABLE transactions (
        id SERIAL PRIMARY KEY,
        user_id user_id_type NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        from_account_id INTEGER REFERENCES accounts(id) ON DELETE SET NULL,
        to_account_id INTEGER REFERENCES accounts(id) ON DELETE SET NULL,
        category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
        amount DECIMAL(15, 2) NOT NULL,
        type transaction_type NOT NULL,
        description TEXT,
        transaction_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

        CONSTRAINT check_at_least_one_account CHECK (from_account_id IS NOT NULL OR to_account_id IS NOT NULL),

        CONSTRAINT transactions_amount_validation CHECK (
             amount != 0
        ),

        CONSTRAINT transactions_description_validation CHECK (
            description IS NULL OR (
                length(trim(description)) != 0
                AND description = btrim(description)
            )
        )
    );

    CREATE TABLE transaction_tags (
        transaction_id INTEGER NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
        tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
        PRIMARY KEY (transaction_id, tag_id)
    );

    CREATE INDEX idx_transactions_user ON transactions(user_id);
    CREATE INDEX idx_transactions_from_account ON transactions(from_account_id);
    CREATE INDEX idx_transactions_to_account ON transactions(to_account_id);
    CREATE INDEX idx_transactions_date ON transactions(transaction_date);
END;
$$ LANGUAGE plpgsql;

select reinitialize_schema();
