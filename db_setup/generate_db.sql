-- Database: wishtree

CREATE TABLE category (
	id int IDENTITY(0,1) NOT NULL,
	name varchar(100) COLLATE SQL_Latin1_General_CP1_CI_AS NOT NULL,
	description varchar(100) COLLATE SQL_Latin1_General_CP1_CI_AS NOT NULL,
	CONSTRAINT category_PK PRIMARY KEY (id)
);

CREATE TABLE tree_status (
	isOpen tinyint DEFAULT 1 NOT NULL
);


CREATE TABLE wish (
	id int IDENTITY(0,1) NOT NULL,
	x float NOT NULL,
	y float NOT NULL,
	[text] varchar(512) COLLATE SQL_Latin1_General_CP1_CI_AS NULL,
	author varchar(100) COLLATE SQL_Latin1_General_CP1_CI_AS NULL,
	zipCode varchar(10) COLLATE SQL_Latin1_General_CP1_CI_AS NULL,
	createdAt datetimeoffset NULL,
	category_id int NULL,
	isArchived tinyint DEFAULT 0 NOT NULL,
	CONSTRAINT wish_PK PRIMARY KEY (id)
);

ALTER TABLE wish ADD CONSTRAINT wish_FK FOREIGN KEY (category_id) REFERENCES category(id);