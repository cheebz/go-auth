-- public.user_groups definition

-- Drop tables

-- DROP TABLE public.user_groups;

-- DROP TABLE public."groups";

-- DROP TABLE public.users;

-- public.users definition

CREATE TABLE public.users (
	id serial NOT NULL,
	username varchar NOT NULL,
	"password" varchar NOT NULL,
	created timestamptz NOT NULL,
	uuid text NOT NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id)
);

-- public."groups" definition

CREATE TABLE public."groups" (
	id serial NOT NULL,
	"name" varchar NOT NULL,
	CONSTRAINT groups_pkey PRIMARY KEY (id)
);

-- public."user_groups" definition

CREATE TABLE public.user_groups (
	id serial NOT NULL,
	user_id int4 NOT NULL,
	group_id int4 NOT NULL,
	CONSTRAINT user_groups_pkey PRIMARY KEY (id)
);
CREATE INDEX fki_user_groups_group_id ON public.user_groups USING btree (group_id);
CREATE INDEX fki_user_groups_user_id ON public.user_groups USING btree (user_id);


-- public.user_groups foreign keys

ALTER TABLE public.user_groups ADD CONSTRAINT user_groups_group_id FOREIGN KEY (group_id) REFERENCES public."groups"(id);
ALTER TABLE public.user_groups ADD CONSTRAINT user_groups_user_id FOREIGN KEY (user_id) REFERENCES public.users(id);
