DROP TABLE public.raw_emails;
CREATE TABLE public.raw_emails
(
  id serial,
  sender varchar(256),
  topic varchar(256),
  subject varchar(256),
  message text
);


