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