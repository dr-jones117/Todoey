--liquibase formatted sql

--changeset your.name:1 labels:example-label context:example-context
--comment: example comment
create table person (
    id int primary key not null,
    name varchar(50) not null,
    address1 varchar(50),
    address2 varchar(50),
    city varchar(30)
)
--rollback DROP TABLE person;

--changeset your.name:2 labels:example-label context:example-context
--comment: example comment
create table company (
    id int primary key not null,
    name varchar(50) not null,
    address1 varchar(50),
    address2 varchar(50),
    city varchar(30)
)
--rollback DROP TABLE company;

--changeset other.dev:3 labels:example-label context:example-context
--comment: example comment
alter table person add column country varchar(2)
--rollback ALTER TABLE person DROP COLUMN country;


--changeset jonah:4 labels:init
CREATE TABLE todo_list (
    id SERIAL PRIMARY KEY,
    title VARCHAR(127) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE todo (
    id SERIAL PRIMARY KEY, 
    completed BOOLEAN NOT NULL,
    value TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    todo_list_id INT NOT NULL, 

    FOREIGN KEY (todo_list_id) REFERENCES todo_list(id) ON DELETE CASCADE
);
--rollback DROP TABLE todo
--rollback DROP TABLE todo_list

--changeset jonah:5 labels:remove-tables
DROP TABLE person;
DROP TABLE company;
