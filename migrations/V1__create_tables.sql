CREATE TABLE balance (
	balanceid uuid,
	profileid uuid,
	operation double precision,
	operationtime timestamp DEFAULT NOW(),
	primary key (balanceid)
);