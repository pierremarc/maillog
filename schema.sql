
CREATE TABLE public.records (
  id serial,
  ts timestamp WITH TIME ZONE DEFAULT NOW(),
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


CREATE TABLE domains (
    id serial,
    http_name   varchar(256),
    mx_name varchar(256)
);