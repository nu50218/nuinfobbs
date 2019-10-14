create table posts(
	number int not null,
	title text not null,
	url text not null,
	tweeted boolean not null default false,
	primary key (number)
);