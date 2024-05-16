CREATE TABLE roles (
      id bigserial NOT NULL,
      name varchar(100) NOT NULL,
      description varchar(255) NULL,
      created_at timestamp NULL DEFAULT now(),
      created_by int8 NOT NULL,
      updated_at timestamp NULL,
      updated_by int8 NULL,
      deleted_at timestamp NULL,
      deleted_by int8 NULL,
      CONSTRAINT roles_pkey PRIMARY KEY (id)
);