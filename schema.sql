CREATE TABLE public.raw_emails
(
  id serial,
  sender varchar(256),
  topic varchar(256),
  subject varchar(256),
  message text
);


CREATE TABLE public.answers
(
  id serial,
  parent int,
  child int
);


CREATE TABLE public.raw_emails2
(
  id serial,
  sender varchar(256),
  topic varchar(256),
  subject varchar(256),
  message bytea
);

CREATE TABLE public.records (
  id serial,
  ts timestamp DEFAULT NOW(),
  sender varchar(256),
  recipient varchar(256),
  topic  varchar(256),
  domain varchar(256),
  header_date varchar(256),
  header_subject varchar(256),
  body text,
  parent int,
  payload bytea
);

CREATE TABLE public.attachments (
  id serial,
  record_id int,
  content_type varchar(256),
  file_name varchar(256)
);