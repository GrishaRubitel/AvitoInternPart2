\c postgres;

SET search_path TO public;

CREATE TABLE public.employee (
	id varchar(100) NOT NULL,
	username varchar(50) NOT NULL,
	first_name varchar(50) NULL,
	last_name varchar(50) NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT employee_pkey PRIMARY KEY (id),
	CONSTRAINT employee_username_key UNIQUE (username)
);

CREATE TYPE public.organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
);

CREATE TABLE public.organization (
	id varchar(100) NOT NULL,
	"name" varchar(100) NOT NULL,
	description text NULL,
	"type" public.organization_type NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT organization_pkey PRIMARY KEY (id)
);

CREATE TABLE public.organization_responsible (
	id serial4 NOT NULL,
	organization_id varchar(100) NULL,
	user_id varchar(100) NULL,
	CONSTRAINT organization_responsible_pkey PRIMARY KEY (id),
	CONSTRAINT organization_responsible_organization_id_fkey FOREIGN KEY (organization_id) REFERENCES public.organization(id) ON DELETE CASCADE,
	CONSTRAINT organization_responsible_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.employee(id) ON DELETE CASCADE
);


-------------------------------------
---DDL для таблиц тендеров и бидов---
-------------------------------------

CREATE TYPE public.tender_status AS ENUM ('Created', 'Published', 'Closed');

CREATE TYPE public.tender_service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');

CREATE TABLE public.tenders (
	id varchar(100) NOT NULL,
	"name" varchar(100) NOT NULL,
	description varchar(500) NOT NULL,
	service_type public.tender_service_type NOT NULL,
	status public.tender_status NOT NULL,
	organization_id varchar(100) NOT NULL,
	"version" int4 DEFAULT 1 NOT NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT tenders_pkey PRIMARY KEY (id),
	CONSTRAINT tenders_version_check CHECK ((version >= 1))
);

CREATE TYPE public.bid_status AS ENUM ('Created', 'Published', 'Canceled', 'Approved', 'Rejected');

CREATE TYPE public.bid_decision AS ENUM ('Approved', 'Rejected'); --- In openapi.yaml there was no info about bid_decision's usage, so I took responsibility upon myself to add it to bid_reviews

CREATE TYPE public.bid_author_type AS ENUM ('Organization', 'User');

CREATE TABLE public.bids (
	id varchar(100) NOT NULL,
	"name" varchar(100) NOT NULL,
	description varchar(500) NOT NULL,
	status public.bid_status NOT NULL,
	tender_id varchar(100) NULL,
	author_type public.bid_author_type NOT NULL,
	author_id varchar(100) NOT NULL,
	"version" int4 DEFAULT 1 NOT NULL,
	created_at timestamp NOT NULL,
	CONSTRAINT bids_pkey PRIMARY KEY (id),
	CONSTRAINT bids_author_id_fkey FOREIGN KEY (author_id) REFERENCES public.employee(id),
	CONSTRAINT bids_tender_id_fkey FOREIGN KEY (tender_id) REFERENCES public.tenders(id)
);

CREATE TABLE public.bid_reviews (
	id varchar(100) NOT NULL,
	description varchar(1000) NOT NULL,
	created_at timestamp NOT NULL,
	bid_id varchar(100) NULL,
	decision public.bid_decision NOT NULL,
	CONSTRAINT bid_reviews_pkey PRIMARY KEY (id),
	CONSTRAINT bid_reviews_bid_id_fkey FOREIGN KEY (bid_id) REFERENCES public.bids(id)
);
------------------------
---DML with test data---
------------------------

INSERT INTO public.employee (id, username, first_name, last_name)
VALUES 
('1', 'john_doe', 'John', 'Doe'),
('2', 'jane_smith', 'Jane', 'Smith'),
('3', 'alex_brown', 'Alex', 'Brown'),
('4', 'lisa_jones', 'Lisa', 'Jones');

INSERT INTO public.organization (id, name, description, type)
VALUES 
('1', 'Tech Innovators', 'A technology company specializing in software solutions', 'LLC'),
('2', 'Green Builders', 'Eco-friendly construction company', 'IE'),
('3', 'Fast Delivery Co', 'Delivery services across the country', 'JSC');

INSERT INTO public.organization_responsible (organization_id, user_id)
VALUES 
('1', '1'),
('2', '2'),
('3', '3');

INSERT INTO public.tenders (id, name, description, service_type, status, organization_id, created_at)
VALUES 
('TND001', 'New Office Construction', 'Construction of a new office building', 'Construction', 'Created', '1', '2023-09-01 10:00:00'),
('TND002', 'Warehouse Delivery', 'Delivery of materials to the warehouse', 'Delivery', 'Published', '3', '2023-09-02 12:00:00'),
('TND003', 'Product Manufacture', 'Manufacture of electronic devices', 'Manufacture', 'Closed', '1', '2023-09-05 14:00:00');

INSERT INTO public.bids (id, name, description, status, tender_id, author_type, author_id, created_at)
VALUES 
('BID001', 'Bid for Office Construction', 'Our proposal for the office construction', 'Published', 'TND001', 'Organization', '1', '2023-09-03 09:00:00'),
('BID002', 'Delivery Bid', 'Proposal for the warehouse delivery', 'Created', 'TND002', 'User', '2', '2023-09-04 11:00:00'),
('BID003', 'Manufacturing Proposal', 'Proposal for product manufacturing', 'Approved', 'TND003', 'Organization', '3', '2023-09-06 16:00:00');

INSERT INTO public.bid_reviews (id, description, created_at, bid_id, decision)
VALUES 
('REV001', 'Approved the construction bid for meeting all requirements', '2023-09-07 10:00:00', 'BID001', 'Approved'),
('REV002', 'Rejected the delivery bid due to pricing issues', '2023-09-08 12:00:00', 'BID002', 'Rejected'),
('REV003', 'Approved the manufacturing proposal as it fits the timeline', '2023-09-09 14:00:00', 'BID003', 'Approved');
