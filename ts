                                        Table "public.users"
   Column   |            Type             | Collation | Nullable |              Default              
------------+-----------------------------+-----------+----------+-----------------------------------
 id         | integer                     |           | not null | nextval('users_id_seq'::regclass)
 username   | character varying(50)       |           | not null | 
 email      | character varying(100)      |           | not null | 
 created_at | timestamp without time zone |           |          | CURRENT_TIMESTAMP
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
    "idx_users_username" btree (username)
    "users_email_key" UNIQUE CONSTRAINT, btree (email)
    "users_username_key" UNIQUE CONSTRAINT, btree (username)
Referenced by:
    TABLE "posts" CONSTRAINT "posts_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE

