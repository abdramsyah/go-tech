CREATE TYPE role_types AS ENUM ('user','non-user');
CREATE TYPE role_statuses AS ENUM ('active','inactive');

CREATE TABLE roles (
	id bigserial NOT NULL,
	"name" varchar(100) NOT NULL,
	description varchar(255) NULL,
	created_at timestamp DEFAULT now() NULL,
	created_by int8 NULL,
	updated_at timestamp NULL,
	updated_by int8 NULL,
	deleted_at timestamp NULL,
	deleted_by int8 NULL,
	jobdesc varchar(255) NULL,
	role_superior_id int4 NULL,
	status role_statuses DEFAULT 'active'::role_statuses NOT NULL,
	role_type role_types DEFAULT 'user'::role_types NOT NULL,
	CONSTRAINT roles_pkey PRIMARY KEY (id)
);