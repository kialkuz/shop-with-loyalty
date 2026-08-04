-- +goose Up
SELECT 'up SQL query';

CREATE TABLE IF NOT EXISTS public.users (
	id serial4 NOT NULL,
	"name" varchar(50) NOT NULL,
	"login" varchar(50) NOT NULL,
	"password" varchar(50) NOT NULL,
	PRIMARY KEY (id)
);

CREATE TABLE active_tokens (
    jti varchar(255),
    user_id int NOT NULL,
    issued_at timestamp with time zone DEFAULT NOW(),
    expires_at timestamp with time zone NOT NULL,
	PRIMARY KEY (jti),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX active_tokens_user_id_idx ON public.active_tokens USING btree (user_id);

CREATE TABLE IF NOT EXISTS public.account (
	id serial4 NOT NULL,
	user_id serial4 NOT NULL,
	balance int NULL,
	PRIMARY KEY (id),
    FOREIGN KEY (id) REFERENCES public.users(id)
);
CREATE INDEX account_user_id_idx ON public.account USING btree (user_id);

CREATE TABLE IF NOT EXISTS public.orders (
	id serial4 NOT NULL,
	user_id serial4 NOT NULL,
	"number" bpchar(16) NOT NULL,
	PRIMARY KEY (id),
    FOREIGN KEY (id) REFERENCES public.users(id)
);
CREATE INDEX orders_user_id_idx ON public.orders USING btree (user_id);

-- +goose Down
SELECT 'down SQL query';

DROP TABLE IF EXISTS public.active_tokens;
DROP TABLE IF EXISTS public.account;
DROP TABLE IF EXISTS public.orders;
DROP TABLE IF EXISTS public.users;
