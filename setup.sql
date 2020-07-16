DROP TABLE IF EXISTS users CASCADE ;
DROP TABLE IF EXISTS contact;
DROP TABLE IF EXISTS positive_users;

CREATE TABLE users (
    id bigserial primary key,
    joined bigint default FLOOR(EXTRACT(epoch FROM NOW() at time zone 'utc')*1000)
);

CREATE TABLE contact (
    person1 bigserial references users(id),
    person2 bigserial references users(id),
    contact_time bigint default (extract(epoch from now()) * 1000)
);

CREATE TABLE positive_users (
    id bigserial references users(id) unique,
    diagnosed bigint default FLOOR(EXTRACT(epoch FROM NOW() at time zone 'utc')*1000)
);


-- INSERT INTO users DEFAULT VALUES RETURNING id;

SELECT * FROM contact;
--
-- SELECT * FROM positive_users WHERE id=(SELECT person1 FROM contact WHERE person2=1) OR id=(SELECT person2 FROM contact WHERE person1=1);
--
-- SELECT FLOOR(EXTRACT(epoch FROM NOW() at time zone 'utc')*1000)