-- +goose Up
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS public.users (
	id uuid NOT NULL,
	"login" varchar(50) NOT NULL,
	"password" varchar(60) NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE active_tokens (
    jti uuid,
    user_id uuid NOT NULL,
    issued_at timestamp with time zone DEFAULT NOW(),
    expires_at timestamp with time zone NOT NULL,
	PRIMARY KEY (jti),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX active_tokens_user_id_idx ON public.active_tokens USING btree (user_id);

CREATE TABLE IF NOT EXISTS public.balance (
	id uuid NOT NULL,
	user_id uuid NOT NULL,
	current int NULL,
	withdrawn int NULL,
	PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES public.users(id)
);
CREATE INDEX account_user_id_idx ON public.balance USING btree (user_id);

CREATE TABLE IF NOT EXISTS public.orders (
	id uuid NOT NULL,
	user_id uuid NOT NULL,
	"number" varchar(255) NOT NULL,
	"status" varchar(50) NOT NULL,
	"accrual" int NULL,
	uploaded_at timestamp with time zone NOT NULL,
	PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES public.users(id)
);
CREATE INDEX orders_user_id_idx ON public.orders USING btree (user_id);

CREATE TABLE IF NOT EXISTS public.drawals (
	id uuid NOT NULL,
	user_id uuid NOT NULL,
	order_number char(11) NOT NULL,
	sum int NOT NULL,
	processed_at timestamp with time zone NOT NULL,
	PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES public.users(id)
);
CREATE INDEX drawals_user_id_idx ON public.drawals USING btree (user_id);

-- +goose Down
SELECT 'down SQL query';

DROP TABLE IF EXISTS public.active_tokens;
DROP TABLE IF EXISTS public.balance;
DROP TABLE IF EXISTS public.orders;
DROP TABLE IF EXISTS public.drawals;
DROP TABLE IF EXISTS public.users;
