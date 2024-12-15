create table if not exists tags
(
    name varchar(50) not null unique
);

create table if not exists articles
(
    id    serial primary key unique not null,
    title varchar(100)              not null,
    text  text                      not null
);

create table if not exists article_tags
(
    article_id int         not null,
    tag_name   varchar(50) not null,
    foreign key (article_id) references articles (id) on delete cascade,
    foreign key (tag_name) references tags (name) on delete cascade
);

create table if not exists courses
(
    id          serial primary key unique not null,
    title       varchar(100)              not null,
    poster      varchar(255),
    description varchar(255)
);

create table if not exists course_articles
(
    course_id  int not null,
    article_id int not null,
    primary key (course_id, article_id),
    foreign key (course_id) references courses (id) on delete cascade,
    foreign key (article_id) references articles (id) on delete cascade
);

create table if not exists course_tags
(
    course_id int         not null,
    tag_name  varchar(50) not null,
    foreign key (course_id) references courses (id) on delete cascade,
    foreign key (tag_name) references tags (name) on delete cascade
);
