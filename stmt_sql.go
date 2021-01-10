package main

const (
	CreateProfileTable = `
		create table user (id integer not null primary key, name text); 
		delete from user;
	`
	CreatImageTable = `
		create table image (id integer not null primary key, name text, path text, isPrivate integer, createTime text); 
		delete from image;
	`
)
